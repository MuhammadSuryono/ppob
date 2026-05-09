package services

import (
	"context"
	"testing"

	"github.com/yontech/ppob/wallet-service/config"
	"github.com/yontech/ppob/wallet-service/internal/dto"
	"github.com/yontech/ppob/wallet-service/internal/models"
	"github.com/yontech/ppob/wallet-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWalletTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.Wallet{}, &models.WalletEvent{}, &models.Hold{})
	return db
}

func setupWalletTestConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret-key",
		ServerPort: "8080",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "test",
		RedisHost:  "localhost",
		RedisPort:  "6379",
		GinMode:    "test",
	}
}

func TestWalletService_GetBalance(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		HoldAmount: 20000,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	resp, err := walletService.GetBalance(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Balance != 100000 {
		t.Errorf("Expected balance 100000, got %f", resp.Balance)
	}

	if resp.Available != 80000 {
		t.Errorf("Expected available 80000, got %f", resp.Available)
	}
}

func TestWalletService_Credit(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	req := &dto.CreditRequest{
		Amount:        50000,
		ReferenceID:   "test-ref-123",
		ReferenceType: "topup",
	}

	err := walletService.Credit(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var updatedWallet models.Wallet
	db.First(&updatedWallet, 1)

	if updatedWallet.Balance != 150000 {
		t.Errorf("Expected balance 150000, got %f", updatedWallet.Balance)
	}
}

func TestWalletService_Debit(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		HoldAmount: 0,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	req := &dto.DebitRequest{
		Amount:        30000,
		ReferenceID:   "test-ref-123",
		ReferenceType: "purchase",
	}

	err := walletService.Debit(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var updatedWallet models.Wallet
	db.First(&updatedWallet, 1)

	if updatedWallet.Balance != 70000 {
		t.Errorf("Expected balance 70000, got %f", updatedWallet.Balance)
	}
}

func TestWalletService_Debit_InsufficientFunds(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		HoldAmount: 0,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	req := &dto.DebitRequest{
		Amount:        150000,
		ReferenceID:   "test-ref-123",
		ReferenceType: "purchase",
	}

	err := walletService.Debit(context.Background(), 1, req)
	if err != ErrInsufficientFund {
		t.Fatalf("Expected ErrInsufficientFund, got %v", err)
	}
}

func TestWalletService_PlaceHold(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		HoldAmount: 0,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	req := &dto.HoldRequest{
		Amount:        50000,
		ReferenceID:   "order-123",
		ReferenceType: "transaction",
	}

	resp, err := walletService.PlaceHold(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Amount != 50000 {
		t.Errorf("Expected hold amount 50000, got %f", resp.Amount)
	}

	var updatedWallet models.Wallet
	db.First(&updatedWallet, 1)

	if updatedWallet.HoldAmount != 50000 {
		t.Errorf("Expected hold amount 50000, got %f", updatedWallet.HoldAmount)
	}
}

func TestWalletService_ReleaseHold(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		HoldAmount: 50000,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	hold := &models.Hold{
		WalletID:      wallet.ID,
		Amount:        50000,
		ReferenceID:   "order-123",
		ReferenceType: "transaction",
		Status:        "active",
	}
	db.Create(hold)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	req := &dto.ReleaseHoldRequest{
		ReferenceID:   "order-123",
		ReferenceType: "transaction",
	}

	err := walletService.ReleaseHold(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var releasedHold models.Hold
	db.First(&releasedHold, hold.ID)

	if releasedHold.Status != "released" {
		t.Errorf("Expected status released, got %s", releasedHold.Status)
	}
}

func TestWalletService_Transfer(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet1 := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet1)

	wallet2 := &models.Wallet{
		UserID:     2,
		Balance:    50000,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet2)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	err := walletService.Transfer(context.Background(), 1, 2, 30000, "transfer-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var fromWallet, toWallet models.Wallet
	db.First(&fromWallet, wallet1.ID)
	db.First(&toWallet, wallet2.ID)

	if fromWallet.Balance != 70000 {
		t.Errorf("Expected from wallet balance 70000, got %f", fromWallet.Balance)
	}

	if toWallet.Balance != 80000 {
		t.Errorf("Expected to wallet balance 80000, got %f", toWallet.Balance)
	}
}

func TestWalletService_EventSourcing(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	walletService.Credit(context.Background(), 1, &dto.CreditRequest{Amount: 50000, ReferenceID: "ref1", ReferenceType: "topup"})
	walletService.Debit(context.Background(), 1, &dto.DebitRequest{Amount: 20000, ReferenceID: "ref2", ReferenceType: "purchase"})
	walletService.Credit(context.Background(), 1, &dto.CreditRequest{Amount: 10000, ReferenceID: "ref3", ReferenceType: "commission"})

	var events []models.WalletEvent
	db.Where("wallet_id = ?", wallet.ID).Order("created_at asc").Find(&events)

	if len(events) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(events))
	}

	if events[0].EventType != "credit" || events[0].Amount != 50000 {
		t.Errorf("First event should be credit of 50000")
	}

	if events[1].EventType != "debit" || events[1].Amount != 20000 {
		t.Errorf("Second event should be debit of 20000")
	}

	if events[2].EventType != "credit" || events[2].Amount != 10000 {
		t.Errorf("Third event should be credit of 10000")
	}
}

func TestWalletService_Reconciliation(t *testing.T) {
	db := setupWalletTestDB(t)
	cfg := setupWalletTestConfig()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	wallet := &models.Wallet{
		UserID:     1,
		Balance:    100000,
		Currency:   "IDR",
		Status:     "active",
	}
	db.Create(wallet)

	db.Create(&models.WalletEvent{
		WalletID:      wallet.ID,
		EventType:     "credit",
		Amount:        100000,
		BalanceBefore: 0,
		BalanceAfter:  100000,
		ReferenceID:   "initial",
		ReferenceType: "setup",
	})

	walletService := NewWalletService(walletRepo, eventRepo, nil, cfg)

	walletService.Credit(context.Background(), 1, &dto.CreditRequest{Amount: 50000, ReferenceID: "ref1", ReferenceType: "topup"})
	walletService.Debit(context.Background(), 1, &dto.DebitRequest{Amount: 20000, ReferenceID: "ref2", ReferenceType: "purchase"})

	resp, err := walletService.Reconcile(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !resp.IsBalanced {
		t.Errorf("Expected wallet to be balanced, got drift %f", resp.Drift)
	}

	if resp.Drift != 0 {
		t.Errorf("Expected drift 0, got %f", resp.Drift)
	}
}