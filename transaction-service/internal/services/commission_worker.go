package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/transaction-service/internal/repository"
)

const CommissionQueueKey = "q:commission_processing"

type CommissionWorker struct {
	redis             *redis.Client
	transactionRepo   *repository.TransactionRepository
	marginSvc         *MarginService
	commissionSvc     *CommissionService
	stopChan          chan struct{}
}

func NewCommissionWorker(
	redis *redis.Client,
	transactionRepo *repository.TransactionRepository,
	marginSvc *MarginService,
	commissionSvc *CommissionService,
) *CommissionWorker {
	return &CommissionWorker{
		redis:           redis,
		transactionRepo: transactionRepo,
		marginSvc:       marginSvc,
		commissionSvc:   commissionSvc,
		stopChan:        make(chan struct{}),
	}
}

type CommissionTask struct {
	TransactionID string `json:"transaction_id"`
}

func (w *CommissionWorker) Start(ctx context.Context) {
	log.Println("Commission Worker started")
	for {
		select {
		case <-w.stopChan:
			return
		case <-ctx.Done():
			return
		default:
			// BLPop with 5s timeout to allow check for stopChan
			res, err := w.redis.BLPop(ctx, 5*time.Second, CommissionQueueKey).Result()
			if err != nil {
				if err != redis.Nil {
					log.Printf("Worker error: %v", err)
				}
				continue
			}

			var task CommissionTask
			if err := json.Unmarshal([]byte(res[1]), &task); err != nil {
				log.Printf("Failed to unmarshal commission task: %v", err)
				continue
			}

			if err := w.processTask(ctx, task.TransactionID); err != nil {
				log.Printf("Failed to process commission for transaction %s: %v", task.TransactionID, err)
			}
		}
	}
}

func (w *CommissionWorker) Stop() {
	close(w.stopChan)
}

func (w *CommissionWorker) processTask(ctx context.Context, transactionID string) error {
	tx, err := w.transactionRepo.FindByTransactionID(transactionID)
	if err != nil {
		return err
	}

	if tx.Status != "success" {
		return fmt.Errorf("transaction %s is not in success state (current: %s)", transactionID, tx.Status)
	}

	// 1. Calculate margin and commission
	// We need to know if the user is a staff
	// For now, MarginService.CalculateTransactionMargin handles the lookup
	_, commission, err := w.marginSvc.CalculateTransactionMargin(tx.UserID, tx.ProductCode, tx.SellingPrice)
	if err != nil {
		// If no margin setting, it's not an error we should retry forever, maybe they are Mitra
		if err == ErrNoMarginSetting || err == ErrStaffNotFound {
			log.Printf("No commission to process for transaction %s (User ID: %d)", transactionID, tx.UserID)
			return nil
		}
		return err
	}

	if commission <= 0 {
		return nil
	}

	// 2. Create and Credit Commission
	// We need a WalletID. If not in Transaction, we find it.
	var walletID uint
	if tx.WalletID != nil {
		walletID = *tx.WalletID
	} else {
		wallet, err := w.transactionRepo.GetActiveWallet(ctx, tx.UserID)
		if err != nil {
			return fmt.Errorf("failed to find wallet for user %d: %v", tx.UserID, err)
		}
		walletID = wallet.ID
	}

	// Create commission record and credit wallet
	// Using the internal walletSvc implementation
	_, err = w.commissionSvc.CreateAndCreditCommission(ctx, tx.UserID, walletID, transactionID, commission, "MarginShare", 1)
	if err != nil {
		return err
	}

	log.Printf("Commission of %.2f successfully credited to user %d for transaction %s", commission, tx.UserID, transactionID)
	return nil
}
