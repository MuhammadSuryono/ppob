package dto

import "time"

type UpdateUserRequest struct {
	FullName    string     `json:"full_name"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	Address     string     `json:"address"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type UserResponse struct {
	ID            uint       `json:"id"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	FullName      string     `json:"full_name"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	Avatar        string     `json:"avatar"`
	Address       string     `json:"address"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type RoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Permissions string `json:"permissions"`
}

type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions string `json:"permissions"`
	Status      string `json:"status"`
}

type AssignRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int64          `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}