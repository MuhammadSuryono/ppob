package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/shared/errors"
	"github.com/yontech/ppob/user-service/internal/dto"
	"github.com/yontech/ppob/user-service/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	tokenUserID, _ := c.Get("user_id")
	if uint(id) != tokenUserID.(uint) {
		role, _ := c.Get("role")
		if role != "admin" && role != "staff" {
			errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
			return
		}
	}

	resp, err := h.userService.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	tokenUserID, _ := c.Get("user_id")
	if uint(id) != tokenUserID.(uint) {
		role, _ := c.Get("role")
		if role != "admin" {
			errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
			return
		}
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.userService.UpdateUser(c.Request.Context(), uint(id), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetUserRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	resp, err := h.userService.GetUserRoles(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": resp})
}

func (h *UserHandler) AssignRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	if err := h.userService.AssignRole(c.Request.Context(), uint(id), req.RoleID); err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

func (h *UserHandler) ListRoles(c *gin.Context) {
	resp, err := h.userService.ListRoles(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": resp})
}

func (h *UserHandler) CreateRole(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.userService.CreateRole(c.Request.Context(), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusCreated, resp)
}