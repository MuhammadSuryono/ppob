package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
	"github.com/yontech/ppob/integration-service/config"
)

type DigiflazzClient struct {
	cfg      *config.Config
	client   *http.Client
	breaker  *gobreaker.CircuitBreaker
}

func NewDigiflazzClient(cfg *config.Config) *DigiflazzClient {
	settings := gobreaker.Settings{
		Name:        "digiflazz",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 10 && failureRatio >= 0.5
		},
	}

	return &DigiflazzClient{
		cfg:     cfg,
		client:  &http.Client{Timeout: 30 * time.Second},
		breaker: gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *DigiflazzClient) generateSignature(payload map[string]interface{}) string {
	delete(payload, "sign")
	delete(payload, "username")

	jsonBytes, _ := json.Marshal(payload)
	hash := md5.Sum(jsonBytes)
	md5String := hex.EncodeToString(hash[:])

	signature := fmt.Sprintf("%s%s%s", c.cfg.DigiflazzKey, md5String, c.cfg.DigiflazzSecret)
	signatureHash := md5.Sum([]byte(signature))
	return hex.EncodeToString(signatureHash[:])
}

func (c *DigiflazzClient) generateWebhookSignature(payload string, timestamp string) string {
	message := payload + timestamp + c.cfg.DigiflazzSecret
	hash := md5.Sum([]byte(message))
	return hex.EncodeToString(hash[:])
}

func (c *DigiflazzClient) maskData(data map[string]interface{}) map[string]interface{} {
	masked := make(map[string]interface{})
	for k, v := range data {
		switch k {
		case "username", "sign", "key", "secret":
			masked[k] = "****"
		case "phone", "customer_no", "customer_id", "hp":
			if str, ok := v.(string); ok && len(str) > 8 {
				masked[k] = str[:4] + "****" + str[len(str)-4:]
			} else {
				masked[k] = "****"
			}
		default:
			masked[k] = v
		}
	}
	return masked
}

func (c *DigiflazzClient) doRequest(ctx context.Context, endpoint string, payload map[string]interface{}) ([]byte, error) {
	payload["username"] = c.cfg.DigiflazzKey
	payload["sign"] = c.generateSignature(payload)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	fullURL := c.cfg.DigiflazzURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Log Request
	maskedPayload := c.maskData(payload)
	log.Printf("[Digiflazz] Outbound Request: URL=%s, Payload=%v", fullURL, maskedPayload)

	result, err := c.breaker.Execute(func() (interface{}, error) {
		start := time.Now()
		resp, err := c.client.Do(req)
		duration := time.Since(start)

		if err != nil {
			log.Printf("[Digiflazz] Request Failed: URL=%s, Error=%v, Duration=%v", fullURL, err, duration)
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		// Log Response
		var respData interface{}
		_ = json.Unmarshal(body, &respData)
		
		var maskedResp interface{}
		if m, ok := respData.(map[string]interface{}); ok {
			maskedResp = c.maskData(m)
		} else {
			maskedResp = respData
		}
		
		log.Printf("[Digiflazz] Inbound Response: Status=%d, Duration=%v, Body=%v", resp.StatusCode, duration, maskedResp)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("digiflazz API error: status=%d, body=%s", resp.StatusCode, string(body))
		}

		return body, nil
	})

	if err != nil {
		return nil, fmt.Errorf("digiflazz request failed: %w", err)
	}

	return result.([]byte), nil
}

type PriceListRequest struct {
	Cmd string `json:"cmd"`
}

type PriceListResponse struct {
	Success bool                   `json:"success"`
	Data    []DigiflazzProduct     `json:"data"`
	Message string                 `json:"message"`
	Error   DigiflazzErrorResponse `json:"error"`
}

type DigiflazzProduct struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Provider     string  `json:"provider"`
	Status       string  `json:"status"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
	Brand        string  `json:"brand"`
}

type DigiflazzErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *DigiflazzClient) GetPriceList(ctx context.Context, productType string) (*PriceListResponse, error) {
	payload := map[string]interface{}{
		"cmd": productType,
	}

	body, err := c.doRequest(ctx, "/pricelist", payload)
	if err != nil {
		return nil, err
	}

	var response PriceListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return &response, fmt.Errorf("digiflazz error: %s", response.Error.Message)
	}

	return &response, nil
}

type TransactionRequest struct {
	Code       string `json:"code"`
	Phone      string `json:"phone"`
	RefID      string `json:"ref_id"`
	CustomerNo string `json:"customer_no,omitempty"`
}

type TransactionResponse struct {
	Success   bool                `json:"success"`
	Data      DigiflazzTransaction `json:"data"`
	Message   string              `json:"message"`
	Error     DigiflazzErrorResponse `json:"error"`
}

type DigiflazzTransaction struct {
	RefID      string `json:"ref_id"`
	TrxID      string `json:"trx_id"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	Price      string `json:"price"`
	ScCode     string `json:"sc_code"`
	ScMessage  string `json:"sc_message"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

func (c *DigiflazzClient) CreateTransaction(ctx context.Context, req *TransactionRequest) (*TransactionResponse, error) {
	payload := map[string]interface{}{
		"cmd":    "inq-pasca",
		"code":   req.Code,
		"phone":  req.Phone,
		"ref_id": req.RefID,
	}

	if req.CustomerNo != "" {
		payload["customer_no"] = req.CustomerNo
	}

	body, err := c.doRequest(ctx, "/transaction", payload)
	if err != nil {
		return nil, err
	}

	var response TransactionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

type TopUpRequest struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
	RefID string `json:"ref_id"`
}

func (c *DigiflazzClient) TopUp(ctx context.Context, req *TopUpRequest) (*TransactionResponse, error) {
	payload := map[string]interface{}{
		"cmd":    "topup",
		"code":   req.Code,
		"phone":  req.Phone,
		"ref_id": req.RefID,
	}

	body, err := c.doRequest(ctx, "/transaction", payload)
	if err != nil {
		return nil, err
	}

	var response TransactionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

type BalanceResponse struct {
	Success bool    `json:"success"`
	Data    Balance `json:"data"`
}

type Balance struct {
	Username   string  `json:"username"`
	Saldo      float64 `json:"saldo"`
	Admin      float64 `json:"admin"`
	Lifetime   float64 `json:"lifetime"`
	LastUpdate string  `json:"last_update"`
}

func (c *DigiflazzClient) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	payload := map[string]interface{}{
		"cmd": "deposit",
	}

	body, err := c.doRequest(ctx, "/profile", payload)
	if err != nil {
		return nil, err
	}

	var response BalanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

func (c *DigiflazzClient) VerifyWebhookSignature(payload string, timestamp string, signature string) bool {
	expected := c.generateWebhookSignature(payload, timestamp)
	return expected == signature
}

type DigiflazzRC int

const (
	RCSuccess       DigiflazzRC = 00
	RCPending       DigiflazzRC = 03
	RCInvalidPhone  DigiflazzRC = 02
	RCInsufficient  DigiflazzRC = 39
	RCSystemError   DigiflazzRC = 69
	RCTimeout       DigiflazzRC = 99
)

func MapRCToStatus(rc string) (string, string) {
	switch rc {
	case "00":
		return "success", "Transaction successful"
	case "03":
		return "pending", "Transaction is being processed"
	case "02":
		return "failed", "Invalid phone number or customer ID"
	case "39":
		return "failed", "Insufficient balance"
	case "69":
		return "failed", "System error, please try again"
	case "99":
		return "failed", "Transaction timeout"
	default:
		return "failed", "Unknown error"
	}
}