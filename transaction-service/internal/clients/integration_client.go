package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IntegrationClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewIntegrationClient(baseURL string) *IntegrationClient {
	return &IntegrationClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
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
	url := fmt.Sprintf("%s/v1/transaction/topup", c.BaseURL)
	
	body, _ := json.Marshal(req)
	hReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	hReq.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(hReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result IntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return &result, nil
}
