package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/shared/events"
	"github.com/yontech/ppob/transaction-service/config"
	"github.com/yontech/ppob/transaction-service/internal/clients"
	"github.com/yontech/ppob/transaction-service/internal/dto"
	"github.com/yontech/ppob/transaction-service/internal/models"
	"github.com/yontech/ppob/transaction-service/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrInvalidAmount            = errors.New("invalid amount")
	ErrTxProductNotFound        = errors.New("product not found")
	ErrProductInactive          = errors.New("product inactive")
	ErrInvalidState             = errors.New("invalid state transition")
	ErrIdempotencyKeyUsed       = errors.New("idempotency key already used")
	ErrMissingRefID             = errors.New("ref_id required for transition")
	ErrAmountMustBePositive     = errors.New("amount must be greater than 0")
	ErrCannotSucceedAfterExpiry = errors.New("cannot succeed after expiry; should be Expired first")
	ErrRefundWindowClosed       = errors.New("refund window closed (24h)")
	ErrTransactionCancelled     = errors.New("transaction cancelled")
	ErrHoldNotFound             = errors.New("hold not found for this transaction")
	ErrHoldFailed               = errors.New("failed to place hold on wallet")
	ErrUnauthorizedTransaction  = errors.New("unauthorized transaction")
)

type TransactionState string

const (
	StateInitiated  TransactionState = "initiated"
	StatePending    TransactionState = "pending"
	StateProcessing TransactionState = "processing"
	StateSuccess    TransactionState = "success"
	StateFailed     TransactionState = "failed"
	StateExpired    TransactionState = "expired"
	StateCancelled  TransactionState = "cancelled"
	StateRefunded   TransactionState = "refunded"
)

type Transaction = models.Transaction

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	marginService   *MarginService
	redis           *redis.Client
	cfg             *config.Config
	db              *gorm.DB
	webhookSecret   string
	walletClient    *clients.WalletClient
	productClient   *clients.ProductClient
	integrationClient *clients.IntegrationClient
	eventPublisher    *events.EventPublisher
}

func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	marginService *MarginService,
	redis *redis.Client,
	cfg *config.Config,
	db *gorm.DB,
	walletClient *clients.WalletClient,
	productClient *clients.ProductClient,
	integrationClient *clients.IntegrationClient,
) *TransactionService {
	return &TransactionService{
		transactionRepo:   transactionRepo,
		marginService:     marginService,
		redis:             redis,
		cfg:               cfg,
		db:                db,
		webhookSecret:     cfg.DigiflazzWebhookSecret,
		walletClient:      walletClient,
		productClient:     productClient,
		integrationClient: integrationClient,
		eventPublisher:    events.NewEventPublisher(redis),
	}
}

func (s *TransactionService) logStateTransition(ctx context.Context, tx *Transaction, oldStatus, newStatus, reason string) {
	if s.db == nil {
		return
	}

	details, _ := json.Marshal(map[string]interface{}{
		"old_status":   oldStatus,
		"new_status":   newStatus,
		"reason":       reason,
		"provider_ref": tx.ProviderRef,
		"digiflazz_rc": tx.DigiflazzRC,
	})

	s.db.WithContext(ctx).Exec(`
		INSERT INTO audit_logs (user_id, action, resource_type, resource_id, details, severity, created_at)
		VALUES (?, ?, 'transaction', ?, ?::jsonb, 'INFO', NOW())
	`, tx.UserID, "transaction_status_change", tx.ID, details)
}

func (s *TransactionService) InitiateTransaction(ctx context.Context, userID uint, req *dto.CreateTransactionRequest, idempotencyKey string) (*dto.TransactionResponse, error) {
	// Validate AuthorizeID
	if s.redis != nil {
		authKey := fmt.Sprintf("transaction_authorize:%s", req.AuthorizeID)
		storedUserID, err := s.redis.Get(ctx, authKey).Result()
		if err != nil || storedUserID != fmt.Sprintf("%d", userID) {
			return nil, ErrUnauthorizedTransaction
		}
		// Consume authorizeID (single-use)
		s.redis.Del(ctx, authKey)
	}

	if idempotencyKey != "" {
		existingKey := fmt.Sprintf("idempotency:%s", idempotencyKey)
		existing, err := s.redis.Get(ctx, existingKey).Result()
		if err == nil && existing != "" {
			tx, _ := s.transactionRepo.FindByTransactionID(existing)
			if tx != nil {
				return s.toResponse(tx), ErrIdempotencyKeyUsed
			}
		}
	}

	var platformPrice float64
	var productID uint

	if productResp, err := s.productClient.GetProductByCode(ctx, req.ProductCode); err != nil || productResp == nil {
		return nil, ErrTxProductNotFound
	} else if !productResp.IsActive {
		return nil, ErrProductInactive
	} else {
		platformPrice = productResp.Price
		productID = uint(productResp.Id)
	}

	sellingPrice := req.SellingPrice
	if sellingPrice <= 0 {
		sellingPrice = platformPrice
	}
	margin := sellingPrice - platformPrice

	transactionID := uuid.New().String()

	// 1. Place Hold on Wallet
	if s.walletClient != nil {
		if err := s.walletClient.PlaceHoldForTransaction(ctx, userID, sellingPrice, transactionID); err != nil {
			log.Printf("Failed to place hold for user %d: %v", userID, err)
			return nil, ErrHoldFailed
		}
	}

	tx := &models.Transaction{
		TransactionID:  transactionID,
		UserID:         userID,
		ProductID:      productID,
		ProductCode:    req.ProductCode,
		CustomerNumber: req.CustomerNumber,
		Amount:         platformPrice,
		Price:          platformPrice,
		SellingPrice:   sellingPrice,
		Margin:         margin,
		Status:         string(StateInitiated),
	}

	if err := s.transactionRepo.Create(tx); err != nil {
		// Compensate: release hold
		if s.walletClient != nil {
			_ = s.walletClient.ReleaseHoldForTransaction(ctx, userID, transactionID)
		}
		return nil, err
	}

	if idempotencyKey != "" {
		key := fmt.Sprintf("idempotency:%s", idempotencyKey)
		s.redis.Set(ctx, key, transactionID, 24*time.Hour)
	}

	// 2. Call Integration Service (Async)
	go s.processExternalTransaction(context.Background(), tx)

	return s.toResponse(tx), nil
}

func (s *TransactionService) processExternalTransaction(ctx context.Context, tx *models.Transaction) {
	if s.integrationClient == nil {
		return
	}

	// For prepaid, we use TopUp
	resp, err := s.integrationClient.TopUp(ctx, &clients.TopUpRequest{
		Code:  tx.ProductCode,
		Phone: tx.CustomerNumber,
		RefID: tx.TransactionID,
	})

	if err != nil {
		log.Printf("Integration request failed for %s: %v", tx.TransactionID, err)
		s.UpdateTransactionStatus(ctx, tx.ID, &dto.UpdateStatusRequest{
			Status:  string(StateFailed),
			Message: "Provider communication failed",
		})
		return
	}

	if resp.Success {
		s.UpdateTransactionStatus(ctx, tx.ID, &dto.UpdateStatusRequest{
			Status:      string(StatePending),
			ProviderRef: resp.Data.TrxID,
			Message:     resp.Message,
		})
	} else {
		s.UpdateTransactionStatus(ctx, tx.ID, &dto.UpdateStatusRequest{
			Status:  string(StateFailed),
			Message: resp.Message,
		})
	}
}

func (s *TransactionService) GetTransaction(ctx context.Context, id uint) (*dto.TransactionResponse, error) {
	tx, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, ErrTransactionNotFound
	}

	return s.toResponse(tx), nil
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, transactionID string) (*dto.TransactionResponse, error) {
	tx, err := s.transactionRepo.FindByTransactionID(transactionID)
	if err != nil {
		return nil, ErrTransactionNotFound
	}

	return s.toResponse(tx), nil
}

func (s *TransactionService) UpdateTransactionStatus(ctx context.Context, id uint, req *dto.UpdateStatusRequest) (*dto.TransactionResponse, error) {
	tx, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, ErrTransactionNotFound
	}

	oldStatus := tx.Status
	newState := TransactionState(req.Status)
	if !s.isValidTransition(TransactionState(tx.Status), newState) {
		return nil, ErrInvalidState
	}

	tx.Status = string(newState)
	if req.ProviderRef != "" {
		tx.ProviderRef = req.ProviderRef
	}
	if req.ProviderStatus != "" {
		tx.ProviderStatus = req.ProviderStatus
	}
	if req.Message != "" {
		tx.Message = req.Message
	}

	if newState == StateSuccess || newState == StateFailed || newState == StateCancelled || newState == StateExpired {
		now := time.Now()
		tx.CompletedAt = &now
	}

	if err := s.transactionRepo.Update(tx); err != nil {
		return nil, err
	}

	// 3. Wallet Operations based on final state
	if newState == StateSuccess {
		if s.walletClient != nil {
			if err := s.walletClient.DebitForTransaction(ctx, tx.UserID, tx.SellingPrice, tx.TransactionID); err != nil {
				log.Printf("CRITICAL: Failed to debit wallet for successful transaction %s: %v", tx.TransactionID, err)
			}
		}

		// Calculate commissions for the event
		var staffCommission float64
		if s.marginService != nil {
			_, staffCommission, _ = s.marginService.CalculateTransactionMargin(tx.UserID, tx.ProductCode, tx.SellingPrice)
		}

		// Publish transaction.success event with commission info
		s.eventPublisher.Publish(ctx, "transaction_stream", "transaction.success", map[string]interface{}{
			"transaction_id":   tx.TransactionID,
			"user_id":          tx.UserID,
			"amount":           tx.SellingPrice,
			"product_code":     tx.ProductCode,
			"staff_commission": staffCommission,
		})

		// Legacy task - we'll eventually remove this
		s.publishCommissionTask(ctx, tx.TransactionID)
	} else if newState == StateFailed || newState == StateCancelled || newState == StateExpired {
		if s.walletClient != nil {
			if err := s.walletClient.ReleaseHoldForTransaction(ctx, tx.UserID, tx.TransactionID); err != nil {
				log.Printf("Failed to release hold for transaction %s: %v", tx.TransactionID, err)
			}
		}
	}

	s.logStateTransition(ctx, tx, oldStatus, string(newState), "api_update")

	return s.toResponse(tx), nil
}


func (s *TransactionService) UpdateTransactionStatusByProviderRef(ctx context.Context, providerRef string, status string, providerStatus string, message string) (*dto.TransactionResponse, error) {
	tx, err := s.transactionRepo.FindByProviderRef(providerRef)
	if err != nil {
		return nil, ErrTransactionNotFound
	}

	oldStatus := tx.Status
	newState := TransactionState(status)
	if !s.isValidTransition(TransactionState(tx.Status), newState) {
		return nil, ErrInvalidState
	}

	tx.Status = string(newState)
	tx.ProviderStatus = providerStatus
	tx.Message = message

	if newState == StateSuccess || newState == StateFailed {
		now := time.Now()
		tx.CompletedAt = &now
	}

	if err := s.transactionRepo.Update(tx); err != nil {
		return nil, err
	}

	if newState == StateSuccess {
		s.publishCommissionTask(ctx, tx.TransactionID)
	}

	s.logStateTransition(ctx, tx, oldStatus, string(newState), "webhook")

	return s.toResponse(tx), nil
}

func (s *TransactionService) publishCommissionTask(ctx context.Context, transactionID string) {
	if s.redis == nil {
		return
	}

	task := map[string]string{
		"transaction_id": transactionID,
	}
	payload, _ := json.Marshal(task)

	if err := s.redis.RPush(ctx, CommissionQueueKey, payload).Err(); err != nil {
		log.Printf("Failed to publish commission task for %s: %v", transactionID, err)
	} else {
		log.Printf("Commission task published for %s", transactionID)
	}
}

func (s *TransactionService) isValidTransition(from, to TransactionState) bool {
	validTransitions := map[TransactionState][]TransactionState{
		StateInitiated: {StatePending, StateSuccess, StateFailed},
		StatePending:   {StateSuccess, StateFailed, StateExpired, StateCancelled},
		StateSuccess:   {StateRefunded},
		StateFailed:    {},
		StateExpired:   {},
		StateCancelled: {},
		StateRefunded:  {},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, state := range allowed {
		if state == to {
			return true
		}
	}

	return false
}

func (s *TransactionService) ListTransactions(ctx context.Context, userID uint, status string, startDate, endDate string, limit, offset int) (*dto.ListTransactionsResponse, error) {
	transactions, total, err := s.transactionRepo.List(userID, status, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = *s.toResponse(&tx)
	}

	return &dto.ListTransactionsResponse{
		Transactions: responses,
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	}, nil
}

func (s *TransactionService) toResponse(tx *models.Transaction) *dto.TransactionResponse {
	resp := &dto.TransactionResponse{
		ID:             tx.ID,
		TransactionID:  tx.TransactionID,
		UserID:         tx.UserID,
		ProductCode:    tx.ProductCode,
		CustomerNumber: tx.CustomerNumber,
		Amount:         tx.Amount,
		Price:          tx.Price,
		Margin:         tx.Margin,
		Status:         tx.Status,
		ProviderRef:    tx.ProviderRef,
		Message:        tx.Message,
		CompletedAt:    tx.CompletedAt,
		CreatedAt:      tx.CreatedAt,
	}

	if tx.UpdatedAt.After(tx.CreatedAt) {
		resp.UpdatedAt = &tx.UpdatedAt
	}

	return resp
}

func (s *TransactionService) GetTransactionHistory(ctx context.Context, userID uint, cursor string, limit int) (*dto.TransactionHistoryResponse, error) {
	if cursor == "" {
		cursor = "0"
	}

	offset, _ := strconv.Atoi(cursor)

	transactions, total, err := s.transactionRepo.List(userID, "", "", "", limit, offset)
	if err != nil {
		return nil, err
	}

	var nextCursor string
	if len(transactions) == limit {
		nextCursor = fmt.Sprintf("%d", offset+limit)
	}

	responses := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = *s.toResponse(&tx)
	}

	return &dto.TransactionHistoryResponse{
		Transactions: responses,
		NextCursor:   nextCursor,
		Total:        total,
		HasMore:      nextCursor != "",
	}, nil
}

func (s *TransactionService) ProcessWebhook(ctx context.Context, payload map[string]interface{}) (*dto.TransactionResponse, error) {
	refID, ok := payload["ref_id"].(string)
	if !ok {
		return nil, errors.New("invalid webhook payload: missing ref_id")
	}

	trxID, ok := payload["trx_id"].(string)
	if !ok {
		trxID = refID
	}

	webhookKey := fmt.Sprintf("webhook:processed:%s:%s", trxID, refID)
	exists, err := s.redis.Exists(ctx, webhookKey).Result()
	if err == nil && exists == 1 {
		tx, _ := s.transactionRepo.FindByProviderRef(trxID)
		if tx != nil {
			return s.toResponse(tx), nil
		}
	}

	status, _ := payload["status"].(string)
	message, _ := payload["message"].(string)
	scCode, _ := payload["sc_code"].(string)

	if scCode != "" {
		webhookStatus, retryable, retryAfterSecs := MapRCToStatus(scCode)
		_ = retryable
		_ = retryAfterSecs

		if scCode == "74" {
			webhookStatus = string(StateRefunded)
		}

		resp, err := s.UpdateTransactionStatusByProviderRef(ctx, trxID, webhookStatus, status, message)
		if err == nil {
			s.redis.Set(ctx, webhookKey, "1", 72*time.Hour)
		} else if err != nil && !s.isRetryableError(err) {
			s.pushToDeadLetterQueue(ctx, refID, trxID, payload, err.Error())
		}
		return resp, err
	}

	webhookStatus := mapDigiflazzStatus(status, scCode)

	resp, err := s.UpdateTransactionStatusByProviderRef(ctx, trxID, webhookStatus, status, message)
	if err == nil {
		s.redis.Set(ctx, webhookKey, "1", 72*time.Hour)
	} else if err != nil && !s.isRetryableError(err) {
		s.pushToDeadLetterQueue(ctx, refID, trxID, payload, err.Error())
	}

	return resp, err
}

func (s *TransactionService) VerifyWebhookSignature(payload []byte, timestamp, signature string) bool {
	if s.webhookSecret == "" {
		return true
	}
	message := string(payload) + timestamp + s.webhookSecret
	hash := md5.Sum([]byte(message))
	expected := hex.EncodeToString(hash[:])
	return expected == signature
}

func (s *TransactionService) pushToDeadLetterQueue(ctx context.Context, refID, trxID string, payload map[string]interface{}, errorMsg string) {
	if s.redis == nil {
		return
	}
	entry := map[string]interface{}{
		"ref_id":    refID,
		"trx_id":    trxID,
		"payload":   payload,
		"error":     errorMsg,
		"failed_at": time.Now().Unix(),
	}
	entryJSON, _ := json.Marshal(entry)
	s.redis.LPush(ctx, "digiflazz_webhook_dlq", entryJSON)
}

func (s *TransactionService) isRetryableError(err error) bool {
	return err == ErrInvalidState
}

func mapDigiflazzStatus(status, scCode string) string {
	switch status {
	case "Sukses", "Success":
		return string(StateSuccess)
	case "Gagal", "Failed":
		return string(StateFailed)
	case "Pending", "Processing":
		return string(StatePending)
	case "Expired", "expired":
		return string(StateExpired)
	case "Cancelled", "cancelled":
		return string(StateCancelled)
	case "Refund", "refund":
		return string(StateRefunded)
	default:
		if scCode == "00" {
			return string(StateSuccess)
		}
		if scCode == "03" {
			return string(StatePending)
		}
		return string(StateFailed)
	}
}

type DigiflazzRCStatus struct {
	InternalStatus string
	Retryable      bool
	RetryAfterSecs int
	UserMessageID  string
}

var digiflazzRCStatusMap = map[string]DigiflazzRCStatus{
	"00": {string(StateSuccess), false, 0, "DIGIFLAZZ_SUCCESS"},
	"01": {string(StateFailed), true, 2, "DIGIFLAZZ_TIMEOUT"},
	"02": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"03": {string(StatePending), false, 0, "DIGIFLAZZ_PENDING"},
	"40": {string(StateFailed), false, 0, "DIGIFLAZZ_PAYLOAD_ERROR"},
	"41": {string(StateFailed), false, 0, "DIGIFLAZZ_INVALID_SIGNATURE"},
	"42": {string(StateFailed), false, 0, "DIGIFLAZZ_SELLER_NOT_FOUND"},
	"43": {string(StateFailed), false, 0, "DIGIFLAZZ_SKU_NOT_FOUND"},
	"44": {string(StateFailed), false, 0, "DIGIFLAZZ_INSUFFICIENT_DEPOSIT"},
	"45": {string(StateFailed), false, 0, "DIGIFLAZZ_IP_BLOCKED"},
	"47": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"49": {string(StateFailed), false, 0, "DIGIFLAZZ_REF_ID_DUPLICATE"},
	"50": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"51": {string(StateFailed), false, 0, "DIGIFLAZZ_NUMBER_BLOCKED"},
	"52": {string(StateFailed), false, 0, "DIGIFLAZZ_PREFIX_MISMATCH"},
	"53": {string(StateFailed), false, 0, "DIGIFLAZZ_PRODUCT_UNAVAILABLE"},
	"54": {string(StateFailed), false, 0, "DIGIFLAZZ_INVALID_NUMBER"},
	"55": {string(StateFailed), false, 0, "DIGIFLAZZ_PRODUCT_DISRUPTION"},
	"57": {string(StateFailed), false, 0, "DIGIFLAZZ_INVALID_NUMBER"},
	"58": {string(StateFailed), true, 900, "DIGIFLAZZ_CUT_OFF"},
	"59": {string(StateFailed), false, 0, "DIGIFLAZZ_AREA_UNSUPPORTED"},
	"60": {string(StateFailed), false, 0, "DIGIFLAZZ_BILL_NOT_AVAILABLE"},
	"61": {string(StateFailed), false, 0, "DIGIFLAZZ_INSUFFICIENT_DEPOSIT"},
	"62": {string(StateFailed), true, 30, "DIGIFLAZZ_SELLER_DISRUPTION"},
	"63": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"64": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"66": {string(StateFailed), false, 0, "DIGIFLAZZ_MAINTENANCE"},
	"67": {string(StateFailed), false, 0, "DIGIFLAZZ_SELLER_NOT_VERIFIED"},
	"68": {string(StateFailed), false, 0, "DIGIFLAZZ_OUT_OF_STOCK"},
	"69": {string(StateFailed), true, 10, "DIGIFLAZZ_PRICE_MISMATCH"},
	"70": {string(StateFailed), true, 2, "DIGIFLAZZ_BILLER_TIMEOUT"},
	"71": {string(StateFailed), true, 5, "DIGIFLAZZ_PRODUCT_UNSTABLE"},
	"72": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"73": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"74": {string(StateRefunded), false, 0, "DIGIFLAZZ_REFUND"},
	"80": {string(StateFailed), false, 0, "DIGIFLAZZ_ACCOUNT_BLOCKED"},
	"81": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"82": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"83": {string(StateFailed), true, 240, "DIGIFLAZZ_PRICELIST_LIMIT"},
	"84": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"85": {string(StateFailed), true, 60, "DIGIFLAZZ_RATE_LIMIT"},
	"86": {string(StateFailed), true, 60, "DIGIFLAZZ_PLN_INQUIRY_LIMIT"},
	"87": {string(StateFailed), false, 0, "DIGIFLAZZ_EMONEY_MULTIPLE"},
	"88": {string(StateFailed), false, 0, "DIGIFLAZZ_GENERAL_FAILURE"},
	"99": {string(StatePending), true, 5, "DIGIFLAZZ_ROUTER_ERROR"},
}

func MapRCToStatus(rc string) (status string, retryable bool, retryAfterSecs int) {
	if rcInfo, ok := digiflazzRCStatusMap[rc]; ok {
		return rcInfo.InternalStatus, rcInfo.Retryable, rcInfo.RetryAfterSecs
	}
	return string(StateFailed), false, 0
}

func GetUserMessage(rc string) string {
	if rcInfo, ok := digiflazzRCStatusMap[rc]; ok {
		return rcInfo.UserMessageID
	}
	return "DIGIFLAZZ_UNKNOWN_ERROR"
}

type TransactionStateMachine struct {
	currentState TransactionState
}

func NewTransactionStateMachine(currentState string) *TransactionStateMachine {
	return &TransactionStateMachine{currentState: TransactionState(currentState)}
}

func (sm *TransactionStateMachine) CanTransitionTo(newState TransactionState) bool {
	ts := &TransactionService{}
	return ts.isValidTransition(sm.currentState, newState)
}

func (sm *TransactionStateMachine) TransitionTo(newState TransactionState) error {
	ts := &TransactionService{}
	if !ts.isValidTransition(sm.currentState, newState) {
		return ErrInvalidState
	}
	sm.currentState = newState
	return nil
}

func (sm *TransactionStateMachine) GetCurrentState() TransactionState {
	return sm.currentState
}

func (sm *TransactionStateMachine) IsTerminal() bool {
	return sm.currentState == StateSuccess ||
		sm.currentState == StateFailed ||
		sm.currentState == StateExpired ||
		sm.currentState == StateCancelled ||
		sm.currentState == StateRefunded
}

func (s *TransactionService) CancelTransaction(ctx context.Context, id uint, userID uint, reason string) (*dto.TransactionResponse, error) {
	tx, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, ErrTransactionNotFound
	}

	if tx.UserID != userID {
		return nil, errors.New("unauthorized: cannot cancel another user's transaction")
	}

	oldStatus := tx.Status
	newState := StateCancelled

	if !s.isValidTransition(TransactionState(tx.Status), newState) {
		return nil, ErrInvalidState
	}

	tx.Status = string(newState)
	tx.PreviousStatus = oldStatus
	tx.StatusChangeReason = reason

	if err := s.transactionRepo.Update(tx); err != nil {
		return nil, err
	}

	s.logStateTransition(ctx, tx, oldStatus, string(newState), "user_cancel")

	return s.toResponse(tx), nil
}

type TransitionValidator struct{}

func (v *TransitionValidator) CanTransition(from, to TransactionState, tx *Transaction) error {
	switch from {
	case StateInitiated:
		if to == StatePending && (tx.ProviderRef == "" || tx.ProviderRef == "-") {
			return ErrMissingRefID
		}
		if to == StateSuccess && tx.Amount <= 0 {
			return ErrAmountMustBePositive
		}
	case StatePending:
		if to == StateSuccess && time.Since(tx.CreatedAt) > 15*time.Minute {
			return ErrCannotSucceedAfterExpiry
		}
	case StateSuccess:
		if to == StateRefunded && time.Since(tx.UpdatedAt) > 24*time.Hour {
			return ErrRefundWindowClosed
		}
	}
	return nil
}

// GetReports returns aggregated KPIs, sales trend, and staff performance for the given date range
func (s *TransactionService) GetReports(ctx context.Context, startDate, endDate string, userID uint) (*dto.ReportsResponse, error) {
	// Parse dates to ensure valid format
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Get KPIs
	kpiResult, err := s.transactionRepo.GetKPIs(startDate, endDate, userID)
	if err != nil {
		return nil, err
	}

	kpi := dto.ReportKPIResponse{
		TotalSales:       kpiResult.TotalSales,
		PlatformProfit:   kpiResult.PlatformProfit,
		TransactionCount: kpiResult.TotalCount,
		SuccessRate:      s.calculateSuccessRate(kpiResult.SuccessCount, kpiResult.TotalCount),
		PeriodStart:      startDate,
		PeriodEnd:        endDate,
	}
	// Staff count from result
	kpi.StaffCount = int(kpiResult.StaffCount)

	// Sales Trend
	salesTrend, err := s.transactionRepo.GetSalesTrend(startDate, endDate)
	if err != nil {
		return nil, err
	}
	salesTrendItems := make([]dto.ReportSalesTrendItem, len(salesTrend))
	for i, item := range salesTrend {
		salesTrendItems[i] = dto.ReportSalesTrendItem{
			Date:  item.Date,
			Sales: item.Sales,
			Count: item.Count,
		}
	}

	// Staff Performance
	staffPerf, err := s.transactionRepo.GetStaffPerformance(startDate, endDate)
	if err != nil {
		return nil, err
	}
	staffPerfItems := make([]dto.ReportStaffPerformanceItem, len(staffPerf))
	for i, item := range staffPerf {
		staffPerfItems[i] = dto.ReportStaffPerformanceItem{
			StaffID:          item.StaffID,
			StaffName:        item.StaffName,
			TransactionCount: item.TransactionCount,
			TotalSales:       0, // can compute from transactions if needed
			TotalCommission:  item.TotalCommission,
			SuccessRate:      0, // could compute from commissions
		}
	}

	return &dto.ReportsResponse{
		KPIs:             []dto.ReportKPIResponse{kpi},
		SalesTrend:       salesTrendItems,
		StaffPerformance: staffPerfItems,
	}, nil
}

func (s *TransactionService) calculateSuccessRate(successCount, totalCount int) float64 {
	if totalCount == 0 {
		return 0
	}
	return float64(successCount) / float64(totalCount) * 100
}
