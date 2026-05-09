package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "ppob"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := getEnvVal(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvVal(key string) string {
	return ""
}

type InvariantCheckResult struct {
	CheckName   string
	Status     string
	AffectedRows int64
	Message    string
}

func main() {
	cfg := Load()

	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	reconciliationService := NewReconciliationService(db, redisClient)

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	log.Println("Starting invariant monitoring job...")
	
	for {
		log.Println("Running invariant checks...")
		results := reconciliationService.RunInvariantChecks(context.Background())
		
		for _, result := range results {
			if result.Status == "CRITICAL" {
				log.Printf("CRITICAL: %s - %s (affected: %d)", result.CheckName, result.Message, result.AffectedRows)
			} else if result.Status == "WARNING" {
				log.Printf("WARNING: %s - %s (affected: %d)", result.CheckName, result.Message, result.AffectedRows)
			}
		}
		
		log.Println("Invariant checks completed")
		<-ticker.C
	}
}

type ReconciliationService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewReconciliationService(db *gorm.DB, redis *redis.Client) *ReconciliationService {
	return &ReconciliationService{db: db, redis: redis}
}

func (s *ReconciliationService) RunInvariantChecks(ctx context.Context) []InvariantCheckResult {
	var results []InvariantCheckResult

	results = append(results, s.checkWalletBalanceConsistency(ctx)...)
	results = append(results, s.checkNegativeBalances(ctx)...)
	results = append(results, s.checkDuplicateRefIDs(ctx)...)
	results = append(results, s.checkOrphanedTransactions(ctx)...)
	results = append(results, s.checkOrphanedWallets(ctx)...)
	results = append(results, s.checkInvalidTransactionStatus(ctx)...)

	return results
}

func (s *ReconciliationService) checkWalletBalanceConsistency(ctx context.Context) []InvariantCheckResult {
	type WalletBalance struct {
		WalletID          string  `gorm:"wallet_id"`
		BalanceAvailable  float64 `gorm:"balance_available"`
		BalanceHeld       float64 `gorm:"balance_held"`
		ComputedBalance   float64 `gorm:"computed_balance"`
	}

	var results []InvariantCheckResult

	var mismatches []WalletBalance
	s.db.Raw(`
		SELECT 
			w.wallet_id,
			w.balance_available,
			w.balance_held,
			COALESCE(SUM(
				CASE 
					WHEN we.event_type IN ('Credited', 'TopupAdded', 'Refunded', 'HoldReleased') THEN we.amount
					WHEN we.event_type IN ('Debited', 'Held') THEN -we.amount
					ELSE 0
				END
			), 0) AS computed_balance
		FROM wallets w
		LEFT JOIN wallet_events we ON w.wallet_id = we.wallet_id
		GROUP BY w.wallet_id, w.balance_available, w.balance_held
		HABSING ABS((w.balance_available + w.balance_held) - computed_balance) > 1
	`).Scan(&mismatches)

	if len(mismatches) > 0 {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 1: Wallet Balance Consistency",
			Status:       "CRITICAL",
			AffectedRows: int64(len(mismatches)),
			Message:      "Wallet balance does not match event log sum",
		})
	} else {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 1: Wallet Balance Consistency",
			Status:       "OK",
			AffectedRows: 0,
			Message:      "All wallet balances are consistent",
		})
	}

	return results
}

func (s *ReconciliationService) checkNegativeBalances(ctx context.Context) []InvariantCheckResult {
	type NegativeWallet struct {
		WalletID         string  `gorm:"wallet_id"`
		BalanceAvailable float64 `gorm:"balance_available"`
	}

	var results []InvariantCheckResult
	var negatives []NegativeWallet

	s.db.Raw("SELECT wallet_id, balance_available FROM wallets WHERE balance_available < 0").Scan(&negatives)

	if len(negatives) > 0 {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 2: No Negative Available Balance",
			Status:       "CRITICAL",
			AffectedRows: int64(len(negatives)),
			Message:      "Found wallets with negative balance",
		})
	} else {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 2: No Negative Available Balance",
			Status:       "OK",
			AffectedRows: 0,
			Message:      "All wallets have non-negative balance",
		})
	}

	return results
}

func (s *ReconciliationService) checkDuplicateRefIDs(ctx context.Context) []InvariantCheckResult {
	type DuplicateRef struct {
		RefID   string `gorm:"ref_id"`
		Count   int64  `gorm:"count"`
	}

	var results []InvariantCheckResult
	var duplicates []DuplicateRef

	s.db.Raw("SELECT ref_id, COUNT(*) as count FROM transactions GROUP BY ref_id HAVING COUNT(*) > 1").Scan(&duplicates)

	if len(duplicates) > 0 {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 6: Reference ID Uniqueness",
			Status:       "CRITICAL",
			AffectedRows: int64(len(duplicates)),
			Message:      "Found duplicate reference IDs in transactions",
		})
	} else {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 6: Reference ID Uniqueness",
			Status:       "OK",
			AffectedRows: 0,
			Message:      "All transaction reference IDs are unique",
		})
	}

	return results
}

func (s *ReconciliationService) checkOrphanedTransactions(ctx context.Context) []InvariantCheckResult {
	var results []InvariantCheckResult
	var count int64

	s.db.Raw("SELECT COUNT(*) FROM transactions t LEFT JOIN wallets w ON t.wallet_id = w.wallet_id WHERE w.wallet_id IS NULL").Scan(&count)

	if count > 0 {
		results = append(results, InvariantCheckResult{
			CheckName:    "Orphaned Transactions",
			Status:       "CRITICAL",
			AffectedRows: count,
			Message:      "Found transactions without valid wallet",
		})
	} else {
		results = append(results, InvariantCheckResult{
			CheckName:    "Orphaned Transactions",
			Status:       "OK",
			AffectedRows: 0,
			Message:      "All transactions have valid wallets",
		})
	}

	return results
}

func (s *ReconciliationService) checkOrphanedWallets(ctx context.Context) []InvariantCheckResult {
	var results []InvariantCheckResult
	var count int64

	s.db.Raw("SELECT COUNT(*) FROM users u LEFT JOIN wallets w ON u.user_id = w.owner_id WHERE w.wallet_id IS NULL AND u.is_active = true").Scan(&count)

	if count > 0 {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 7: Wallet Ownership",
			Status:       "WARNING",
			AffectedRows: count,
			Message:      "Found active users without wallet",
		})
	} else {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 7: Wallet Ownership",
			Status:       "OK",
			AffectedRows: 0,
			Message:      "All active users have wallets",
		})
	}

	return results
}

func (s *ReconciliationService) checkInvalidTransactionStatus(ctx context.Context) []InvariantCheckResult {
	var results []InvariantCheckResult
	var count int64

	s.db.Raw("SELECT COUNT(*) FROM transactions WHERE status NOT IN ('Initiated', 'Pending', 'Success', 'Failed', 'Expired', 'Cancelled', 'Refunded')").Scan(&count)

	if count > 0 {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 8: All Transaction Statuses are Valid",
			Status:       "CRITICAL",
			AffectedRows: count,
			Message:      "Found transactions with invalid status",
		})
	} else {
		results = append(results, InvariantCheckResult{
			CheckName:    "Invariant 8: All Transaction Statuses are Valid",
			Status:       "OK",
			AffectedRows: 0,
			Message:      "All transaction statuses are valid",
		})
	}

	return results
}