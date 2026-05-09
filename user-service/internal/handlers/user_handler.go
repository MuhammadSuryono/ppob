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

// Staff Management Handlers

func (h *UserHandler) ListStaff(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.DefaultQuery("status", "active")

	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	var mitraID uint
	if role == "mitra" {
		mitraID = userID.(uint)
	} else if role == "admin" {
		mitraID = 0 // admin sees all staff
	} else {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	resp, err := h.userService.ListStaff(c.Request.Context(), mitraID, limit, offset, status)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetStaff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	resp, err := h.userService.GetStaff(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("USER_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) CreateStaff(c *gin.Context) {
	var req dto.StaffCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	// Only Mitra can create staff
	role, _ := c.Get("role")
	if role != "mitra" {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}

	mitraID, _ := c.Get("user_id")

	resp, err := h.userService.CreateStaff(c.Request.Context(), mitraID.(uint), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) UpdateStaff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	var req dto.StaffUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_JSON_INVALID", nil))
		return
	}

	resp, err := h.userService.UpdateStaff(c.Request.Context(), uint(id), &req)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetStaffStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("VALIDATION_MISSING_FIELD", map[string]interface{}{"field": "id"}))
		return
	}

	resp, err := h.userService.GetStaffStats(c.Request.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("USER_NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPendingStaffCount returns number of staff awaiting approval (status = 'pending')
func (h *UserHandler) GetPendingStaffCount(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "mitra" && role != "admin" {
		errors.RespondWithError(c, errors.NewAppError("AUTH_INSUFFICIENT_PERMISSION", nil))
		return
	}
	count, err := h.userService.GetPendingStaffCount(c.Request.Context())
	if err != nil {
		errors.RespondWithError(c, errors.NewAppError("SYSTEM_INTERNAL", nil))
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}