package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,phone_id"`
	FullName string `json:"full_name" binding:"required,min=2,max=255"`
	Password string `json:"password" binding:"required,min=8,password_complex"`
	PIN      string `json:"pin" binding:"required,pinformat"`
}

type RegisterResponse struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required_without=Phone"`
	Phone    string `json:"phone" binding:"required_without=Email"`
	Password string `json:"password" binding:"required"`
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

type VerifyOTPRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Type   string `json:"type" binding:"required,oneof=login register"`
}

type VerifyOTPResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
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
	OldPIN string `json:"old_pin" binding:"required,len=6"`
	NewPIN string `json:"new_pin" binding:"required,len=6"`
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