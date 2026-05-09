package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/auth-service/config"
	"github.com/yontech/ppob/auth-service/internal/dto"
	"github.com/yontech/ppob/auth-service/internal/models"
	"github.com/yontech/ppob/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidOTP        = errors.New("invalid OTP")
	ErrOTPExpired        = errors.New("OTP expired")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
)

type AuthService struct {
	userRepo  *repository.UserRepository
	otpRepo   *repository.OTPRepository
	walletRepo *repository.WalletRepository
	redis     *redis.Client
	cfg       *config.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	otpRepo *repository.OTPRepository,
	walletRepo *repository.WalletRepository,
	redis *redis.Client,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		otpRepo:    otpRepo,
		walletRepo: walletRepo,
		redis:      redis,
		cfg:        cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existingUser, _ := s.userRepo.FindByEmailOrPhone(req.Email, req.Phone)
	if existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash PIN: %w", err)
	}

	user := &models.User{
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     string(hashedPassword),
		FullName:    req.FullName,
		PIN:          string(hashedPIN),
		Role:         "user",
		Status:       "active",
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	wallet := &models.Wallet{
		UserID:   user.ID,
		Balance: 0,
		Status:  "active",
	}
	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	otpCode := s.generateOTP()
	otp := &models.OTP{
		UserID:    user.ID,
		Code:     otpCode,
		Type:     "register",
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.OTPExpireMinutes) * time.Minute),
	}
	if err := s.otpRepo.Create(otp); err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	fmt.Printf("OTP for registration: %s (expires in %d minutes)\n", otpCode, s.cfg.OTPExpireMinutes)

	return &dto.RegisterResponse{
		UserID:    user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
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
		UserID:        user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		FullName:      user.FullName,
		Token:         accessToken,
		RefreshToken:  refreshToken,
		ExpiresAt:     accessExpiresAt.Unix(),
	}, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, req *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	otp, err := s.otpRepo.FindValidOTP(req.UserID, req.Code, req.Type)
	if err != nil {
		return nil, ErrInvalidOTP
	}

	if err := s.otpRepo.MarkAsUsed(otp); err != nil {
		return nil, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Type == "register" {
		user.EmailVerified = true
		user.PhoneVerified = true
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
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

	return &dto.VerifyOTPResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	claims, err := s.validateToken(refreshToken)
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

func (s *AuthService) ChangePIN(ctx context.Context, userID uint, oldPIN, newPIN string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(oldPIN)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(newPIN), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash PIN: %w", err)
	}

	user.PIN = string(hashedPIN)
	return s.userRepo.Update(user)
}

func (s *AuthService) ValidateToken(tokenString string) (*dto.TokenClaims, error) {
	return s.validateToken(tokenString)
}

func (s *AuthService) generateToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.cfg.JWTExpire)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"phone":   user.Phone,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.cfg.RefreshExpire)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) validateToken(tokenString string) (*dto.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
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

		return &dto.TokenClaims{
			UserID:   uint(claims["user_id"].(float64)),
			Email:    claims["email"].(string),
			Phone:    claims["phone"].(string),
			Role:     claims["role"].(string),
			Exp:      int64(exp),
			IssuedAt: int64(claims["iat"].(float64)),
		}, nil
	}

	return nil, ErrInvalidToken
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