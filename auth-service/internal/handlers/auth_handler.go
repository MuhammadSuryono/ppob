package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yontech/ppob/auth-service/internal/dto"
	"github.com/yontech/ppob/auth-service/internal/errors"
	"github.com/yontech/ppob/auth-service/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Initiate(c *gin.Context) {
	var req dto.InitiateAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.authService.InitiateAuth(c.Request.Context(), &req)
	if err != nil {
		errors.InternalError(c, "SYSTEM_INTERNAL", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) VerifyPassword(c *gin.Context) {
	var req dto.VerifyPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.authService.VerifyPassword(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			errors.Unauthorized(c, "AUTH_INVALID_CREDENTIALS", "Invalid password")
		case services.ErrVerificationRequired:
			errors.BadRequest(c, "AUTH_OTP_NOT_VERIFIED", "OTP belum diverifikasi atau tidak valid")
		default:
			errors.InternalError(c, "SYSTEM_INTERNAL", err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) VerifyPINLogin(c *gin.Context) {
	var req dto.PINLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.authService.VerifyPINLogin(c.Request.Context(), req.Phone, req.PIN, req.DeviceID)
	if err != nil {
		switch {
		case err == services.ErrInvalidCredentials:
			errors.Unauthorized(c, "AUTH_INVALID_CREDENTIALS", "Invalid PIN")
		case err == services.ErrDeviceNotTrusted:
			errors.Forbidden(c, "AUTH_DEVICE_NOT_TRUSTED", "Device not trusted for PIN login")
		default:
			errors.InternalError(c, "SYSTEM_INTERNAL", err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) VerifyCredential(c *gin.Context) {
	var req dto.VerifyCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.authService.VerifyCredential(c.Request.Context(), &req)
	if err != nil {
		switch {
		case err == services.ErrInvalidCredentials:
			errors.Unauthorized(c, "AUTH_INVALID_CREDENTIALS", "Invalid credentials")
		case err == services.ErrVerificationRequired:
			errors.BadRequest(c, "AUTH_OTP_NOT_VERIFIED", "OTP belum diverifikasi atau tidak valid")
		default:
			errors.InternalError(c, "SYSTEM_INTERNAL", err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "VALIDATION_MISSING_FIELD", err.Error())
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrUserExists:
			errors.BadRequest(c, "AUTH_USER_EXISTS", "User dengan nomor HP ini sudah terdaftar")
		case services.ErrVerificationRequired:
			errors.BadRequest(c, "AUTH_OTP_NOT_VERIFIED", "OTP belum diverifikasi atau tidak valid")
		default:
			errors.InternalError(c, "SYSTEM_INTERNAL", errors.GetErrorMessage("SYSTEM_INTERNAL"))
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email/phone or password",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.authService.SendOTP(c.Request.Context(), &req)
	if err != nil {
		errors.InternalError(c, "SYSTEM_INTERNAL", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.authService.VerifyOTP(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrInvalidOTP:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_otp",
				Message: "Invalid or expired OTP",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrInvalidToken, services.ErrTokenExpired:
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "invalid_token",
				Message: "Invalid or expired refresh token",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	tokenJTI := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token, _, err := new(jwt.Parser).ParseUnverified(parts[1], nil)
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if jti, ok := claims["jti"].(string); ok {
						tokenJTI = jti
					}
				}
			}
		}
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(uint), tokenJTI); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.LogoutResponse{Message: "Logged out successfully"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Old password is incorrect",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *AuthHandler) ChangePIN(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID not found in token",
		})
		return
	}

	var req dto.ChangePINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.authService.ChangePIN(c.Request.Context(), userID.(uint), req.OldPIN, req.NewPIN); err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Old PIN is incorrect",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN changed successfully"})
}

func (h *AuthHandler) AuthorizeTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.Unauthorized(c, "AUTH_UNAUTHORIZED", "User ID not found in token")
		return
	}

	var req dto.AuthorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.authService.AuthorizeTransaction(c.Request.Context(), userID.(uint), req.PIN)
	if err != nil {
		switch {
		case err == services.ErrInvalidCredentials:
			errors.Unauthorized(c, "AUTH_INVALID_CREDENTIALS", "Invalid PIN")
		case err == services.ErrUserNotFound:
			errors.NotFound(c, "AUTH_USER_NOT_FOUND", "User not found")
		default:
			errors.InternalError(c, "SYSTEM_INTERNAL", err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
