package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/shared/proto/product"
	"github.com/yontech/ppob/transaction-service/config"
	"github.com/yontech/ppob/transaction-service/internal/clients"
	"github.com/yontech/ppob/transaction-service/internal/dto"
	"github.com/yontech/ppob/transaction-service/internal/models"
	"github.com/yontech/ppob/transaction-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockProductClient struct {
	product *product.GetProductResponse
	err     error
}

func (m *mockProductClient) GetProductByCode(ctx context.Context, skuCode string) (*product.GetProductResponse, error) {
	return m.product, m.err
}

func (m *mockProductClient) ValidateProduct(ctx context.Context, productID uint, expectedPrice float64) (*product.ValidateProductResponse, error) {
	return nil, nil
}

type mockWalletClient struct {
	err error
}

func (m *mockWalletClient) PlaceHoldForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error {
	return m.err
}

func (m *mockWalletClient) ReleaseHoldForTransaction(ctx context.Context, userID uint, transactionID string) error {
	return m.err
}

func (m *mockWalletClient) DebitForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error {
	return m.err
}

func (m *mockWalletClient) CreditWallet(ctx context.Context, userID uint, amount float64, referenceID, referenceType string) error {
	return m.err
}

func (m *mockWalletClient) DebitWallet(ctx context.Context, userID uint, amount float64, referenceID, referenceType string) error {
	return m.err
}

type mockIntegrationClient struct {
	resp *clients.IntegrationResponse
	err  error
}

func (m *mockIntegrationClient) TopUp(ctx context.Context, req *clients.TopUpRequest) (*clients.IntegrationResponse, error) {
	if m.resp == nil && m.err == nil {
		return &clients.IntegrationResponse{Success: true}, nil
	}
	return m.resp, m.err
}

func setupTransactionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.Transaction{}, &models.Product{})
	return db
}

func setupTransactionTestConfig() *config.Config {
	return &config.Config{
		ServerPort: "8080",
		DBHost:        "localhost",
		DBPort:        "5432",
		DBUser:        "postgres",
		DBPassword:    "postgres",
		DBName:        "test",
		RedisHost:     "localhost",
		RedisPort:     "6379",
		GinMode:       "test",
		DigiflazzURL:  "https://api.digiflazz.com/v1",
		DigiflazzKey:  "",
		DigiflazzSecret: "test-secret",
	}
}

func TestTransactionStateMachine_ValidTransitions(t *testing.T) {
	tests := []struct {
		name      string
		fromState string
		toState   string
		valid     bool
	}{
		{"initiated to pending", "initiated", "pending", true},
		{"initiated to success", "initiated", "success", true},
		{"initiated to failed", "initiated", "failed", true},
		{"pending to success", "pending", "success", true},
		{"pending to failed", "pending", "failed", true},
		{"pending to expired", "pending", "expired", true},
		{"pending to cancelled", "pending", "cancelled", true},
		{"success to refunded", "success", "refunded", true},
		{"failed to success", "failed", "success", false},
		{"success to pending", "success", "pending", false},
		{"pending to processing", "pending", "processing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewTransactionStateMachine(tt.fromState)
			newState := TransactionState(tt.toState)
			result := sm.CanTransitionTo(newState)

			if result != tt.valid {
				t.Errorf("Expected transition from %s to %s to be %v, got %v", tt.fromState, tt.toState, tt.valid, result)
			}
		})
	}
}

func TestTransactionStateMachine_IsTerminal(t *testing.T) {
	nonTerminal := []string{"initiated", "pending"}
	for _, state := range nonTerminal {
		sm := NewTransactionStateMachine(state)
		if sm.IsTerminal() {
			t.Errorf("Expected %s to be non-terminal", state)
		}
	}

	terminal := []string{"success", "failed", "expired", "cancelled", "refunded"}
	for _, state := range terminal {
		sm := NewTransactionStateMachine(state)
		if !sm.IsTerminal() {
			t.Errorf("Expected %s to be terminal", state)
		}
	}
}

func TestTransactionService_InitiateTransaction(t *testing.T) {
	db := setupTransactionTestDB(t)
	cfg := setupTransactionTestConfig()

	transactionRepo := repository.NewTransactionRepository(db)
	mockProd := &mockProductClient{
		product: &product.GetProductResponse{
			Id:       1,
			SkuCode:  "PREPAID_SIMPATIS_5K",
			Price:    5500,
			IsActive: true,
		},
	}
	marginService := NewMarginService(db, cfg, mockProd)
	transactionService := NewTransactionService(transactionRepo, marginService, nil, cfg, db, &mockWalletClient{}, mockProd, &mockIntegrationClient{})

	// Add product to test DB
	productModel := &models.Product{
		Code:     "PREPAID_SIMPATIS_5K",
		Name:     "Simpati 5000",
		Price:    5500,
		Status:   "active",
		IsActive: true,
	}
	db.Create(productModel)

	req := &dto.CreateTransactionRequest{
		ProductCode:    "PREPAID_SIMPATIS_5K",
		CustomerNumber: "081234567890",
		Amount:         5500,
		AuthorizeID:    "test-auth-id",
	}

	resp, err := transactionService.InitiateTransaction(context.Background(), 1, req, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.TransactionID == "" {
		t.Error("Expected non-empty transaction ID")
	}

	if resp.Status != "initiated" {
		t.Errorf("Expected status 'initiated', got %s", resp.Status)
	}

	if resp.Amount != 5500 {
		t.Errorf("Expected amount 5500, got %f", resp.Amount)
	}
}

func TestTransactionService_InitiateTransaction_AuthorizeID_Success(t *testing.T) {
	db := setupTransactionTestDB(t)
	cfg := setupTransactionTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	transactionRepo := repository.NewTransactionRepository(db)
	mockProd := &mockProductClient{
		product: &product.GetProductResponse{
			Id:       1,
			SkuCode:  "PREPAID_SIMPATIS_5K",
			Price:    5500,
			IsActive: true,
		},
	}
	marginService := NewMarginService(db, cfg, mockProd)
	transactionService := NewTransactionService(transactionRepo, marginService, redisClient, cfg, db, &mockWalletClient{}, mockProd, &mockIntegrationClient{})

	// Add product to test DB
	productModel := &models.Product{
		Code:     "PREPAID_SIMPATIS_5K",
		Name:     "Simpati 5000",
		Price:    5500,
		Status:   "active",
		IsActive: true,
	}
	db.Create(productModel)

	userID := uint(1)
	authorizeID := "valid-auth-id"
	s.Set(fmt.Sprintf("transaction_authorize:%s", authorizeID), fmt.Sprintf("%d", userID))

	req := &dto.CreateTransactionRequest{
		ProductCode:    "PREPAID_SIMPATIS_5K",
		CustomerNumber: "081234567890",
		Amount:         5500,
		AuthorizeID:    authorizeID,
	}

	resp, err := transactionService.InitiateTransaction(context.Background(), userID, req, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.TransactionID == "" {
		t.Error("Expected non-empty transaction ID")
	}

	// Verify key is consumed
	if s.Exists(fmt.Sprintf("transaction_authorize:%s", authorizeID)) {
		t.Error("Expected authorizeID to be consumed (deleted from redis)")
	}
}

func TestTransactionService_InitiateTransaction_AuthorizeID_Invalid(t *testing.T) {
	db := setupTransactionTestDB(t)
	cfg := setupTransactionTestConfig()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to run miniredis: %v", err)
	}
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	transactionRepo := repository.NewTransactionRepository(db)
	mockProd := &mockProductClient{
		product: &product.GetProductResponse{
			Id:       1,
			SkuCode:  "PREPAID_SIMPATIS_5K",
			Price:    5500,
			IsActive: true,
		},
	}
	marginService := NewMarginService(db, cfg, mockProd)
	transactionService := NewTransactionService(transactionRepo, marginService, redisClient, cfg, db, &mockWalletClient{}, mockProd, &mockIntegrationClient{})

	userID := uint(1)
	authorizeID := "invalid-auth-id"
	// Don't set in redis

	req := &dto.CreateTransactionRequest{
		ProductCode:    "PREPAID_SIMPATIS_5K",
		CustomerNumber: "081234567890",
		Amount:         5500,
		AuthorizeID:    authorizeID,
	}

	_, err = transactionService.InitiateTransaction(context.Background(), userID, req, "")
	if err != ErrUnauthorizedTransaction {
		t.Fatalf("Expected ErrUnauthorizedTransaction, got %v", err)
	}
}

func TestTransactionService_InitiateWithIdempotency(t *testing.T) {
	t.Skip("requires Redis connection")
}

func TestTransactionService_UpdateStatus_ValidTransition(t *testing.T) {
	db := setupTransactionTestDB(t)
	cfg := setupTransactionTestConfig()

	transactionRepo := repository.NewTransactionRepository(db)
	mockProd := &mockProductClient{}
	marginService := NewMarginService(db, cfg, mockProd)
	transactionService := NewTransactionService(transactionRepo, marginService, nil, cfg, db, &mockWalletClient{}, mockProd, &mockIntegrationClient{})

	tx := &models.Transaction{
		TransactionID:  "test-123",
		UserID:          1,
		ProductCode:    "PREPAID_SIMPATIS_5K",
		CustomerNumber: "081234567890",
		Amount:         5500,
		Price:          5500,
		Status:         "initiated",
	}
	db.Create(tx)

	updateReq := &dto.UpdateStatusRequest{
		Status:    "pending",
	}

	_, err := transactionService.UpdateTransactionStatus(context.Background(), tx.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var updated models.Transaction
	db.First(&updated, tx.ID)

	if updated.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", updated.Status)
	}
}

func TestTransactionService_UpdateStatus_InvalidTransition(t *testing.T) {
	db := setupTransactionTestDB(t)
	cfg := setupTransactionTestConfig()

	transactionRepo := repository.NewTransactionRepository(db)
	mockProd := &mockProductClient{}
	marginService := NewMarginService(db, cfg, mockProd)
	transactionService := NewTransactionService(transactionRepo, marginService, nil, cfg, db, &mockWalletClient{}, mockProd, &mockIntegrationClient{})

	tx := &models.Transaction{
		TransactionID:  "test-123",
		UserID:          1,
		ProductCode:    "PREPAID_SIMPATIS_5K",
		CustomerNumber: "081234567890",
		Amount:         5500,
		Price:          5500,
		Status:         "success",
	}
	db.Create(tx)

	updateReq := &dto.UpdateStatusRequest{
		Status: "pending",
	}

	_, err := transactionService.UpdateTransactionStatus(context.Background(), tx.ID, updateReq)
	if err != ErrInvalidState {
		t.Fatalf("Expected ErrInvalidState, got %v", err)
	}
}

func TestTransactionService_ListTransactions(t *testing.T) {
	db := setupTransactionTestDB(t)
	cfg := setupTransactionTestConfig()

	transactionRepo := repository.NewTransactionRepository(db)
	mockProd := &mockProductClient{}
	marginService := NewMarginService(db, cfg, mockProd)
	transactionService := NewTransactionService(transactionRepo, marginService, nil, cfg, db, &mockWalletClient{}, mockProd, &mockIntegrationClient{})

	for i := 0; i < 5; i++ {
		tx := &models.Transaction{
			TransactionID:  "test-123-" + string(rune('0'+i)),
			UserID:          1,
			ProductCode:    "PREPAID_SIMPATIS_5K",
			CustomerNumber: "081234567890",
			Amount:         5500,
			Price:          5500,
			Status:         "success",
		}
		db.Create(tx)
	}

	resp, err := transactionService.ListTransactions(context.Background(), 1, "", "", "", 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Transactions) != 5 {
		t.Errorf("Expected 5 transactions, got %d", len(resp.Transactions))
	}

	if resp.Total != 5 {
		t.Errorf("Expected total 5, got %d", resp.Total)
	}
}