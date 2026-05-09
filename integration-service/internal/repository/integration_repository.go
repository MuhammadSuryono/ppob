package repository

import (
	"github.com/yontech/ppob/integration-service/internal/models"
	"gorm.io/gorm"
)

type IntegrationLogRepository struct {
	db *gorm.DB
}

func NewIntegrationLogRepository(db *gorm.DB) *IntegrationLogRepository {
	return &IntegrationLogRepository{db: db}
}

func (r *IntegrationLogRepository) Create(log *models.IntegrationLog) error {
	return r.db.Create(log).Error
}

func (r *IntegrationLogRepository) Update(log *models.IntegrationLog) error {
	return r.db.Save(log).Error
}

func (r *IntegrationLogRepository) FindByRequestID(requestID string) (*models.IntegrationLog, error) {
	var log models.IntegrationLog
	err := r.db.Where("request_id = ?", requestID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *IntegrationLogRepository) FindByTransactionID(transactionID string) (*models.IntegrationLog, error) {
	var log models.IntegrationLog
	err := r.db.Where("transaction_id = ?", transactionID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

type ProviderConfigRepository struct {
	db *gorm.DB
}

func NewProviderConfigRepository(db *gorm.DB) *ProviderConfigRepository {
	return &ProviderConfigRepository{db: db}
}

func (r *ProviderConfigRepository) FindByProvider(provider string) (*models.ProviderConfig, error) {
	var config models.ProviderConfig
	err := r.db.Where("provider = ? AND is_active = ?", provider, true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ProviderConfigRepository) List() ([]models.ProviderConfig, error) {
	var configs []models.ProviderConfig
	err := r.db.Where("is_active = ?", true).Order("priority ASC").Find(&configs).Error
	return configs, err
}

func (r *ProviderConfigRepository) Create(config *models.ProviderConfig) error {
	return r.db.Create(config).Error
}

func (r *ProviderConfigRepository) Update(config *models.ProviderConfig) error {
	return r.db.Save(config).Error
}