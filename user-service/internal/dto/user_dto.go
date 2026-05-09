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
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type StaffResponse struct {
	ID                          uint      `json:"id"`
	Email                       string    `json:"email"`
	Phone                       string    `json:"phone"`
	FullName                    string    `json:"full_name"`
	Status                      string    `json:"status"`
	Avatar                      string    `json:"avatar"`
	TodayTransactionCount       int       `json:"today_transaction_count"`
	TodayTransactionAmount      float64   `json:"today_transaction_amount"`
	WalletBalance               float64   `json:"wallet_balance"`
	DailyLimitCount             int       `json:"daily_limit_count"`
	DailyLimitAmount            float64   `json:"daily_limit_amount"`
	MarginSchemeType            string    `json:"margin_scheme_type"`
	MarginValue                 float64   `json:"margin_value"`
	CreatedAt                   time.Time `json:"created_at"`
}

type StaffDetailResponse struct {
	ID                          uint      `json:"id"`
	Email                       string    `json:"email"`
	Phone                       string    `json:"phone"`
	FullName                    string    `json:"full_name"`
	Status                      string    `json:"status"`
	Avatar                      string    `json:"avatar"`
	Address                     string    `json:"address"`
	DateOfBirth                 *time.Time `json:"date_of_birth"`
	TodayTransactionCount       int       `json:"today_transaction_count"`
	TodayTransactionAmount      float64   `json:"today_transaction_amount"`
	WalletBalance               float64   `json:"wallet_balance"`
	DailyLimit                  struct {
		Count    int     `json:"count"`
		MaxCount int     `json:"max_count"`
		Amount   float64 `json:"amount"`
		MaxAmount float64 `json:"max_amount"`
	} `json:"daily_limit"`
	MarginSetting struct {
		SchemeType            string  `json:"scheme_type"`
		GlobalMarginPercent   float64 `json:"global_margin_percent"`
		FixedAllowance        float64 `json:"fixed_allowance"`
	} `json:"margin_setting"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

type StaffCreateRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"required"`
	FullName     string `json:"full_name" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	Pin          string `json:"pin" binding:"required,min=6"`
	DailyLimitCount    *int     `json:"daily_limit_count"`
	DailyLimitAmount   *float64 `json:"daily_limit_amount"`
	MarginSchemeType   string   `json:"margin_scheme_type"` // FixedAllowance or MarginShare
	MarginValue        float64  `json:"margin_value"`
}

type StaffUpdateRequest struct {
	FullName              string  `json:"full_name"`
	Phone                 string  `json:"phone"`
	Avatar                string  `json:"avatar"`
	Address               string  `json:"address"`
	DateOfBirth           *time.Time `json:"date_of_birth"`
	Status                string  `json:"status"` // active/inactive
	DailyLimitCount       *int    `json:"daily_limit_count"`
	DailyLimitAmount      *float64 `json:"daily_limit_amount"`
	MarginSchemeType      string  `json:"margin_scheme_type"`
	MarginValue           float64 `json:"margin_value"`
}

type StaffListQuery struct {
	Limit  int    `form:"limit" binding:"min=1,max=100"`
	Offset int    `form:"offset" binding:"min=0"`
	Search string `form:"search"`
	Status string `form:"status"` // active/inactive/all
}

type ReportKPIResponse struct {
	TotalSales          float64 `json:"total_sales"`
	PlatformProfit      float64 `json:"platform_profit"`
	StaffCount          int     `json:"staff_count"`
	SuccessRate         float64 `json:"success_rate"`
	TransactionCount    int     `json:"transaction_count"`
	PeriodStart         string  `json:"period_start"`
	PeriodEnd           string  `json:"period_end"`
}

type ReportSalesTrendResponse struct {
	Date    string  `json:"date"`
	Sales   float64 `json:"sales"`
	Count   int     `json:"count"`
}

type ReportStaffPerformanceResponse struct {
	StaffID       uint    `json:"staff_id"`
	StaffName     string  `json:"staff_name"`
	TransactionCount int `json:"transaction_count"`
	TotalSales    float64 `json:"total_sales"`
	TotalCommission float64 `json:"total_commission"`
	SuccessRate   float64 `json:"success_rate"`
}

type ReportsResponse struct {
	KPIs            []ReportKPIResponse             `json:"kpis"`
	SalesTrend     []ReportSalesTrendResponse     `json:"sales_trend"`
	StaffPerformance []ReportStaffPerformanceResponse `json:"staff_performance"`
}

type NotificationResponse struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type UnreadCountResponse struct {
	UnreadCount int `json:"unread_count"`
}