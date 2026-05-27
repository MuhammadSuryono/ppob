package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	logRepo         *repository.IntegrationLogRepository
	providerRepo    *repository.ProviderConfigRepository
	redis           *redis.Client
	cfg             *config.Config
	digiflazzClient *DigiflazzClient
}

func NewIntegrationService(
	logRepo *repository.IntegrationLogRepository,
	providerRepo *repository.ProviderConfigRepository,
	redis *redis.Client,
	cfg *config.Config,
	digiflazzClient *DigiflazzClient,
) *IntegrationService {
	return &IntegrationService{
		logRepo:         logRepo,
		providerRepo:    providerRepo,
		redis:           redis,
		cfg:             cfg,
		digiflazzClient: digiflazzClient,
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

	resp, err := s.digiflazzClient.TopUp(ctx, &TopUpRequest{
		Code:  req.ProductCode,
		Phone: req.CustomerNumber,
		RefID: req.RefID,
	})

	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = err.Error()
		logEntry.DurationMs = int(time.Since(startTime).Milliseconds())
		s.logRepo.Update(logEntry)
		return nil, err
	}

	response := &dto.DigiflazzTransactionResponse{
		Success:      resp.Success,
		Message:      resp.Message,
		RefID:        resp.Data.RefID,
		TrxID:        resp.Data.TrxID,
		Status:       resp.Data.Status,
		Price:        0, // Should be parsed if needed
		ScCode:       resp.Data.ScCode,
		ScMessage:    resp.Data.ScMessage,
		CustomerName: resp.Data.Message,
	}

	logEntry.ResponseData = s.mustMarshal(resp)
	if response.Success {
		logEntry.Status = "success"
	} else {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = response.Message
	}
	logEntry.DurationMs = int(time.Since(startTime).Milliseconds())
	s.logRepo.Update(logEntry)

	return response, nil
}

func (s *IntegrationService) InquiryDigiflazz(ctx context.Context, req *dto.DigiflazzTransactionRequest) (*dto.DigiflazzTransactionResponse, error) {
	startTime := time.Now()

	logEntry := &models.IntegrationLog{
		Provider:      "digiflazz",
		Action:        "inquiry",
		RequestID:     req.RefID,
		TransactionID: req.RefID,
		Status:        "pending",
		RequestData:   s.mustMarshal(req),
	}
	s.logRepo.Create(logEntry)

	resp, err := s.digiflazzClient.PostpaidInquiry(ctx, &TransactionRequest{
		Code:  req.ProductCode,
		Phone: req.CustomerNumber,
		RefID: req.RefID,
	})

	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = err.Error()
		logEntry.DurationMs = int(time.Since(startTime).Milliseconds())
		s.logRepo.Update(logEntry)
		return nil, err
	}

	response := &dto.DigiflazzTransactionResponse{
		Success:      resp.Success,
		Message:      resp.Message,
		RefID:        resp.Data.RefID,
		TrxID:        resp.Data.TrxID,
		Status:       resp.Data.Status,
		Price:        0, // Should be bill amount
		ScCode:       resp.Data.ScCode,
		ScMessage:    resp.Data.ScMessage,
		CustomerName: resp.Data.Message,
	}

	logEntry.ResponseData = s.mustMarshal(resp)
	if response.Success {
		logEntry.Status = "success"
	} else {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = response.Message
	}
	logEntry.DurationMs = int(time.Since(startTime).Milliseconds())
	s.logRepo.Update(logEntry)

	return response, nil
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

func (s *IntegrationService) mustMarshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}