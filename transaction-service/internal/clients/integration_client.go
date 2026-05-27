package clients

import (
	"context"
	"fmt"

	"github.com/yontech/ppob/shared/proto/integration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IntegrationClient struct {
	client integration.IntegrationServiceClient
	conn   *grpc.ClientConn
}

func NewIntegrationClient(address string) (*IntegrationClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to integration service: %w", err)
	}

	return &IntegrationClient{
		client: integration.NewIntegrationServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *IntegrationClient) Close() error {
	return c.conn.Close()
}

type TopUpRequest struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
	RefID string `json:"ref_id"`
}

type IntegrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		RefID     string `json:"ref_id"`
		TrxID     string `json:"trx_id"`
		Status    string `json:"status"`
		ScCode    string `json:"sc_code"`
		ScMessage string `json:"sc_message"`
	} `json:"data"`
}

func (c *IntegrationClient) TopUp(ctx context.Context, req *TopUpRequest) (*IntegrationResponse, error) {
	resp, err := c.client.TopUp(ctx, &integration.TopUpRequest{
		ProductCode:    req.Code,
		CustomerNumber: req.Phone,
		RefId:          req.RefID,
	})

	if err != nil {
		return nil, err
	}

	result := &IntegrationResponse{
		Success: resp.Success,
		Message: resp.Message,
	}

	if resp.Data != nil {
		result.Data.RefID = resp.Data.RefId
		result.Data.TrxID = resp.Data.TrxId
		result.Data.Status = resp.Data.Status
		result.Data.ScCode = resp.Data.ScCode
		result.Data.ScMessage = resp.Data.ScMessage
	}

	return result, nil
}

func (c *IntegrationClient) PostpaidInquiry(ctx context.Context, productCode, customerNumber, refID string) (*IntegrationResponse, error) {
	// Temporarily returning mock while waiting for gRPC regeneration
	return &IntegrationResponse{
		Success: true,
		Message: "Inquiry Success (Mock)",
	}, nil
}
