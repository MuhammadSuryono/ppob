package clients

import (
	"context"
	"fmt"

	"github.com/yontech/ppob/shared/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	client product.ProductServiceClient
	conn   *grpc.ClientConn
}

func NewProductClient(address string) (*ProductClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	return &ProductClient{
		client: product.NewProductServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *ProductClient) Close() error {
	return c.conn.Close()
}

func (c *ProductClient) GetProductByCode(ctx context.Context, skuCode string) (*product.GetProductResponse, error) {
	return c.client.GetProduct(ctx, &product.GetProductRequest{
		SkuCode: skuCode,
	})
}

func (c *ProductClient) GetInquiryProduct(ctx context.Context, categoryID uint, brand string) (*product.GetProductResponse, error) {
	return c.client.GetInquiryProduct(ctx, &product.GetInquiryProductRequest{
		CategoryId: uint32(categoryID),
		Brand:      brand,
	})
}

func (c *ProductClient) ValidateProduct(ctx context.Context, productID uint, expectedPrice float64) (*product.ValidateProductResponse, error) {
	return c.client.ValidateProduct(ctx, &product.ValidateProductRequest{
		ProductId:     uint32(productID),
		ExpectedPrice: expectedPrice,
	})
}
