package services

import (
	"errors"

	"github.com/yontech/ppob/user-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrNoMarginSetting    = errors.New("no margin setting found for staff")
	ErrMarginSettingExists = errors.New("margin setting already exists")
)

type MarginService struct {
	db *gorm.DB
}

func NewMarginService(db *gorm.DB) *MarginService {
	return &MarginService{db: db}
}

type CommissionResult struct {
	SchemeType  string  `json:"scheme_type"`
	Value       float64 `json:"value"`
	Commission  float64 `json:"commission"`
	Margin      float64 `json:"margin"`
}

func (s *MarginService) GetStaffCommission(staffID, productID uint, productCode string, sellingPrice, platformPrice float64) (*CommissionResult, error) {
	margin := sellingPrice - platformPrice
	if margin < 0 {
		margin = 0
	}

	override := &models.StaffProductMarginOverride{}
	err := s.db.Where("staff_id = ? AND product_id = ? AND is_active = ?", staffID, productID, true).First(override).Error

	schemeType := "MarginShare"
	value := float64(60)

	if err == nil {
		if override.FixedMargin > 0 {
			schemeType = "FixedAllowance"
			value = override.FixedMargin
		} else if override.MarginPercent > 0 {
			schemeType = "MarginShare"
			value = override.MarginPercent
		}
	} else {
		global := &models.StaffGlobalMarginSetting{}
		err = s.db.Where("staff_id = ? AND is_active = ?", staffID, true).First(global).Error
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrNoMarginSetting
			}
			return nil, err
		}

		schemeType = global.SchemeType
		if schemeType == "MarginShare" {
			value = global.GlobalMarginPercent
		} else {
			value = global.FixedAllowance
		}
	}

	var commission float64
	if schemeType == "MarginShare" {
		commission = margin * (value / 100)
	} else {
		commission = value
		if commission > margin {
			commission = margin
		}
	}

	return &CommissionResult{
		SchemeType: schemeType,
		Value:      value,
		Commission: commission,
		Margin:     margin,
	}, nil
}

func (s *MarginService) SetStaffMargin(mitraID, staffID uint, schemeType string, marginValue float64) (*models.StaffGlobalMarginSetting, error) {
	var existing models.StaffGlobalMarginSetting
	err := s.db.Where("mitra_id = ? AND staff_id = ?", mitraID, staffID).First(&existing).Error

	if err == nil {
		existing.SchemeType = schemeType
		if schemeType == "MarginShare" {
			existing.GlobalMarginPercent = marginValue
			existing.FixedAllowance = 0
		} else {
			existing.FixedAllowance = marginValue
			existing.GlobalMarginPercent = 0
		}
		existing.IsActive = true
		err = s.db.Save(&existing).Error
		return &existing, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		setting := &models.StaffGlobalMarginSetting{
			MitraID:            mitraID,
			StaffID:            staffID,
			SchemeType:         schemeType,
			GlobalMarginPercent: 0,
			FixedAllowance:     0,
			IsActive:           true,
		}

		if schemeType == "MarginShare" {
			setting.GlobalMarginPercent = marginValue
		} else {
			setting.FixedAllowance = marginValue
		}

		err = s.db.Create(setting).Error
		return setting, err
	}

	return nil, err
}

func (s *MarginService) SetProductMarginOverride(mitraID, staffID, productID uint, productCode string, marginValue float64, schemeType string) (*models.StaffProductMarginOverride, error) {
	var existing models.StaffProductMarginOverride
	err := s.db.Where("mitra_id = ? AND staff_id = ? AND product_id = ?", mitraID, staffID, productID).First(&existing).Error

	if err == nil {
		existing.MarginPercent = marginValue
		existing.ProductCode = productCode
		existing.IsActive = true
		err = s.db.Save(&existing).Error
		return &existing, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		override := &models.StaffProductMarginOverride{
			MitraID:       mitraID,
			StaffID:       staffID,
			ProductID:     productID,
			ProductCode:   productCode,
			MarginPercent: marginValue,
			IsActive:      true,
		}

		if schemeType == "FixedAllowance" {
			override.FixedMargin = marginValue
		}

		err = s.db.Create(override).Error
		return override, err
	}

	return nil, err
}

func (s *MarginService) GetStaffMarginSettings(staffID uint) (*models.StaffGlobalMarginSetting, error) {
	var setting models.StaffGlobalMarginSetting
	err := s.db.Where("staff_id = ? AND is_active = ?", staffID, true).First(&setting).Error
	return &setting, err
}

func (s *MarginService) GetStaffProductOverrides(staffID uint) ([]models.StaffProductMarginOverride, error) {
	var overrides []models.StaffProductMarginOverride
	err := s.db.Where("staff_id = ? AND is_active = ?", staffID, true).Find(&overrides).Error
	return overrides, err
}

func (s *MarginService) DeactivateStaffMargin(mitraID, staffID uint) error {
	return s.db.Model(&models.StaffGlobalMarginSetting{}).
		Where("mitra_id = ? AND staff_id = ?", mitraID, staffID).
		Update("is_active", false).Error
}