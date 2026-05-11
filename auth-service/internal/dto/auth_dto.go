package dto

import "time"

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,phone_id"`
	FullName  string `json:"full_name" binding:"required,min=2,max=255"`
	Password  string `json:"password" binding:"required,min=8,password_complex"`
	PIN       string `json:"pin" binding:"required,pinformat"`
	DeviceID  string `json:"device_id"`
	RequestID string `json:"request_id" binding:"required"`
}

type RegisterResponse struct {
	UserID           uint   `json:"user_id"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	FullName         string `json:"full_name"`
	Token            string `json:"token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresAt        int64  `json:"expires_at"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required_without=Phone"`
	Phone    string `json:"phone" binding:"required_without=Email"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id"`
}

type InitiateAuthRequest struct {
	Phone      string `json:"phone" binding:"required,phone_id"`
	DeviceID   string `json:"device_id" binding:"required"`
	Fingerprint string `json:"fingerprint"`
}

type InitiateAuthResponse struct {
	UserID       uint `json:"user_id,omitempty"`
	IsRegistered bool `json:"is_registered"`
	IsTrusted    bool `json:"is_trusted"`
	RequiresOTP  bool `json:"requires_otp"`
}

type VerifyPasswordRequest struct {
	Phone     string `json:"phone" binding:"required,phone_id"`
	Password  string `json:"password" binding:"required"`
	DeviceID  string `json:"device_id" binding:"required"`
	RequestID string `json:"request_id" binding:"required"`
}

type LoginResponse struct {
	UserID      uint      `json:"user_id"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	FullName    string    `json:"full_name"`
	Token       string    `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt   int64     `json:"expires_at"`
}

type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required,phone_id"`
	Type  string `json:"type" binding:"required,oneof=login register"`
}

type SendOTPResponse struct {
	RequestID string `json:"request_id"`
	ExpiresAt int64  `json:"expires_at"`
}

type VerifyOTPRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Phone     string `json:"phone" binding:"required,phone_id"`
	Code      string `json:"code" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=login register"`
}

type VerifyOTPResponse struct {
	RequestID  string `json:"request_id"`
	IsVerified bool   `json:"is_verified"`
	IsNewUser  bool   `json:"is_new_user"`
}

type VerifyCredentialRequest struct {
	Phone      string `json:"phone" binding:"required,phone_id"`
	RequestID  string `json:"request_id" binding:"required"`
	DeviceID   string `json:"device_id" binding:"required"`
	AuthMethod string `json:"auth_method" binding:"required,oneof=password pin"`
	Value      string `json:"value" binding:"required"`
}

type PINLoginRequest struct {
	Phone     string `json:"phone" binding:"required,phone_id"`
	PIN       string `json:"pin" binding:"required,pinformat"`
	DeviceID  string `json:"device_id" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	Token       string `json:"token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ChangePINRequest struct {
	OldPIN string `json:"old_pin" binding:"required,pinformat"`
	NewPIN string `json:"new_pin" binding:"required,pinformat"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
	IssuedAt int64  `json:"iat"`
	JTI      string `json:"jti"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type LoginWithOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type LoginWithOTPResponse struct {
	UserID    uint   `json:"user_id"`
	NeedOTP   bool   `json:"need_otp"`
	ExpiresAt int64  `json:"expires_at"`
}

type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

type AuthResult struct {
	User         interface{}
	AccessToken  *TokenInfo
	RefreshToken *TokenInfo
}