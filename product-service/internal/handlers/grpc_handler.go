package handlers

import (
	"context"

	"github.com/yontech/ppob/product-service/internal/dto"
	"github.com/yontech/ppob/product-service/internal/services"
	"github.com/yontech/ppob/shared/proto/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGRPCHandler struct {
	product.UnimplementedProductServiceServer
	productService         *services.ProductService
	priceValidationService *services.PriceValidationService
}

func NewProductGRPCHandler(productService *services.ProductService, priceValidationService *services.PriceValidationService) *ProductGRPCHandler {
	return &ProductGRPCHandler{
		productService:         productService,
		priceValidationService: priceValidationService,
	}
}

func (h *ProductGRPCHandler) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	var p *dto.ProductResponse
	var err error

	if req.ProductId > 0 {
		p, err = h.productService.GetProduct(ctx, uint(req.ProductId))
	} else if req.SkuCode != "" {
		p, err = h.productService.GetProductByCode(ctx, req.SkuCode)
	} else {
		return nil, status.Error(codes.InvalidArgument, "product_id or sku_code is required")
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	return &product.GetProductResponse{
		Id:       uint32(p.ID),
		Name:     p.Name,
		SkuCode:  p.Code,
		Price:    p.Price,
		IsActive: p.Status == "active",
	}, nil
}

/*
func (h *ProductGRPCHandler) GetInquiryProduct(ctx context.Context, req *product.GetInquiryProductRequest) (*product.GetProductResponse, error) {
	p, err := h.productService.GetInquiryProduct(ctx, uint(req.CategoryId), req.Brand)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inquiry product not found for brand %s: %v", req.Brand, err)
	}

	return &product.GetProductResponse{
		Id:       uint32(p.ID),
		Name:     p.Name,
		SkuCode:  p.Code,
		Price:    p.Price,
		IsActive: p.Status == "active",
	}, nil
}
*/

func (h *ProductGRPCHandler) ValidateProduct(ctx context.Context, req *product.ValidateProductRequest) (*product.ValidateProductResponse, error) {
	p, err := h.productService.GetProduct(ctx, uint(req.ProductId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	result := h.priceValidationService.ValidatePrice(p.Code, req.ExpectedPrice)

	return &product.ValidateProductResponse{
		IsValid:      result.Valid,
		Message:      result.ValidationError,
		CurrentPrice: result.PlatformPrice,
	}, nil
}
