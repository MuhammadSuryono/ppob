package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/auth-service/config"
	"github.com/yontech/ppob/auth-service/internal/handlers"
	"github.com/yontech/ppob/auth-service/internal/middleware"
	"github.com/yontech/ppob/auth-service/internal/repository"
	"github.com/yontech/ppob/auth-service/internal/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func init() {
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	testDB.AutoMigrate(
		&User{},
		&OTP{},
		&Wallet{},
	)
}

type User struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Email    string    `gorm:"uniqueIndex" json:"email"`
	Phone    string    `gorm:"uniqueIndex" json:"phone"`
	Password string    `json:"-"`
	FullName string    `json:"full_name"`
	Role     string    `gorm:"default:user" json:"role"`
	Status   string    `gorm:"default:active" json:"status"`
}

type OTP struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
}

type Wallet struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	UserID  uint    `gorm:"uniqueIndex" json:"user_id"`
	Balance float64 `gorm:"default:0" json:"balance"`
}

type TestAuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func (h *TestAuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		PIN      string `json:"pin" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser := User{}
	if err := h.db.Where("email = ? OR phone = ?", req.Email, req.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	user := User{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		FullName: req.FullName,
		Role:     "user",
		Status:   "active",
	}
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	wallet := Wallet{UserID: user.ID, Balance: 0}
	h.db.Create(&wallet)

	otp := OTP{
		UserID:    user.ID,
		Code:      "123456",
		Type:      "register",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	h.db.Create(&otp)

	token := fmt.Sprintf("test-token-%d", user.ID)
	c.JSON(http.StatusCreated, gin.H{
		"user_id":    user.ID,
		"email":      user.Email,
		"phone":      user.Phone,
		"full_name":  user.FullName,
		"token":      token,
		"expires_at": time.Now().Add(15 * time.Minute).Unix(),
	})
}

func (h *TestAuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required_without=Phone"`
		Phone    string `json:"phone" binding:"required_without=Email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	query := "email = ?"
	if req.Email != "" {
		query = "email = ?"
	} else {
		query = "phone = ?"
	}

	if err := h.db.Where(query, req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := fmt.Sprintf("test-token-%d", user.ID)
	refreshToken := fmt.Sprintf("refresh-token-%d", user.ID)

	c.JSON(http.StatusOK, gin.H{
		"user_id":        user.ID,
		"email":          user.Email,
		"phone":          user.Phone,
		"full_name":       user.FullName,
		"token":          token,
		"refresh_token":  refreshToken,
		"expires_at":     time.Now().Add(15 * time.Minute).Unix(),
	})
}

func TestAuthUserIntegration_RegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Load()
	r := gin.New()

	handler := &TestAuthHandler{db: testDB, cfg: cfg}

	r.POST("/api/v1/auth/register", handler.Register)
	r.POST("/api/v1/auth/login", handler.Login)

	registerReq := map[string]interface{}{
		"email":     "integration@test.com",
		"phone":     "081234567890",
		"full_name": "Integration Test",
		"password":  "test123456",
		"pin":       "123456",
	}
	registerBody, _ := json.Marshal(registerReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(string(registerBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var registerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResp)

	userID := registerResp["user_id"].(float64)

	if userID == 0 {
		t.Error("Expected non-zero user ID")
	}

	loginReq := map[string]interface{}{
		"email":    "integration@test.com",
		"password": "test123456",
	}
	loginBody, _ := json.Marshal(loginReq)

	loginReqHTTP := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(string(loginBody)))
	loginReqHTTP.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReqHTTP)

	if loginW.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", loginW.Code, loginW.Body.String())
	}

	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)

	if loginResp["token"] == nil {
		t.Error("Expected non-nil token")
	}

	if loginResp["refresh_token"] == nil {
		t.Error("Expected non-nil refresh token")
	}
}

func TestAuthUserIntegration_MultipleRoleSwitching(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Load()
	r := gin.New()

	handler := &TestAuthHandler{db: testDB, cfg: cfg}

	r.POST("/api/v1/auth/register", handler.Register)
	r.POST("/api/v1/auth/login", handler.Login)

	registerReq := map[string]interface{}{
		"email":     "role@test.com",
		"phone":     "081234567891",
		"full_name": "Role Test",
		"password":  "test123456",
		"pin":       "123456",
	}
	registerBody, _ := json.Marshal(registerReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(string(registerBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", w.Code)
	}

	var registerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResp)

	userID := uint(registerResp["user_id"].(float64))

	var user User
	testDB.First(&user, userID)

	roles := []string{"user", "mitra", "reseller", "admin"}

	for _, role := range roles {
		user.Role = role
		testDB.Save(&user)

		loginReq := map[string]interface{}{
			"email":    "role@test.com",
			"password": "test123456",
		}
		loginBody, _ := json.Marshal(loginReq)

		loginReqHTTP := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(string(loginBody)))
		loginReqHTTP.Header.Set("Content-Type", "application/json")

		loginW := httptest.NewRecorder()
		r.ServeHTTP(loginW, loginReqHTTP)

		if loginW.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for role %s, got %d", role, loginW.Code)
		}
	}
}

func TestAuthUserIntegration_InvalidLoginAttempts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Load()
	r := gin.New()

	handler := &TestAuthHandler{db: testDB, cfg: cfg}

	r.POST("/api/v1/auth/login", handler.Login)

	invalidLogins := []struct {
		name     string
		email    string
		password string
	}{
		{"wrong email", "wrong@test.com", "password123"},
		{"wrong password", "test@test.com", "wrongpassword"},
		{"empty email", "", "password123"},
		{"empty password", "test@test.com", ""},
	}

	for _, tc := range invalidLogins {
		t.Run(tc.name, func(t *testing.T) {
			loginReq := map[string]interface{}{
				"email":    tc.email,
				"password": tc.password,
			}
			loginBody, _ := json.Marshal(loginReq)

			loginReqHTTP := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(string(loginBody)))
			loginReqHTTP.Header.Set("Content-Type", "application/json")

			loginW := httptest.NewRecorder()
			r.ServeHTTP(loginW, loginReqHTTP)

			if loginW.Code != http.StatusUnauthorized {
				t.Errorf("Expected status 401, got %d", loginW.Code)
			}
		})
	}
}

func TestAuthUserIntegration_DuplicateRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Load()
	r := gin.New()

	handler := &TestAuthHandler{db: testDB, cfg: cfg}

	r.POST("/api/v1/auth/register", handler.Register)

	registerReq := map[string]interface{}{
		"email":     "duplicate@test.com",
		"phone":     "081234567892",
		"full_name": "Duplicate Test",
		"password":  "test123456",
		"pin":       "123456",
	}
	registerBody, _ := json.Marshal(registerReq)

	req1 := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(string(registerBody)))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 for first registration, got %d", w1.Code)
	}

	req2 := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(string(registerBody)))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate registration, got %d", w2.Code)
	}
}

func TestAuthUserIntegration_PasswordChange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Load()
	r := gin.New()

	handler := &TestAuthHandler{db: testDB, cfg: cfg}

	r.POST("/api/v1/auth/register", handler.Register)
	r.POST("/api/v1/auth/login", handler.Login)

	registerReq := map[string]interface{}{
		"email":     "password@test.com",
		"phone":     "081234567893",
		"full_name": "Password Test",
		"password":  "oldpassword",
		"pin":       "123456",
	}
	registerBody, _ := json.Marshal(registerReq)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(string(registerBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	loginOld := map[string]interface{}{
		"email":    "password@test.com",
		"password": "oldpassword",
	}
	loginOldBody, _ := json.Marshal(loginOld)

	loginOldReq := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(string(loginOldBody)))
	loginOldReq.Header.Set("Content-Type", "application/json")

	loginOldW := httptest.NewRecorder()
	r.ServeHTTP(loginOldW, loginOldReq)

	if loginOldW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for old password login, got %d", loginOldW.Code)
	}

	var user User
	testDB.Where("email = ?", "password@test.com").First(&user)
	user.Password = "newpassword"
	testDB.Save(&user)

	loginNew := map[string]interface{}{
		"email":    "password@test.com",
		"password": "newpassword",
	}
	loginNewBody, _ := json.Marshal(loginNew)

	loginNewReq := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(string(loginNewBody)))
	loginNewReq.Header.Set("Content-Type", "application/json")

	loginNewW := httptest.NewRecorder()
	r.ServeHTTP(loginNewW, loginNewReq)

	if loginNewW.Code != http.StatusOK {
		t.Errorf("Expected status 200 for new password login, got %d", loginNewW.Code)
	}
}