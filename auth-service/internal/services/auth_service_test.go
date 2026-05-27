package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/auth-service/config"
	"github.com/yontech/ppob/auth-service/internal/dto"
	"github.com/yontech/ppob/auth-service/internal/models"
	"github.com/yontech/ppob/auth-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.OTP{})
	db.AutoMigrate(&models.Wallet{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.UserRole{})
	db.AutoMigrate(&models.DeviceFingerprint{})

	// Seed required roles
	db.Create(&models.Role{Name: "Mitra", Status: "active"})

	return db
}

func setupTestConfig() *config.Config {
	// Generate temporary RSA key for testing
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	return &config.Config{
		JWTPrivateKey:    string(privatePEM),
		JWTPublicKey:     string(publicPEM),
		JWTExpire:        15 * time.Minute,
		RefreshExpire:    7 * 24 * time.Hour,
		OTPLength:        6,
		OTPExpireMinutes: 5,
		ServerPort:       "8080",
		DBHost:           "localhost",
		DBPort:           "5432",
		DBUser:           "postgres",
		DBPassword:       "postgres",
		DBName:           "test",
		RedisHost:        "localhost",
		RedisPort:        "6379",
		GinMode:          "test",
	}
}

type mockUserRepository struct {
	db *gorm.DB
}

func (m *mockUserRepository) Create(user *models.User) error {
	return m.db.Create(user).Error
}

func (m *mockUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := m.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *mockUserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := m.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *mockUserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := m.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *mockUserRepository) Update(user *models.User) error {
	return m.db.Save(user).Error
}

func (m *mockUserRepository) FindByEmailOrPhone(email, phone string) (*models.User, error) {
	var user models.User
	err := m.db.Where("email = ? OR phone = ?", email, phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type mockOTPRepository struct {
	db *gorm.DB
}

func (m *mockOTPRepository) Create(otp *models.OTP) error {
	return m.db.Create(otp).Error
}

func (m *mockOTPRepository) FindValidOTP(userID uint, code, otpType string) (*models.OTP, error) {
	var otp models.OTP
	err := m.db.Where("user_id = ? AND code = ? AND type = ? AND used_at IS NULL AND expires_at > ?",
		userID, code, otpType, time.Now()).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (m *mockOTPRepository) MarkAsUsed(otp *models.OTP) error {
	now := time.Now()
	otp.UsedAt = &now
	return m.db.Save(otp).Error
}

func (m *mockOTPRepository) DeleteExpired(userID uint, otpType string) error {
	return m.db.Where("user_id = ? AND type = ? AND expires_at < ?", userID, otpType, time.Now()).Delete(&models.OTP{}).Error
}

func (m *mockOTPRepository) FindLatestOTP(userID uint, otpType string) (*models.OTP, error) {
	var otp models.OTP
	err := m.db.Where("user_id = ? AND type = ?", userID, otpType).Order("created_at DESC").First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

type mockWalletRepository struct {
	db *gorm.DB
}

func (m *mockWalletRepository) Create(wallet *models.Wallet) error {
	return m.db.Create(wallet).Error
}

func (m *mockWalletRepository) FindByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := m.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (m *mockWalletRepository) Update(wallet *models.Wallet) error {
	return m.db.Save(wallet).Error
}

func (m *mockWalletRepository) FindByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := m.db.First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()
	req := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}

	// Mock OTP verification in redis
	s.Set("verified:"+req.RequestID, req.Phone)

	resp, err := authService.Register(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, resp.Email)
	}

	if resp.UserID == 0 {
		t.Error("Expected non-zero user ID")
	}

	if resp.Token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()
	req := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}

	// Mock OTP verification in redis
	s.Set("verified:"+req.RequestID, req.Phone)

	authService.Register(ctx, req)

	// Set again for second call because it's consumed
	s.Set("verified:"+req.RequestID, req.Phone)
	resp, err := authService.Register(ctx, req)

	if err != ErrUserExists {
		t.Fatalf("Expected ErrUserExists, got %v", err)
	}

	if resp != nil {
		t.Error("Expected nil response on duplicate user")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	authService.Register(ctx, registerReq)

	loginReq := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := authService.Login(ctx, loginReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Error("Expected non-empty access token")
	}

	if resp.RefreshToken == "" {
		t.Error("Expected non-empty refresh token")
	}

	if resp.UserID == 0 {
		t.Error("Expected non-zero user ID")
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	authService.Register(ctx, registerReq)

	loginReq := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	resp, err := authService.Login(ctx, loginReq)

	if err != ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got %v", err)
	}

	if resp != nil {
		t.Error("Expected nil response on invalid credentials")
	}
}

func TestAuthService_Login_NonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, nil, cfg)

	ctx := context.Background()

	loginReq := &dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	resp, err := authService.Login(ctx, loginReq)

	if err != ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got %v", err)
	}

	if resp != nil {
		t.Error("Expected nil response for non-existent user")
	}
}

func TestAuthService_VerifyOTP_Success(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "reg-req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	resp, _ := authService.Register(ctx, registerReq)

	var otpModel models.OTP
	db.Where("user_id = ? AND type = ?", resp.UserID, "register").First(&otpModel)

	verifyReq := &dto.VerifyOTPRequest{
		RequestID: "otp-req-1",
		Phone:     "081234567890",
		Code:   otpModel.Code,
		Type:   "register",
	}

	// Mock OTP in redis for VerifyOTP
	otpData := fmt.Sprintf("%s:%s:%s", verifyReq.Phone, otpModel.Code, verifyReq.Type)
	s.Set("otp:"+verifyReq.RequestID, otpData)

	verifyResp, err := authService.VerifyOTP(ctx, verifyReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if verifyResp != nil && !verifyResp.IsVerified {
		t.Error("Expected IsVerified to be true")
	}
}

func TestAuthService_VerifyOTP_InvalidCode(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "reg-req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	_, _ = authService.Register(ctx, registerReq)

	verifyReq := &dto.VerifyOTPRequest{
		RequestID: "req-1",
		Phone:     "081234567890",
		Code:   "000000",
		Type:   "register",
	}

	_, err = authService.VerifyOTP(ctx, verifyReq)

	if err != ErrInvalidOTP {
		t.Fatalf("Expected ErrInvalidOTP, got %v", err)
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	resp, _ := authService.Register(ctx, registerReq)

	claims, err := authService.ValidateToken(resp.Token)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if claims.Email != registerReq.Email {
		t.Errorf("Expected email %s, got %s", registerReq.Email, claims.Email)
	}

	if claims.UserID != resp.UserID {
		t.Errorf("Expected userID %d, got %d", resp.UserID, claims.UserID)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	cfg := setupTestConfig()

	authService := NewAuthService(nil, nil, nil, nil, nil, nil, cfg)

	_, err := authService.ValidateToken("invalid-token")

	if err != ErrInvalidToken {
		t.Fatalf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestAuthService_ChangePassword_Success(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	resp, _ := authService.Register(ctx, registerReq)

	err = authService.ChangePassword(ctx, resp.UserID, "password123", "newpassword456")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	loginReq := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "newpassword456",
	}

	_, err = authService.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Expected login with new password to succeed, got %v", err)
	}
}

func TestAuthService_ChangePassword_WrongOldPassword(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	resp, _ := authService.Register(ctx, registerReq)

	err = authService.ChangePassword(ctx, resp.UserID, "wrongpassword", "newpassword456")

	if err != ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestGenerateToken(t *testing.T) {
	cfg := setupTestConfig()

	email := "test@example.com"
	user := &models.User{
		ID:    1,
		Email: &email,
		Phone: "081234567890",
		Role:  "user",
	}

	token, expiresAt, err := generateToken(cfg, user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	if expiresAt.IsZero() {
		t.Error("Expected non-zero expiration time")
	}
}

func TestGenerateOTP(t *testing.T) {
	otp1 := generateOTP()
	otp2 := generateOTP()

	if len(otp1) != 6 {
		t.Errorf("Expected OTP length 6, got %d", len(otp1))
	}

	if otp1 == otp2 {
		t.Error("Expected different OTPs on consecutive calls")
	}
}

func generateToken(cfg *config.Config, user *models.User) (string, time.Time, error) {
	pk, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.JWTPrivateKey))
	expiresAt := time.Now().Add(cfg.JWTExpire)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"phone":   user.Phone,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(pk)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func TestAuthService_AuthorizeTransaction_Success(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	resp, _ := authService.Register(ctx, registerReq)

	authResp, err := authService.AuthorizeTransaction(ctx, resp.UserID, "123456")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if authResp.AuthorizeID == "" {
		t.Error("Expected non-empty AuthorizeID")
	}

	// Verify it's in redis
	key := fmt.Sprintf("transaction_authorize:%s", authResp.AuthorizeID)
	storedUserID, _ := redisClient.Get(ctx, key).Result()
	if storedUserID != fmt.Sprintf("%d", resp.UserID) {
		t.Errorf("Expected stored user ID %d, got %s", resp.UserID, storedUserID)
	}
}

func TestAuthService_AuthorizeTransaction_InvalidPIN(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	resp, _ := authService.Register(ctx, registerReq)

	_, err = authService.AuthorizeTransaction(ctx, resp.UserID, "000000")

	if err != ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := NewAuthService(userRepo, otpRepo, roleRepo, walletRepo, deviceRepo, redisClient, cfg)

	ctx := context.Background()

	registerReq := &dto.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "081234567890",
		Name:     "Test User",
		Password: "password123",
		PIN:      "123456",
		RequestID: "req-1",
	}
	s.Set("verified:"+registerReq.RequestID, registerReq.Phone)
	resp, _ := authService.Register(ctx, registerReq)

	refreshResp, err := authService.RefreshToken(ctx, resp.RefreshToken)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if refreshResp.Token == "" {
		t.Error("Expected non-empty new access token")
	}
}

func generateOTP() string {
	max := big.NewInt(999999)
	result, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", result)
}