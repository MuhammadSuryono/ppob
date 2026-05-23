package handlers

import (
	"context"

	"github.com/yontech/ppob/shared/proto/wallet"
	"github.com/yontech/ppob/wallet-service/internal/dto"
	"github.com/yontech/ppob/wallet-service/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletGRPCHandler struct {
	wallet.UnimplementedWalletServiceServer
	walletService *services.WalletService
}

func NewWalletGRPCHandler(walletService *services.WalletService) *WalletGRPCHandler {
	return &WalletGRPCHandler{
		walletService: walletService,
	}
}

func (h *WalletGRPCHandler) GetBalance(ctx context.Context, req *wallet.GetBalanceRequest) (*wallet.GetBalanceResponse, error) {
	resp, err := h.walletService.GetBalance(ctx, uint(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wallet not found: %v", err)
	}

	return &wallet.GetBalanceResponse{
		WalletId:   uint32(resp.WalletID),
		UserId:     uint32(resp.UserID),
		Balance:    resp.Balance,
		HoldAmount: resp.HoldAmount,
		Available:  resp.Available,
		Currency:   resp.Currency,
	}, nil
}

func (h *WalletGRPCHandler) PlaceHold(ctx context.Context, req *wallet.PlaceHoldRequest) (*wallet.PlaceHoldResponse, error) {
	holdReq := &dto.HoldRequest{
		Amount:        req.Amount,
		ReferenceID:   req.ReferenceId,
		ReferenceType: req.ReferenceType,
		ExpiresAt:     req.ExpiresAt,
	}

	resp, err := h.walletService.PlaceHold(ctx, uint(req.UserId), holdReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to place hold: %v", err)
	}

	return &wallet.PlaceHoldResponse{
		HoldId:        uint32(resp.HoldID),
		WalletId:      uint32(resp.WalletID),
		Amount:        resp.Amount,
		ReferenceId:   resp.ReferenceID,
		ReferenceType: resp.ReferenceType,
		Status:        resp.Status,
	}, nil
}

func (h *WalletGRPCHandler) ReleaseHold(ctx context.Context, req *wallet.ReleaseHoldRequest) (*wallet.ReleaseHoldResponse, error) {
	releaseReq := &dto.ReleaseHoldRequest{
		ReferenceID:   req.ReferenceId,
		ReferenceType: req.ReferenceType,
	}

	err := h.walletService.ReleaseHold(ctx, uint(req.UserId), releaseReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release hold: %v", err)
	}

	return &wallet.ReleaseHoldResponse{Success: true}, nil
}

func (h *WalletGRPCHandler) Debit(ctx context.Context, req *wallet.DebitRequest) (*wallet.DebitResponse, error) {
	debitReq := &dto.DebitRequest{
		Amount:        req.Amount,
		ReferenceID:   req.ReferenceId,
		ReferenceType: req.ReferenceType,
		ReleaseHold:   req.ReleaseHold,
	}

	err := h.walletService.Debit(ctx, uint(req.UserId), debitReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to debit: %v", err)
	}

	return &wallet.DebitResponse{Success: true}, nil
}

func (h *WalletGRPCHandler) Credit(ctx context.Context, req *wallet.CreditRequest) (*wallet.CreditResponse, error) {
	creditReq := &dto.CreditRequest{
		Amount:        req.Amount,
		ReferenceID:   req.ReferenceId,
		ReferenceType: req.ReferenceType,
	}

	err := h.walletService.Credit(ctx, uint(req.UserId), creditReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to credit: %v", err)
	}

	return &wallet.CreditResponse{Success: true}, nil
}

func (h *WalletGRPCHandler) Transfer(ctx context.Context, req *wallet.TransferRequest) (*wallet.TransferResponse, error) {
	err := h.walletService.Transfer(ctx, uint(req.FromUserId), uint(req.ToUserId), req.Amount, req.ReferenceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transfer: %v", err)
	}

	return &wallet.TransferResponse{Success: true}, nil
}
