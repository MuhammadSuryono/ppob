package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/integration-service/config"
	"github.com/yontech/ppob/integration-service/internal/dto"
	"github.com/yontech/ppob/integration-service/internal/models"
	"github.com/yontech/ppob/integration-service/internal/repository"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrRequestFailed   = errors.New("external request failed")
)

type IntegrationService struct {
	logRepo     *repository.IntegrationLogRepository
	providerRepo *repository.ProviderConfigRepository
	redis       *redis.Client
	cfg         *config.Config
}

func NewIntegrationService(
	logRepo *repository.IntegrationLogRepository,
	providerRepo *repository.ProviderConfigRepository,
	redis *redis.Client,
	cfg *config.Config,
) *IntegrationService {
	return &IntegrationService{
		logRepo:      logRepo,
		providerRepo: providerRepo,
		redis:        redis,
		cfg:          cfg,
	}
}

func (s *IntegrationService) InitiateDigiflazzTransaction(ctx context.Context, userID uint, req *dto.DigiflazzTransactionRequest) (*dto.DigiflazzTransactionResponse, error) {
	startTime := time.Now()

	logEntry := &models.IntegrationLog{
		Provider:      "digiflazz",
		Action:        "transaction",
		RequestID:     req.RefID,
		TransactionID: req.RefID,
		Status:        "pending",
		RequestData:   s.mustMarshal(req),
	}
	s.logRepo.Create(logEntry)

	payload := map[string]interface{}{
		"username": s.cfg.DigiflazzKey,
		"code":     req.ProductCode,
		"phone":    req.CustomerNumber,
		"ref_id":   req.RefID,
		"sign":     s.generateSignature(req.RefID),
	}

	body, err := s.callDigiflazzAPI(ctx, "/transaction", payload)
	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = err.Error()
		logEntry.DurationMs = int(time.Since(startTime).Milliseconds())
		s.logRepo.Update(logEntry)
		return nil, err
	}

	var response dto.DigiflazzTransactionResponse
	json.Unmarshal(body, &response)

	logEntry.ResponseData = string(body)
	if response.Success {
		logEntry.Status = "success"
	} else {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = response.Message
	}
	logEntry.DurationMs = int(time.Since(startTime).Milliseconds())
	s.logRepo.Update(logEntry)

	return &response, nil
}

func (s *IntegrationService) HandleDigiflazzCallback(ctx context.Context, req *dto.DigiflazzCallbackRequest) error {
	logEntry, err := s.logRepo.FindByRequestID(req.RefID)
	if err != nil {
		return err
	}

	logEntry.ResponseData = s.mustMarshal(req)
	if req.Status == "Sukses" {
		logEntry.Status = "success"
	} else {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = req.Message
	}

	return s.logRepo.Update(logEntry)
}

func (s *IntegrationService) ListProviders(ctx context.Context) ([]dto.ProviderResponse, error) {
	providers, err := s.providerRepo.List()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = dto.ProviderResponse{
			Provider:   p.Provider,
			Name:       p.Name,
			APIURL:     p.APIURL,
			IsActive:   p.IsActive,
			RateLimit:  p.RateLimit,
			TimeoutSec: p.TimeoutSec,
		}
	}

	return responses, nil
}

func (s *IntegrationService) TryLockWebhook(ctx context.Context, key string) (bool, error) {
	lockKey := fmt.Sprintf("webhook_lock:%s", key)
	result, err := s.redis.SetNX(ctx, lockKey, "locked", 5*time.Minute).Result()
	return result, err
}

func (s *IntegrationService) ReleaseWebhookLock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("webhook_lock:%s", key)
	return s.redis.Del(ctx, lockKey).Err()
}

func (s *IntegrationService) UpdateTransactionStatus(ctx context.Context, refID string, req *dto.UpdateStatusRequest) error {
	updateKey := fmt.Sprintf("transaction_status:%s", refID)
	s.redis.Set(ctx, updateKey, req.Status, 24*time.Hour)
	return nil
}

func (s *IntegrationService) GetTransactionStatus(ctx context.Context, refID string) (string, error) {
	updateKey := fmt.Sprintf("transaction_status:%s", refID)
	return s.redis.Get(ctx, updateKey).Result()
}

func (s *IntegrationService) callDigiflazzAPI(ctx context.Context, endpoint string, payload map[string]interface{}) ([]byte, error) {
	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, "POST", s.cfg.DigiflazzURL+endpoint, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("digiflazz API error: %s", string(body))
	}

	return body, nil
}

func (s *IntegrationService) generateSignature(refID string) string {
	return fmt.Sprintf("%s%s%s", s.cfg.DigiflazzKey, refID, s.cfg.DigiflazzSecret)
}

func (s *IntegrationService) mustMarshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}