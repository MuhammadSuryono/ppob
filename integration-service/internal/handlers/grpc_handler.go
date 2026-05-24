package handlers

import (
	"context"

	"github.com/yontech/ppob/integration-service/internal/dto"
	"github.com/yontech/ppob/integration-service/internal/services"
	"github.com/yontech/ppob/shared/proto/integration"
)

type IntegrationGRPCHandler struct {
	integration.UnimplementedIntegrationServiceServer
	integrationService *services.IntegrationService
}

func NewIntegrationGRPCHandler(integrationService *services.IntegrationService) *IntegrationGRPCHandler {
	return &IntegrationGRPCHandler{
		integrationService: integrationService,
	}
}

func (h *IntegrationGRPCHandler) TopUp(ctx context.Context, req *integration.TopUpRequest) (*integration.TopUpResponse, error) {
	// UserID is not strictly needed for the integration call itself, 
	// but the service method signature expects it. 
	// For internal system calls, we can use a system user ID or 0.
	var userID uint = 0 

	resp, err := h.integrationService.InitiateDigiflazzTransaction(ctx, userID, &dto.DigiflazzTransactionRequest{
		ProductCode:    req.ProductCode,
		CustomerNumber: req.CustomerNumber,
		RefID:          req.RefId,
	})

	if err != nil {
		return &integration.TopUpResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &integration.TopUpResponse{
		Success: resp.Success,
		Message: resp.Message,
		Data: &integration.TopUpData{
			RefId:     resp.RefID,
			TrxId:     resp.TrxID,
			Status:    "pending", // Initial status from provider call
			ScCode:    resp.ScCode,
			ScMessage: resp.ScMessage,
		},
	}, nil
}
