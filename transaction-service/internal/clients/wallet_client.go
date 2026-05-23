package clients

import (
	"context"
	"fmt"

	"github.com/yontech/ppob/shared/proto/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WalletClient struct {
	client wallet.WalletServiceClient
	conn   *grpc.ClientConn
}

func NewWalletClient(address string) (*WalletClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to wallet service: %w", err)
	}

	return &WalletClient{
		client: wallet.NewWalletServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *WalletClient) Close() error {
	return c.conn.Close()
}

func (c *WalletClient) CreditWallet(ctx context.Context, userID uint, amount float64, referenceID, referenceType string) error {
	_, err := c.client.Credit(ctx, &wallet.CreditRequest{
		UserId:        uint32(userID),
		Amount:        amount,
		ReferenceId:   referenceID,
		ReferenceType: referenceType,
	})
	return err
}

func (c *WalletClient) DebitWallet(ctx context.Context, userID uint, amount float64, referenceID, referenceType string) error {
	_, err := c.client.Debit(ctx, &wallet.DebitRequest{
		UserId:        uint32(userID),
		Amount:        amount,
		ReferenceId:   referenceID,
		ReferenceType: referenceType,
	})
	return err
}

func (c *WalletClient) PlaceHoldForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error {
	_, err := c.client.PlaceHold(ctx, &wallet.PlaceHoldRequest{
		UserId:        uint32(userID),
		Amount:        amount,
		ReferenceId:   fmt.Sprintf("hold_tx_%s", transactionID),
		ReferenceType: "transaction",
	})
	return err
}

func (c *WalletClient) ReleaseHoldForTransaction(ctx context.Context, userID uint, transactionID string) error {
	_, err := c.client.ReleaseHold(ctx, &wallet.ReleaseHoldRequest{
		UserId:        uint32(userID),
		ReferenceId:   fmt.Sprintf("hold_tx_%s", transactionID),
		ReferenceType: "transaction",
	})
	return err
}

func (c *WalletClient) DebitForTransaction(ctx context.Context, userID uint, amount float64, transactionID string) error {
	_, err := c.client.Debit(ctx, &wallet.DebitRequest{
		UserId:        uint32(userID),
		Amount:        amount,
		ReferenceId:   fmt.Sprintf("debit_tx_%s", transactionID),
		ReferenceType: "transaction",
		ReleaseHold:   true,
	})
	return err
}
