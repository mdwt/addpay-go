package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/example/addpay-go/auth"
	"github.com/example/addpay-go/logger"
	"github.com/example/addpay-go/types"
)

// Client represents the AddPay API client
type Client struct {
	config     *types.Config
	httpClient *http.Client
	auth       *auth.RSAAuth
	logger     types.Logger
}

// New creates a new AddPay client
func New(config *types.Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.AppID == "" {
		return nil, fmt.Errorf("app_id is required")
	}

	if config.GatewayURL == "" {
		return nil, fmt.Errorf("gateway_url is required")
	}

	if len(config.MerchantPrivateKey) == 0 {
		return nil, fmt.Errorf("merchant_private_key is required")
	}

	if len(config.GatewayPublicKey) == 0 {
		return nil, fmt.Errorf("gateway_public_key is required")
	}

	// Set default timeout if not provided
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Set default logger if not provided
	if config.Logger == nil {
		config.Logger = logger.NewDefaultLogger(logger.INFO)
	}

	// Initialize RSA authentication
	rsaAuth, err := auth.NewRSAAuth(config.MerchantPrivateKey, config.GatewayPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RSA auth: %w", err)
	}

	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		auth:   rsaAuth,
		logger: config.Logger,
	}

	return client, nil
}

// HostedCheckout creates a hosted checkout request
func (c *Client) HostedCheckout(ctx context.Context, req *types.CheckoutRequest) (*types.CheckoutResponse, error) {
	c.logger.Info("Creating hosted checkout", types.Field{Key: "merchant_order_no", Value: req.MerchantOrderNo})

	response := &types.CheckoutResponse{}
	err := c.makeRequest(ctx, "POST", "/api/entry/checkout", req, response)
	if err != nil {
		c.logger.Error("Hosted checkout failed", types.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	c.logger.Info("Hosted checkout created successfully", types.Field{Key: "pay_url", Value: response.PayURL})
	return response, nil
}

// QueryToken queries token information
func (c *Client) QueryToken(ctx context.Context, req *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	c.logger.Info("Querying token", types.Field{Key: "token", Value: req.Token})

	response := &types.QueryTokenResponse{}
	err := c.makeRequest(ctx, "POST", "/api/entry/query-token", req, response)
	if err != nil {
		c.logger.Error("Query token failed", types.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	c.logger.Info("Token queried successfully", types.Field{Key: "status", Value: response.TokenStatus})
	return response, nil
}

// TokenizedPay processes a tokenized payment
func (c *Client) TokenizedPay(ctx context.Context, req *types.TokenizedPayRequest) (*types.TokenizedPayResponse, error) {
	c.logger.Info("Processing tokenized payment", types.Field{Key: "merchant_order_no", Value: req.MerchantOrderNo})

	response := &types.TokenizedPayResponse{}
	err := c.makeRequest(ctx, "POST", "/api/entry/tokenized-pay", req, response)
	if err != nil {
		c.logger.Error("Tokenized payment failed", types.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	c.logger.Info("Tokenized payment processed successfully",
		types.Field{Key: "transaction_id", Value: response.TransactionID},
		types.Field{Key: "status", Value: response.TransactionStatus})
	return response, nil
}

// DebitCheck creates a debit check request
func (c *Client) DebitCheck(ctx context.Context, req *types.DebitCheckRequest) (*types.DebitCheckResponse, error) {
	c.logger.Info("Creating debit check", types.Field{Key: "merchant_order_no", Value: req.MerchantOrderNo})

	response := &types.DebitCheckResponse{}
	err := c.makeRequest(ctx, "POST", "/api/entry/debit-check", req, response)
	if err != nil {
		c.logger.Error("Debit check failed", types.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	c.logger.Info("Debit check created successfully",
		types.Field{Key: "mandate_id", Value: response.MandateID},
		types.Field{Key: "status", Value: response.MandateStatus})
	return response, nil
}

// makeRequest makes an HTTP request to the AddPay API
func (c *Client) makeRequest(ctx context.Context, method, path string, request, response interface{}) error {
	// Marshal request body
	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	// Create HTTP request
	url := c.config.GatewayURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "addpay-go/1.0.0")
	req.Header.Set("X-App-ID", c.config.AppID)

	// Sign the request if we have a body
	if len(body) > 0 {
		signature, err := c.auth.Sign(body)
		if err != nil {
			return fmt.Errorf("failed to sign request: %w", err)
		}
		req.Header.Set("X-Signature", signature)
	}

	// Log request details
	c.logger.Debug("Making API request",
		types.Field{Key: "method", Value: method},
		types.Field{Key: "url", Value: url},
		types.Field{Key: "body_length", Value: len(body)})

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log response details
	c.logger.Debug("Received API response",
		types.Field{Key: "status_code", Value: resp.StatusCode},
		types.Field{Key: "body_length", Value: len(respBody)})

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		var apiResp types.APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err == nil && apiResp.Error != nil {
			return apiResp.Error
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response
	if response != nil {
		// Try to unmarshal as APIResponse first
		var apiResp types.APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err == nil {
			if !apiResp.Success && apiResp.Error != nil {
				return apiResp.Error
			}
			if apiResp.Data != nil {
				// Re-marshal the data field and unmarshal into our response type
				dataBytes, err := json.Marshal(apiResp.Data)
				if err != nil {
					return fmt.Errorf("failed to re-marshal data: %w", err)
				}
				if err := json.Unmarshal(dataBytes, response); err != nil {
					return fmt.Errorf("failed to unmarshal data: %w", err)
				}
			}
		} else {
			// Direct unmarshal into response type
			if err := json.Unmarshal(respBody, response); err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}
		}
	}

	return nil
}

// SetLogger allows changing the logger after client creation
func (c *Client) SetLogger(logger types.Logger) {
	c.logger = logger
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *types.Config {
	return c.config
}
