package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WalletClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewWalletClient(baseURL string) *WalletClient {
	return &WalletClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *WalletClient) Credit(amount float64, referenceID, referenceType string) error {
	// Based on wallet-service/cmd/main.go:
	// wallets.POST("/:id/credit", walletHandler.Credit)
	// Wait, the credit endpoint in wallet-service requires a wallet ID.
	// But in CommissionService, we only have userID? 
	// No, we need to know WHICH wallet to credit.
	// Usually, for commission, it's the staff's wallet.
	
	return fmt.Errorf("use CreditWallet instead with specific ID")
}

func (c *WalletClient) CreditWallet(walletID uint, amount float64, referenceID, referenceType string) error {
	url := fmt.Sprintf("%s/api/v1/wallets/%d/credit", c.BaseURL, walletID)
	
	payload := map[string]interface{}{
		"amount":         amount,
		"reference_id":   referenceID,
		"reference_type": referenceType,
	}
	
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Note: Internal service calls should ideally have a shared secret or be in a private network
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}
	
	return nil
}

func (c *WalletClient) DebitWallet(walletID uint, amount float64, referenceID, referenceType string) error {
	url := fmt.Sprintf("%s/api/v1/wallets/%d/debit", c.BaseURL, walletID)
	
	payload := map[string]interface{}{
		"amount":         amount,
		"reference_id":   referenceID,
		"reference_type": referenceType,
	}
	
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}
	
	return nil
}
