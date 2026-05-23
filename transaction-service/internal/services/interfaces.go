package services

import (
	"context"
	"github.com/yontech/ppob/shared/proto/product"
	"github.com/yontech/ppob/transaction-service/internal/clients"
)

type ProductClient interface {
	GetProductByCode(ctx context.Context, skuCode string) (*product.GetProductResponse, error)
	ValidateProduct(ctx context.Context, productID uint, expectedPrice float64) (*product.ValidateProductResponse, error)
}

type WalletClient interface {
	PlaceHoldForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error
	ReleaseHoldForTransaction(ctx context.Context, userID uint, transactionID string) error
	DebitForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error
	CreditWallet(ctx context.Context, userID uint, amount float64, referenceID, referenceType string) error
	DebitWallet(ctx context.Context, userID uint, amount float64, referenceID, referenceType string) error
}

type IntegrationClient interface {
	TopUp(ctx context.Context, req *clients.TopUpRequest) (*clients.IntegrationResponse, error)
}
