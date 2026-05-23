package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/auth-service/config"
	"github.com/yontech/ppob/auth-service/internal/dto"
	"github.com/yontech/ppob/auth-service/internal/models"
	"github.com/yontech/ppob/auth-service/internal/repository"
	"github.com/yontech/ppob/shared/events"
	"github.com/yontech/ppob/shared/security"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidOTP           = errors.New("invalid OTP")
	ErrOTPExpired           = errors.New("OTP expired")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrVerificationRequired = errors.New("OTP verification required")
	ErrDeviceNotTrusted     = errors.New("device not trusted for PIN login")
)

type AuthService struct {
	userRepo       *repository.UserRepository
	otpRepo        *repository.OTPRepository
	walletRepo     *repository.WalletRepository
	deviceRepo     *repository.DeviceRepository
	redis          *redis.Client
	cfg            *config.Config
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	eventPublisher *events.EventPublisher
}

func NewAuthService(
	userRepo *repository.UserRepository,
	otpRepo *repository.OTPRepository,
	walletRepo *repository.WalletRepository,
	deviceRepo *repository.DeviceRepository,
	redis *redis.Client,
	cfg *config.Config,
) *AuthService {
	svc := &AuthService{
		userRepo:       userRepo,
		otpRepo:        otpRepo,
		walletRepo:     walletRepo,
		deviceRepo:     deviceRepo,
		redis:          redis,
		cfg:            cfg,
		eventPublisher: events.NewEventPublisher(redis),
	}

	if cfg.JWTPrivateKey != "" {
		pk, err := security.ParseRSAPrivateKey(cfg.JWTPrivateKey)
		if err != nil {
			fmt.Printf("Warning: failed to parse JWT private key: %v\n", err)
		} else {
			svc.privateKey = pk
		}
	}

	if cfg.JWTPublicKey != "" {
		pub, err := security.ParseRSAPublicKey(cfg.JWTPublicKey)
		if err != nil {
			fmt.Printf("Warning: failed to parse JWT public key: %v\n", err)
		} else {
			svc.publicKey = pub
		}
	}

	return svc
}

func (s *AuthService) InitiateAuth(ctx context.Context, req *dto.InitiateAuthRequest) (*dto.InitiateAuthResponse, error) {
	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		return &dto.InitiateAuthResponse{
			IsRegistered: false,
			IsTrusted:    false,
			RequiresOTP:  true,
		}, nil
	}

	device, err := s.deviceRepo.FindByFingerprint(user.ID, req.DeviceID)
	isTrusted := false
	if err == nil && device != nil && device.IsTrusted {
		isTrusted = true
	}

	return &dto.InitiateAuthResponse{
		UserID:       user.ID,
		IsRegistered: true,
		IsTrusted:    isTrusted,
		RequiresOTP:  !isTrusted,
	}, nil
}

func (s *AuthService) VerifyPassword(ctx context.Context, req *dto.VerifyPasswordRequest) (*dto.LoginResponse, error) {
	// Check if phone matches requestID verification
	verifiedPhone, err := s.redis.Get(ctx, "verified:"+req.RequestID).Result()
	if err != nil || verifiedPhone != req.Phone {
		return nil, ErrVerificationRequired
	}

	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Mark device as trusted after successful password validation
	s.upsertDeviceTrust(ctx, user.ID, req.DeviceID)

	// Consume verification flag
	s.redis.Del(ctx, "verified:"+req.RequestID)

	return s.completeAuth(ctx, user)
}

func (s *AuthService) VerifyPINLogin(ctx context.Context, phone, pin, deviceID string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify device trust
	device, err := s.deviceRepo.FindByFingerprint(user.ID, deviceID)
	if err != nil || device == nil || !device.IsTrusted {
		return nil, ErrDeviceNotTrusted
	}

	if match, err := security.VerifyPIN(pin, user.PinHash); err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	return s.completeAuth(ctx, user)
}

func (s *AuthService) VerifyCredential(ctx context.Context, req *dto.VerifyCredentialRequest) (*dto.LoginResponse, error) {
	// Check if phone matches requestID verification
	verifiedPhone, err := s.redis.Get(ctx, "verified:"+req.RequestID).Result()
	if err != nil || verifiedPhone != req.Phone {
		return nil, ErrVerificationRequired
	}

	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify credential based on auth method
	switch req.AuthMethod {
	case "password":
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Value)); err != nil {
			return nil, ErrInvalidCredentials
		}
	case "pin":
		if match, err := security.VerifyPIN(req.Value, user.PinHash); err != nil || !match {
			return nil, ErrInvalidCredentials
		}
	default:
		return nil, ErrInvalidCredentials
	}

	// Mark device as trusted
	s.upsertDeviceTrust(ctx, user.ID, req.DeviceID)

	// Consume verification flag
	s.redis.Del(ctx, "verified:"+req.RequestID)

	return s.completeAuth(ctx, user)
}

func (s *AuthService) upsertDeviceTrust(ctx context.Context, userID uint, deviceID string) {
	if deviceID == "" {
		return
	}

	device, err := s.deviceRepo.FindByFingerprint(userID, deviceID)
	if err != nil {
		// New device
		if createErr := s.deviceRepo.Create(&models.DeviceFingerprint{
			UserID:          userID,
			FingerprintHash: deviceID,
			IsTrusted:       true,
			FirstSeen:       time.Now(),
			LastSeen:        time.Now(),
		}); createErr != nil {
			fmt.Printf("Failed to create device trust record for user %d: %v\n", userID, createErr)
		}
	} else if !device.IsTrusted {
		device.IsTrusted = true
		device.LastSeen = time.Now()
		if updateErr := s.deviceRepo.Update(device); updateErr != nil {
			fmt.Printf("Failed to update device trust for user %d: %v\n", userID, updateErr)
		}
	}
}

func (s *AuthService) completeAuth(ctx context.Context, user *models.User) (*dto.LoginResponse, error) {
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	accessToken, accessExpiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	s.storeRefreshToken(ctx, user.ID, refreshToken, refreshExpiresAt)

	return &dto.LoginResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Phone:        user.Phone,
		Name:         user.Name,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if phone matches requestID verification
	verifiedPhone, err := s.redis.Get(ctx, "verified:"+req.RequestID).Result()
	if err != nil || verifiedPhone != req.Phone {
		return nil, ErrVerificationRequired
	}

	existingUser, _ := s.userRepo.FindByEmailOrPhone(req.Email, req.Phone)
	if existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	hashedPIN, err := security.HashPIN(req.PIN)
	if err != nil {
		return nil, fmt.Errorf("failed to hash PIN: %w", err)
	}

	user := &models.User{
		Email:         req.Email,
		Phone:         req.Phone,
		PasswordHash:  string(hashedPassword),
		Name:          req.Name,
		PinHash:       hashedPIN,
		Role:          "Mitra",
		Status:        "active",
		EmailVerified: true, // OTP already verified
		PhoneVerified: true, // OTP already verified
	}

	err = s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		txUserRepo := repository.NewUserRepository(tx)
		if err := txUserRepo.Create(user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if req.DeviceID != "" {
			device := &models.DeviceFingerprint{
				UserID:          user.ID,
				FingerprintHash: req.DeviceID,
				IsTrusted:       true, // Initial device after OTP is trusted
				FirstSeen:       time.Now(),
				LastSeen:        time.Now(),
			}
			txDeviceRepo := repository.NewDeviceRepository(tx)
			if err := txDeviceRepo.Create(device); err != nil {
				return fmt.Errorf("failed to create device trust record: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish user.registered event for asynchronous processing (e.g., wallet creation)
	s.eventPublisher.Publish(ctx, "user_stream", "user.registered", map[string]interface{}{
		"user_id": user.ID,
		"phone":   user.Phone,
		"email":   user.Email,
		"role":    user.Role,
	})

	// Consume verification flag
	s.redis.Del(ctx, "verified:"+req.RequestID)

	accessToken, accessExpiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	s.storeRefreshToken(ctx, user.ID, refreshToken, refreshExpiresAt)

	return &dto.RegisterResponse{
		UserID:           user.ID,
		Email:            user.Email,
		Phone:            user.Phone,
		Name:             user.Name,
		Token:            accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        accessExpiresAt.Unix(),
		RefreshExpiresAt: refreshExpiresAt.Unix(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *models.User
	var err error

	if req.Email != "" {
		user, err = s.userRepo.FindByEmail(req.Email)
	} else {
		user, err = s.userRepo.FindByPhone(req.Phone)
	}

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	accessToken, accessExpiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	s.storeRefreshToken(ctx, user.ID, refreshToken, refreshExpiresAt)

	return &dto.LoginResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Phone:        user.Phone,
		Name:         user.Name,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}

func (s *AuthService) SendOTP(ctx context.Context, req *dto.SendOTPRequest) (*dto.SendOTPResponse, error) {
	otpCode := s.generateOTP()
	requestID := uuid.New().String()

	// Store in Redis: otp:request_id -> {phone, code, type}
	otpData := fmt.Sprintf("%s:%s:%s", req.Phone, otpCode, req.Type)
	err := s.redis.Set(ctx, "otp:"+requestID, otpData, time.Duration(s.cfg.OTPExpireMinutes)*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	fmt.Printf("OTP for %s [%s]: %s (expires in %d minutes)\n", req.Phone, req.Type, otpCode, s.cfg.OTPExpireMinutes)

	return &dto.SendOTPResponse{
		RequestID: requestID,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.OTPExpireMinutes) * time.Minute).Unix(),
	}, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, req *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	val, err := s.redis.Get(ctx, "otp:"+req.RequestID).Result()
	if err != nil {
		return nil, ErrInvalidOTP
	}

	parts := strings.Split(val, ":")
	if len(parts) != 3 {
		return nil, ErrInvalidOTP
	}

	storedPhone := parts[0]
	storedCode := parts[1]
	storedType := parts[2]

	if storedPhone != req.Phone || storedCode != req.Code || storedType != req.Type {
		return nil, ErrInvalidOTP
	}

	// Delete OTP after successful verification
	s.redis.Del(ctx, "otp:"+req.RequestID)

	// Mark as verified in Redis: verified:request_id -> phone
	err = s.redis.Set(ctx, "verified:"+req.RequestID, req.Phone, 10*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	// Check if this is a new user or existing user
	_, err = s.userRepo.FindByPhone(req.Phone)
	isNewUser := err != nil

	return &dto.VerifyOTPResponse{
		RequestID:  req.RequestID,
		IsVerified: true,
		IsNewUser:  isNewUser,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	claims, err := ValidateTokenStatic(refreshToken, s.publicKey)
	if err != nil {
		return nil, ErrInvalidToken
	}

	tokenKey := fmt.Sprintf("refresh_token:%d", claims.UserID)
	storedToken, err := s.redis.Get(ctx, tokenKey).Result()
	if err != nil || storedToken != refreshToken {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	newToken, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		Token:     newToken,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uint, tokenJTI string) error {
	tokenKey := fmt.Sprintf("refresh_token:%d", userID)
	if err := s.redis.Del(ctx, tokenKey).Err(); err != nil {
		return err
	}

	if tokenJTI != "" {
		blacklistKey := fmt.Sprintf("access_blacklist:%s", tokenJTI)
		expiration := 24 * time.Hour
		return s.redis.Set(ctx, blacklistKey, "1", expiration).Err()
	}

	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}

func (s *AuthService) ChangePIN(ctx context.Context, userID uint, oldPIN, newPIN string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if match, err := security.VerifyPIN(oldPIN, user.PinHash); err != nil || !match {
		return ErrInvalidCredentials
	}

	hashedPIN, err := security.HashPIN(newPIN)
	if err != nil {
		return fmt.Errorf("failed to hash PIN: %w", err)
	}

	user.PinHash = hashedPIN
	return s.userRepo.Update(user)
}

func (s *AuthService) AuthorizeTransaction(ctx context.Context, userID uint, pin string) (*dto.AuthorizeResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if match, err := security.VerifyPIN(pin, user.PinHash); err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	authorizeID := uuid.New().String()
	ttl := 5 * time.Minute
	expiresAt := time.Now().Add(ttl).Unix()

	// Store in Redis: key "transaction_authorize:<authorizeID>" -> userID
	key := fmt.Sprintf("transaction_authorize:%s", authorizeID)
	if err := s.redis.Set(ctx, key, userID, ttl).Err(); err != nil {
		return nil, fmt.Errorf("failed to store authorize id: %w", err)
	}

	return &dto.AuthorizeResponse{
		AuthorizeID: authorizeID,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*dto.TokenClaims, error) {
	return ValidateTokenStatic(tokenString, s.publicKey)
}

func ValidateTokenStatic(tokenString string, key *rsa.PublicKey) (*dto.TokenClaims, error) {
	if key == nil {
		return nil, fmt.Errorf("public key not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, ErrInvalidToken
		}

		if int64(exp) < time.Now().Unix() {
			return nil, ErrTokenExpired
		}

		userID, _ := claims["user_id"].(float64)
		email, _ := claims["email"].(string)
		phone, _ := claims["phone"].(string)
		role, _ := claims["role"].(string)
		iat, _ := claims["iat"].(float64)

		return &dto.TokenClaims{
			UserID:   uint(userID),
			Email:    email,
			Phone:    phone,
			Role:     role,
			Exp:      int64(exp),
			IssuedAt: int64(iat),
		}, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthService) generateToken(user *models.User) (string, time.Time, error) {
	if s.privateKey == nil {
		return "", time.Time{}, fmt.Errorf("private key not configured")
	}

	expiresAt := time.Now().Add(s.cfg.JWTExpire)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"phone":   user.Phone,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, time.Time, error) {
	if s.privateKey == nil {
		return "", time.Time{}, fmt.Errorf("private key not configured")
	}

	expiresAt := time.Now().Add(s.cfg.RefreshExpire)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) generateOTP() string {
	max := big.NewInt(999999)
	result, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", result)
}

func (s *AuthService) storeRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) {
	tokenKey := fmt.Sprintf("refresh_token:%d", userID)
	s.redis.Set(ctx, tokenKey, token, expiresAt.Sub(time.Now()))
}

func GenerateSalt() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
