package repository

import (
	"strings"

	"github.com/yontech/ppob/product-service/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByCode(code string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("code = ?", code).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) List(categoryID uint, brand string, status string, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if brand != "" {
		query = query.Where("lower(brand) = ?", strings.ToLower(brand))
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&products).Error

	return products, total, err
}

type ProductWithPrice struct {
	models.Product
	MitraSellingPrice *float64 `gorm:"column:mitra_selling_price"`
}

func (r *ProductRepository) ListWithMitraPrice(mitraID uint, categoryID uint, brand string, status string, limit, offset int) ([]ProductWithPrice, int64, error) {
	var results []ProductWithPrice
	var total int64

	query := r.db.Table("products").
		Select("products.*, mpp.selling_price as mitra_selling_price").
		Joins("LEFT JOIN mitra_product_prices mpp ON mpp.product_id = products.id AND mpp.mitra_id = ?", mitraID).
		Where("products.deleted_at IS NULL")

	if categoryID > 0 {
		query = query.Where("products.category_id = ?", categoryID)
	}
	if brand != "" {
		query = query.Where("lower(products.brand) = ?", strings.ToLower(brand))
	}
	if status != "" {
		query = query.Where("products.status = ?", status)
	}

	query.Count(&total)
	err := query.Order("products.created_at DESC").Limit(limit).Offset(offset).Scan(&results).Error

	return results, total, err
}

func (r *ProductRepository) Search(keyword string, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("name ILIKE ? OR code ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	query.Count(&total)
	err := query.Limit(limit).Offset(offset).Find(&products).Error

	return products, total, err
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) FindByCode(code string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("code = ?", code).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("status = ?", "active").Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
