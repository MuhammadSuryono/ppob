package services

import (
	"context"
	"testing"

	"github.com/yontech/ppob/transaction-service/config"
	"github.com/yontech/ppob/transaction-service/internal/dto"
	"github.com/yontech/ppob/transaction-service/internal/models"
	"github.com/yontech/ppob/transaction-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTransactionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.Transaction{})
	return db
}

func setupTransactionTestConfig() *config.Config {
	return &config.Config{
		JWTSecret:     "test-secret-key",
		ServerPort:    "8080",
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
	marginRepo := repository.NewMarginRepository(db)

	transactionService := NewTransactionService(transactionRepo, marginRepo, nil, cfg, db)

	req := &dto.CreateTransactionRequest{
		ProductCode:    "PREPAID_SIMPATIS_5K",
		CustomerNumber: "081234567890",
		Amount:         5500,
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

func TestTransactionService_InitiateWithIdempotency(t *testing.T) {
	t.Skip("requires Redis connection")
}

func TestTransactionService_UpdateStatus_ValidTransition(t *testing.T) {
	db := setupTransactionTestDB(t)
	cfg := setupTransactionTestConfig()

	transactionRepo := repository.NewTransactionRepository(db)
	marginRepo := repository.NewMarginRepository(db)

	transactionService := NewTransactionService(transactionRepo, marginRepo, nil, cfg, db)

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
	marginRepo := repository.NewMarginRepository(db)

	transactionService := NewTransactionService(transactionRepo, marginRepo, nil, cfg, db)

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
	marginRepo := repository.NewMarginRepository(db)

	transactionService := NewTransactionService(transactionRepo, marginRepo, nil, cfg, db)

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