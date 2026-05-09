package services

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/yontech/ppob/product-service/config"
	"github.com/yontech/ppob/product-service/internal/dto"
	"github.com/yontech/ppob/product-service/internal/models"
	"github.com/yontech/ppob/product-service/internal/repository"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	redis        *redis.Client
	cfg          *config.Config
}

func NewProductService(
	productRepo *repository.ProductRepository,
	categoryRepo *repository.CategoryRepository,
	redis *redis.Client,
	cfg *config.Config,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		redis:        redis,
		cfg:          cfg,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		CategoryID:  product.CategoryID,
		Provider:    product.Provider,
		Price:       product.Price,
		PriceAPI:    product.PriceAPI,
		Stock:       product.Stock,
		Status:      product.Status,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
	}, nil
}

func (s *ProductService) GetProductByCode(ctx context.Context, code string) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByCode(code)
	if err != nil {
		return nil, ErrProductNotFound
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		CategoryID:  product.CategoryID,
		Provider:    product.Provider,
		Price:       product.Price,
		PriceAPI:    product.PriceAPI,
		Stock:       product.Stock,
		Status:      product.Status,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
	}, nil
}

func (s *ProductService) ListProducts(ctx context.Context, categoryID uint, status string, limit, offset int) (*dto.ListProductsResponse, error) {
	products, total, err := s.productRepo.List(categoryID, status, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			ProductCode: p.ProductCode,
			ProductName: p.ProductName,
			CategoryID:  p.CategoryID,
			Provider:    p.Provider,
			Price:       p.Price,
			PriceAPI:    p.PriceAPI,
			Stock:       p.Stock,
			Status:      p.Status,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
		}
	}

	return &dto.ListProductsResponse{
		Products: responses,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (s *ProductService) SearchProducts(ctx context.Context, keyword string, limit, offset int) (*dto.ListProductsResponse, error) {
	products, total, err := s.productRepo.Search(keyword, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			ProductCode: p.ProductCode,
			ProductName: p.ProductName,
			CategoryID:  p.CategoryID,
			Provider:    p.Provider,
			Price:       p.Price,
			PriceAPI:    p.PriceAPI,
			Stock:       p.Stock,
			Status:      p.Status,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
		}
	}

	return &dto.ListProductsResponse{
		Products: responses,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &models.Product{
		ProductCode: req.ProductCode,
		ProductName: req.ProductName,
		CategoryID:  req.CategoryID,
		Provider:    req.Provider,
		Price:       req.Price,
		Stock:       req.Stock,
		Description: req.Description,
		Status:      "active",
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		CategoryID:  product.CategoryID,
		Provider:    product.Provider,
		Price:       product.Price,
		Stock:       product.Stock,
		Status:      product.Status,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if req.ProductName != "" {
		product.ProductName = req.ProductName
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if req.Provider != "" {
		product.Provider = req.Provider
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock != 0 {
		product.Stock = req.Stock
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	if req.Description != "" {
		product.Description = req.Description
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		CategoryID:  product.CategoryID,
		Provider:    product.Provider,
		Price:       product.Price,
		Stock:       product.Stock,
		Status:      product.Status,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		return ErrProductNotFound
	}

	return s.productRepo.Delete(id)
}

func (s *ProductService) ListCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = dto.CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Code:        c.Code,
			Description: c.Description,
			Icon:        c.Icon,
			SortOrder:   c.SortOrder,
			Status:      c.Status,
		}
	}

	return responses, nil
}

func (s *ProductService) GetCategory(ctx context.Context, id uint) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
		Icon:        category.Icon,
		SortOrder:   category.SortOrder,
		Status:      category.Status,
	}, nil
}

func (s *ProductService) UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Code != "" {
		category.Code = req.Code
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.SortOrder > 0 {
		category.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		category.Status = req.Status
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
		Icon:        category.Icon,
		SortOrder:   category.SortOrder,
		Status:      category.Status,
	}, nil
}

func (s *ProductService) DeleteCategory(ctx context.Context, id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return ErrCategoryNotFound
	}

	return s.categoryRepo.Delete(id)
}

func (s *ProductService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &models.Category{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
		Icon:        category.Icon,
		SortOrder:   category.SortOrder,
		Status:      category.Status,
	}, nil
}