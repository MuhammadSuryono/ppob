package services

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/user-service/config"
	"github.com/yontech/ppob/user-service/internal/dto"
	"github.com/yontech/ppob/user-service/internal/models"
	"github.com/yontech/ppob/user-service/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
	ErrRoleExists   = errors.New("role already exists")
)

type UserService struct {
	userRepo  *repository.UserRepository
	roleRepo  *repository.RoleRepository
	redis     *redis.Client
	cfg       *config.Config
}

func NewUserService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	redis *redis.Client,
	cfg *config.Config,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		redis:    redis,
		cfg:      cfg,
	}
}

func (s *UserService) GetUser(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &dto.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		FullName:      user.FullName,
		Role:          user.Role,
		Status:        user.Status,
		Avatar:        user.Avatar,
		Address:       user.Address,
		DateOfBirth:   user.DateOfBirth,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	s.invalidateUserCache(ctx, userID)

	return &dto.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		FullName:      user.FullName,
		Role:          user.Role,
		Status:        user.Status,
		Avatar:        user.Avatar,
		Address:       user.Address,
		DateOfBirth:   user.DateOfBirth,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
	users, total, err := s.userRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Phone:         user.Phone,
			FullName:      user.FullName,
			Role:          user.Role,
			Status:        user.Status,
			Avatar:        user.Avatar,
			Address:       user.Address,
			DateOfBirth:   user.DateOfBirth,
			LastLoginAt:   user.LastLoginAt,
			CreatedAt:     user.CreatedAt,
		}
	}

	return &dto.ListUsersResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *UserService) GetUserRoles(ctx context.Context, userID uint) ([]dto.RoleResponse, error) {
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			Status:      role.Status,
		}
	}

	return roleResponses, nil
}

func (s *UserService) AssignRole(ctx context.Context, userID uint, roleID uint) error {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return ErrRoleNotFound
	}

	if err := s.roleRepo.AssignRole(userID, roleID); err != nil {
		return err
	}

	s.invalidateUserCache(ctx, userID)

	user, _ := s.userRepo.FindByID(userID)
	if user != nil {
		user.Role = role.Name
		s.userRepo.Update(user)
	}

	return nil
}

func (s *UserService) ListRoles(ctx context.Context) ([]dto.RoleResponse, error) {
	roles, err := s.roleRepo.List()
	if err != nil {
		return nil, err
	}

	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			Status:      role.Status,
		}
	}

	return roleResponses, nil
}

func (s *UserService) CreateRole(ctx context.Context, req *dto.RoleRequest) (*dto.RoleResponse, error) {
	existing, _ := s.roleRepo.FindByName(req.Name)
	if existing != nil {
		return nil, ErrRoleExists
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		Status:     "active",
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
		Status:      role.Status,
	}, nil
}

func (s *UserService) invalidateUserCache(ctx context.Context, userID uint) {
	key := "user:" + string(rune(userID))
	s.redis.Del(ctx, key)
}