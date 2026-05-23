package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/user-service/config"
	"github.com/yontech/ppob/user-service/internal/dto"
	"github.com/yontech/ppob/user-service/internal/models"
	"github.com/yontech/ppob/user-service/internal/repository"
)

type RoleSwitchHandler struct {
	userRepo  *repository.UserRepository
	roleRepo  *repository.RoleRepository
	cfg       *config.Config
}

func NewRoleSwitchHandler(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository, cfg *config.Config) *RoleSwitchHandler {
	return &RoleSwitchHandler{
		userRepo: userRepo,
		roleRepo: roleRepo,
		cfg:      cfg,
	}
}

func (h *RoleSwitchHandler) SwitchRole(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.RespondWithError(c, errors.NewAppError("AUTH_TOKEN_EXPIRED", nil))
		return
	}

	var req dto.SwitchRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	userRoles, err := h.userRepo.GetUserRoles(userID.(uint))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	var targetRole *models.UserRole
	for _, ur := range userRoles {
		if ur.RoleID == req.RoleID {
			targetRole = &ur
			break
		}
	}

	if targetRole == nil {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	role, err := h.roleRepo.GetRoleByID(req.RoleID)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	wallet, err := h.userRepo.GetWalletByUserAndRole(userID.(uint), req.RoleID)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	if wallet == nil {
		errors.RespondWithError(c, errors.NewAppError("TRANSACTION_NOT_FOUND", nil))
		return
	}

	claims := jwt.MapClaims{
		"user_id":  userID.(uint),
		"email":    c.GetString("email"),
		"phone":    c.GetString("phone"),
		"role":     role.RoleName,
		"role_id":  req.RoleID,
		"wallet_id": wallet.ID,
		"exp":      time.Now().Add(h.cfg.JWTExpire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(h.cfg.PrivateKey)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SwitchRoleResponse{
		Message:   "Role switched successfully",
		RoleID:    req.RoleID,
		RoleName:  role.RoleName,
		WalletID:  wallet.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(h.cfg.JWTExpire).Unix(),
	})
}