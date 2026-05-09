package services

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/user-service/config"
	"github.com/yontech/ppob/user-service/internal/dto"
	"github.com/yontech/ppob/user-service/internal/models"
	"github.com/yontech/ppob/user-service/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
	ErrRoleExists   = errors.New("role already exists")
)

type UserService struct {
	userRepo      *repository.UserRepository
	roleRepo      *repository.RoleRepository
	marginService *MarginService
	redis         *redis.Client
	cfg           *config.Config
	db            *gorm.DB
}

func NewUserService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	marginService *MarginService,
	redis *redis.Client,
	cfg *config.Config,
) *UserService {
	return &UserService{
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		marginService: marginService,
		redis:         redis,
		cfg:           cfg,
		db:            nil,
	}
}

func (s *UserService) SetDB(db *gorm.DB) {
	s.db = db
}

// GetUser retrieves a user by ID
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

// UpdateUser updates user profile
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

// ListUsers returns paginated list of users
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

// GetUserRoles returns roles for a user
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

// AssignRole assigns a role to a user
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

// ListRoles returns all roles
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

// CreateRole creates a new role
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

// ========== Staff Management ==========

// ListStaff returns list of staff users for a Mitra
func (s *UserService) ListStaff(ctx context.Context, mitraID uint, limit, offset int, status string) (*dto.ListUsersResponse, error) {
	users, total, err := s.userRepo.FindStaffUsers(mitraID, limit, offset, status)
	if err != nil {
		return nil, err
	}

	staffResponses := make([]dto.StaffResponse, len(users))
	for i, user := range users {
		balance, _ := s.userRepo.GetWalletBalance(user.ID)
		todayCount, todayAmount, _ := s.userRepo.GetStaffTodayStats(user.ID)
		maxCount, maxAmount, _ := s.userRepo.GetStaffDailyLimit(user.ID)
		marginSetting, _ := s.marginService.GetStaffMarginSettings(user.ID)

		staffResponses[i] = dto.StaffResponse{
			ID:                        user.ID,
			Email:                     user.Email,
			Phone:                     user.Phone,
			FullName:                  user.FullName,
			Status:                    user.Status,
			Avatar:                    user.Avatar,
			TodayTransactionCount:     todayCount,
			TodayTransactionAmount:    todayAmount,
			WalletBalance:             balance,
			DailyLimitCount:           maxCount,
			DailyLimitAmount:          maxAmount,
			MarginSchemeType:          "",
			MarginValue:               0,
			CreatedAt:                 user.CreatedAt,
		}

		if marginSetting != nil {
			staffResponses[i].MarginSchemeType = marginSetting.SchemeType
			if marginSetting.SchemeType == "MarginShare" {
				staffResponses[i].MarginValue = marginSetting.GlobalMarginPercent
			} else {
				staffResponses[i].MarginValue = marginSetting.FixedAllowance
			}
		}
	}

	return &dto.ListUsersResponse{
		Users:  staffResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
	// ... rest unchanged

	staffResponses := make([]dto.StaffResponse, len(users))
	for i, user := range users {
		balance, _ := s.userRepo.GetWalletBalance(user.ID)
		todayCount, todayAmount, _ := s.userRepo.GetStaffTodayStats(user.ID)
		maxCount, maxAmount, _ := s.userRepo.GetStaffDailyLimit(user.ID)
		marginSetting, _ := s.marginService.GetStaffMarginSettings(user.ID)

		staffResponses[i] = dto.StaffResponse{
			ID:                        user.ID,
			Email:                     user.Email,
			Phone:                     user.Phone,
			FullName:                  user.FullName,
			Status:                    user.Status,
			Avatar:                    user.Avatar,
			TodayTransactionCount:     todayCount,
			TodayTransactionAmount:    todayAmount,
			WalletBalance:             balance,
			DailyLimitCount:           maxCount,
			DailyLimitAmount:          maxAmount,
			MarginSchemeType:          "",
			MarginValue:               0,
			CreatedAt:                 user.CreatedAt,
		}

		if marginSetting != nil {
			staffResponses[i].MarginSchemeType = marginSetting.SchemeType
			if marginSetting.SchemeType == "MarginShare" {
				staffResponses[i].MarginValue = marginSetting.GlobalMarginPercent
			} else {
				staffResponses[i].MarginValue = marginSetting.FixedAllowance
			}
		}
	}

	return &dto.ListUsersResponse{
		Users:  staffResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetStaff returns detailed staff info
func (s *UserService) GetStaff(ctx context.Context, staffID uint) (*dto.StaffDetailResponse, error) {
	user, err := s.userRepo.FindStaffByID(staffID)
	if err != nil {
		return nil, err
	}

	balance, _ := s.userRepo.GetWalletBalance(user.ID)
	todayCount, todayAmount, _ := s.userRepo.GetStaffTodayStats(user.ID)
	maxCount, maxAmount, _ := s.userRepo.GetStaffDailyLimit(user.ID)

	marginSetting, err := s.marginService.GetStaffMarginSettings(user.ID)
	if err != nil && err.Error() != "no margin setting found for staff" {
		// log error but continue
	}

	resp := &dto.StaffDetailResponse{
		ID:                        user.ID,
		Email:                     user.Email,
		Phone:                     user.Phone,
		FullName:                  user.FullName,
		Status:                    user.Status,
		Avatar:                    user.Avatar,
		Address:                   user.Address,
		DateOfBirth:               user.DateOfBirth,
		TodayTransactionCount:     todayCount,
		TodayTransactionAmount:    todayAmount,
		WalletBalance:             balance,
		DailyLimit: struct {
			Count    int     `json:"count"`
			MaxCount int     `json:"max_count"`
			Amount   float64 `json:"amount"`
			MaxAmount float64 `json:"max_amount"`
		}{
			Count:     0,
			MaxCount:  maxCount,
			Amount:    0,
			MaxAmount: maxAmount,
		},
		MarginSetting: struct {
			SchemeType            string  `json:"scheme_type"`
			GlobalMarginPercent   float64 `json:"global_margin_percent"`
			FixedAllowance        float64 `json:"fixed_allowance"`
		}{
			SchemeType:          "",
			GlobalMarginPercent: 0,
			FixedAllowance:      0,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if marginSetting != nil {
		resp.MarginSetting.SchemeType = marginSetting.SchemeType
		resp.MarginSetting.GlobalMarginPercent = marginSetting.GlobalMarginPercent
		resp.MarginSetting.FixedAllowance = marginSetting.FixedAllowance
	}

	return resp, nil
}

// CreateStaff creates a new staff user (Mitra only)
func (s *UserService) CreateStaff(ctx context.Context, mitraID uint, req *dto.StaffCreateRequest) (*dto.StaffDetailResponse, error) {
	// Check if email exists
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}
	existing, _ = s.userRepo.FindByPhone(req.Phone)
	if existing != nil {
		return nil, errors.New("phone already exists")
	}

	user := &models.User{
		Email:         req.Email,
		Phone:         req.Phone,
		FullName:      req.FullName,
		Role:          "staff",
		Status:        "active",
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign staff role
	role, err := s.roleRepo.FindByName("staff")
	if err != nil {
		return nil, err
	}
	if err := s.roleRepo.AssignRole(user.ID, role.ID); err != nil {
		return nil, err
	}

	// Set daily limit if provided
	if req.DailyLimitCount != nil || req.DailyLimitAmount != nil {
		limit := &models.DailyLimit{
			UserID:    user.ID,
			Date:      time.Now().Format("2006-01-02"),
			MaxCount:  100,
			MaxAmount: 10000000,
		}
		if req.DailyLimitCount != nil {
			limit.MaxCount = *req.DailyLimitCount
		}
		if req.DailyLimitAmount != nil {
			limit.MaxAmount = *req.DailyLimitAmount
		}
		s.db.Create(limit)
	}

	// Set margin with Mitra ownership
	if req.MarginSchemeType != "" {
		setting := &models.StaffGlobalMarginSetting{
			MitraID:       mitraID,
			StaffID:      user.ID,
			SchemeType:    req.MarginSchemeType,
			IsActive:      true,
		}
		if req.MarginSchemeType == "MarginShare" {
			setting.GlobalMarginPercent = req.MarginValue
		} else {
			setting.FixedAllowance = req.MarginValue
		}
		s.db.Create(setting)
	}

	s.invalidateUserCache(ctx, user.ID)

	return s.GetStaff(ctx, user.ID)
}
	existing, _ = s.userRepo.FindByPhone(req.Phone)
	if existing != nil {
		return nil, errors.New("phone already exists")
	}

	user := &models.User{
		Email:         req.Email,
		Phone:         req.Phone,
		FullName:      req.FullName,
		Role:          "staff",
		Status:        "active",
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByName("staff")
	if err != nil {
		return nil, err
	}
	if err := s.roleRepo.AssignRole(user.ID, role.ID); err != nil {
		return nil, err
	}

	// Daily limit
	if req.DailyLimitCount != nil || req.DailyLimitAmount != nil {
		limit := &models.DailyLimit{
			UserID:   user.ID,
			Date:     time.Now().Format("2006-01-02"),
			MaxCount: 100,
			MaxAmount: 10000000,
		}
		if req.DailyLimitCount != nil {
			limit.MaxCount = *req.DailyLimitCount
		}
		if req.DailyLimitAmount != nil {
			limit.MaxAmount = *req.DailyLimitAmount
		}
		s.db.Create(limit)
	}

	// Margin setting
	if req.MarginSchemeType != "" {
		setting := &models.StaffGlobalMarginSetting{
			MitraID:      0,
			StaffID:      user.ID,
			SchemeType:    req.MarginSchemeType,
			IsActive:      true,
		}
		if req.MarginSchemeType == "MarginShare" {
			setting.GlobalMarginPercent = req.MarginValue
		} else {
			setting.FixedAllowance = req.MarginValue
		}
		s.db.Create(setting)
	}

	s.invalidateUserCache(ctx, user.ID)
	return s.GetStaff(ctx, user.ID)
}

// UpdateStaff updates staff details
func (s *UserService) UpdateStaff(ctx context.Context, staffID uint, req *dto.StaffUpdateRequest) (*dto.StaffDetailResponse, error) {
	user, err := s.userRepo.FindStaffByID(staffID)
	if err != nil {
		return nil, err
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
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Update daily limit if provided
	if req.DailyLimitCount != nil || req.DailyLimitAmount != nil {
		var limit models.DailyLimit
		err := s.db.Where("user_id = ? AND date = ?", staffID, time.Now().Format("2006-01-02")).First(&limit).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			limit = models.DailyLimit{
				UserID: staffID,
				Date:   time.Now().Format("2006-01-02"),
			}
		}
		if req.DailyLimitCount != nil {
			limit.MaxCount = *req.DailyLimitCount
		}
		if req.DailyLimitAmount != nil {
			limit.MaxAmount = *req.DailyLimitAmount
		}
		if limit.ID == 0 {
			s.db.Create(&limit)
		} else {
			s.db.Save(&limit)
		}
	}

	// Update margin if provided
	if req.MarginSchemeType != "" && req.MarginValue > 0 {
		setting := &models.StaffGlobalMarginSetting{}
		err := s.db.Where("staff_id = ? AND is_active = ?", staffID, true).First(setting).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = &models.StaffGlobalMarginSetting{
				MitraID:      0,
				StaffID:      staffID,
				SchemeType:    req.MarginSchemeType,
				IsActive:      true,
			}
			if req.MarginSchemeType == "MarginShare" {
				setting.GlobalMarginPercent = req.MarginValue
			} else {
				setting.FixedAllowance = req.MarginValue
			}
			s.db.Create(setting)
		} else {
			setting.SchemeType = req.MarginSchemeType
			setting.IsActive = true
			if req.MarginSchemeType == "MarginShare" {
				setting.GlobalMarginPercent = req.MarginValue
				setting.FixedAllowance = 0
			} else {
				setting.FixedAllowance = req.MarginValue
				setting.GlobalMarginPercent = 0
			}
			s.db.Save(setting)
		}
	}

	s.invalidateUserCache(ctx, staffID)
	return s.GetStaff(ctx, staffID)
}

// GetStaffStats returns staff stats (alias for GetStaff for now)
func (s *UserService) GetStaffStats(ctx context.Context, staffID uint) (*dto.StaffDetailResponse, error) {
	return s.GetStaff(ctx, staffID)
}

func (s *UserService) GetPendingStaffCount(ctx context.Context) (int64, error) {
	// Count staff with status = 'pending' (invited but not yet active)
	var count int64
	err := s.db.Model(&models.User{}).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ? AND users.status = ?", "staff", "pending").
		Count(&count).Error
	return count, err
}
