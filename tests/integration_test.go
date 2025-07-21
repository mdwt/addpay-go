package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdwt/addpay-go"
	"github.com/mdwt/addpay-go/types"
)

// Mock RSA keys for testing
const (
	testMerchantPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEjWT2btYT9
5eSVfLRxGq8YYkPLrHQo5Zq7vLW3u7U+8pQ7fY6qjIcr4JnRhDGKsMLPGo3ckY6E
n5sGqgAKC2DHoYbGnWQHqOlZ9bG8L7gg/VfK+9QtSqXrSfFfzEOQqfMnZP3s8X5I
7OIWFxRkf9R7G3L8mAjwLqNzaHlhKgZdvFfF+QE6sG+kk8wXaGJk+XiHKmKjvO1I
+pHYFTUJVRsIo7TH9S/kF8M7XDT+l5wN8k/a3wCKUwHFJdFbPYGHgO8tHW/XUwQ
dL7+jKe+6ZqaFIa7JvGNgvDEo4JkSrO3HX6pLawIDAQABAoIBAQCJYTLQqoJ5hWq
vOC1Q8+O4qNYK9iJdKDa7+PiCqGvQ6SV7V+D8YRdj3QVnVf1s+OV6bKjp1j0m2T0
Qx7+Cq3mJb8k7BzZdD+6i0TQ+9jjVwS6Q1Xz7xXeZCp5jX6t7B4mZ4RJ7iH3lkJ
gHzJkzPP6wV4d2mN4WzgD3h6F0nI5YNrh6KNQ1xgZnJ7Js2RwjNrfh+4kJl4Pj2
5MXo6qFvXGaP3Kf6V8tNkJ2Y6l/h3Qd/9CyO0nJ+k7a8e9Qx1Q7qQmLWh5kzgH
3wJ3t9KkNt5wvkSj1V4Pz6EJnkfOh4wLKLH1k2nY5mHJ4wFq3MqZv1K7y5S1j8
zfQ0+P/lCq+tBAoGBAP6vQvfNQ/nK7j0a4D6Z5mJ6pP9Qr2Xq8YPT8pnQj9nM5M
LJl8nF6vN1gxN1G1UhZZJqcUKbGUo2KqDJlj5Y5r6aW2WK+CUa2FXg1qG8+ZnJ
zX4o5Qr2QsqXqfE3J9OE3lK5oF9HkNrg9r3kCCkO6qJ6GZqhJwB+5H6e4J5wQdB
AoGBAN4YNfpGg6LqZ9e9Y6N8hZ2qC0aGK8pN5x5r2W8n6nPd1pZ0mJ7S4v7R3X
8wI3k7aQE5C7o5N+WqKrQ0lV7F5oNe3tQ7aN2l9gvD1I8Vtz4IaZQb7YkF8n8
wJ5o4hN9qnx8yUa1Z3OQJL2ZsLU3SXr2JkGsN9Xh8QvLf4k3rOJ1hQAoGAJ5b
4YnvzxaIZJ2hL8DRrGLo9k4mFN3fDqTJJ2b+Qj6v9KGfKNz7v2O5B8Tx9q8k4V
gZ4hH3qT7h+rKl1kN6lB4mYJ1CtJ8s4jN3z+yT8qnF3+kFjqN9h4Y0wJ4N4sD
AoGBAKHV4TkOaZL2d6JnKHqUJK8Dt0jH5M8J2f+o3D4Y6kHXN6+Tr2Fhq8J7hL
nJkO3l1Jj8j8Z7W3J6qV4D5F6xOzDJZhEsrN6f2Q8K7dC1Y7q5B6WKNHjQf3
v9h0R2rDWNfONUUDfAJ6Z+RrKhGK5yQXCzJ3qTpKFhvJJdNNRhzBMAoGAJ2n
dIYnDq+CxJ6BZ/7HFo0vJqG+fMeF4JhV8G9F/qhqJvEgHJPq8z5vNNVhsWnJ
Y6nLsGK1KN4q6vOFAJ4D5z8jjSJFxDDvYqGqW2H7LPNDvTG6NXJVhF9wLF6J
+i9vNqA3+gOZ8hOY5zNNWFQAH1k4G5wqL3z9jM6rJeOWPLzMEQKBgQCnMzKu
8d8K1EBvfXhzVRfSXr6N4wJ5HgRzJ4cxTW6MhB9qKqKlvczI9QPqjv8jHjJz
GzfvkBZP2k7nVyh5rGkjJ4D4X7rqJ6rGq8yTFzGkOzjNXqBt0F5sHdZvbZQ
1H7pZ8k2f9OqTNDFc6vZtOJVhsNJ7Bk2kOV1TrGOyoYWCJ/8QwjJ3qQ==
-----END RSA PRIVATE KEY-----`

	testGatewayPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEjWT2btYT95eSVfLRxGq8YYkPLrHQo5Zq7vLW3u7U+8pQ7
fY6qjIcr4JnRhDGKsMLPGo3ckY6En5sGqgAKC2DHoYbGnWQHqOlZ9bG8L7gg/VfK
+9QtSqXrSfFfzEOQqfMnZP3s8X5I7OIWFxRkf9R7G3L8mAjwLqNzaHlhKgZdvFfF
+QE6sG+kk8wXaGJk+XiHKmKjvO1I+pHYFTUJVRsIo7TH9S/kF8M7XDT+l5wN8k/a
3wCKUwHFJdFbPYGHgO8tHW/XUwQdL7+jKe+6ZqaFIa7JvGNgvDEo4JkSrO3HX6pL
awIDAQAB
-----END PUBLIC KEY-----`
)

func TestHostedCheckoutIntegration(t *testing.T) {
	// Create a test server to mock the AddPay API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/entry/checkout" && r.Method == "POST" {
			// Mock successful checkout response
			response := map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"pay_url": "https://checkout.addpay.com/pay/12345",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.NotFound(w, r)
	}))
	defer testServer.Close()

	// Create client configuration
	config := types.Config{
		AppID:              "test-app-id",
		GatewayURL:         testServer.URL,
		MerchantPrivateKey: []byte(testMerchantPrivateKey),
		GatewayPublicKey:   []byte(testGatewayPublicKey),
		Timeout:            10 * time.Second,
		Logger:             addpay.NewNoOpLogger(),
	}

	// Create client
	client, err := addpay.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create checkout request with example payload
	checkoutReq := types.CheckoutRequest{
		MerchantNo:      "MERCHANT001",
		StoreNo:         "STORE001",
		MerchantOrderNo: "ORDER-" + time.Now().Format("20060102150405"),
		PriceCurrency:   "USD",
		OrderAmount:     99.99,
		Expires:         time.Now().Add(24 * time.Hour).Unix(),
		NotifyURL:       "https://yourstore.com/webhook/addpay/notify",
		ReturnURL:       "https://yourstore.com/checkout/success",
		Description:     "Test product purchase",
		Geolocation:     "US",
	}

	// Make the request
	ctx := context.Background()
	response, err := client.HostedCheckout(ctx, checkoutReq)
	if err != nil {
		t.Fatalf("HostedCheckout failed: %v", err)
	}

	// Verify response
	if response.PayURL == "" {
		t.Error("Expected PayURL to be set")
	}

	if response.PayURL != "https://checkout.addpay.com/pay/12345" {
		t.Errorf("Expected PayURL to be 'https://checkout.addpay.com/pay/12345', got '%s'", response.PayURL)
	}

	t.Logf("Checkout successful, PayURL: %s", response.PayURL)
}

func TestQueryTokenIntegration(t *testing.T) {
	// Create a test server to mock the AddPay API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/entry/query-token" && r.Method == "POST" {
			// Mock successful token query response
			response := map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"token_status": "ACTIVE",
					"token_info": map[string]interface{}{
						"card_number": "**** **** **** 1234",
						"expiry_date": "12/25",
						"card_type":   "VISA",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.NotFound(w, r)
	}))
	defer testServer.Close()

	// Create client configuration
	config := types.Config{
		AppID:              "test-app-id",
		GatewayURL:         testServer.URL,
		MerchantPrivateKey: []byte(testMerchantPrivateKey),
		GatewayPublicKey:   []byte(testGatewayPublicKey),
		Timeout:            10 * time.Second,
		Logger:             addpay.NewNoOpLogger(),
	}

	// Create client
	client, err := addpay.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create query token request
	tokenReq := types.QueryTokenRequest{
		Token: "tok_1234567890abcdef",
	}

	// Make the request
	ctx := context.Background()
	response, err := client.QueryToken(ctx, tokenReq)
	if err != nil {
		t.Fatalf("QueryToken failed: %v", err)
	}

	// Verify response
	if response.TokenStatus != "ACTIVE" {
		t.Errorf("Expected TokenStatus to be 'ACTIVE', got '%s'", response.TokenStatus)
	}

	if response.TokenInfo.CardNumber != "**** **** **** 1234" {
		t.Errorf("Expected CardNumber to be '**** **** **** 1234', got '%s'", response.TokenInfo.CardNumber)
	}

	t.Logf("Token query successful, Status: %s, Card: %s", response.TokenStatus, response.TokenInfo.CardNumber)
}

func TestTokenizedPayIntegration(t *testing.T) {
	// Create a test server to mock the AddPay API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/entry/tokenized-pay" && r.Method == "POST" {
			// Mock successful tokenized payment response
			response := map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"transaction_id":     "txn_1234567890",
					"transaction_status": "SUCCESS",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.NotFound(w, r)
	}))
	defer testServer.Close()

	// Create client configuration
	config := types.Config{
		AppID:              "test-app-id",
		GatewayURL:         testServer.URL,
		MerchantPrivateKey: []byte(testMerchantPrivateKey),
		GatewayPublicKey:   []byte(testGatewayPublicKey),
		Timeout:            10 * time.Second,
		Logger:             addpay.NewNoOpLogger(),
	}

	// Create client
	client, err := addpay.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create tokenized payment request
	payReq := types.TokenizedPayRequest{
		MerchantNo:      "MERCHANT001",
		StoreNo:         "STORE001",
		MerchantOrderNo: "ORDER-" + time.Now().Format("20060102150405"),
		Token:           "tok_1234567890abcdef",
		PriceCurrency:   "USD",
		OrderAmount:     49.99,
		NotifyURL:       "https://yourstore.com/webhook/addpay/notify",
		Description:     "Recurring subscription payment",
	}

	// Make the request
	ctx := context.Background()
	response, err := client.TokenizedPay(ctx, payReq)
	if err != nil {
		t.Fatalf("TokenizedPay failed: %v", err)
	}

	// Verify response
	if response.TransactionID == "" {
		t.Error("Expected TransactionID to be set")
	}

	if response.TransactionStatus != "SUCCESS" {
		t.Errorf("Expected TransactionStatus to be 'SUCCESS', got '%s'", response.TransactionStatus)
	}

	t.Logf("Tokenized payment successful, TransactionID: %s, Status: %s",
		response.TransactionID, response.TransactionStatus)
}

func TestDebitCheckIntegration(t *testing.T) {
	// Create a test server to mock the AddPay API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/entry/debit-check" && r.Method == "POST" {
			// Mock successful debit check response
			response := map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"mandate_id":     "mnd_1234567890",
					"mandate_status": "PENDING",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.NotFound(w, r)
	}))
	defer testServer.Close()

	// Create client configuration
	config := types.Config{
		AppID:              "test-app-id",
		GatewayURL:         testServer.URL,
		MerchantPrivateKey: []byte(testMerchantPrivateKey),
		GatewayPublicKey:   []byte(testGatewayPublicKey),
		Timeout:            10 * time.Second,
		Logger:             addpay.NewNoOpLogger(),
	}

	// Create client
	client, err := addpay.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create debit check request
	debitReq := types.DebitCheckRequest{
		MerchantNo:      "MERCHANT001",
		StoreNo:         "STORE001",
		MerchantOrderNo: "DEBIT-" + time.Now().Format("20060102150405"),
		AccountNumber:   "1234567890",
		BankCode:        "ABSA",
		Amount:          299.99,
		Currency:        "ZAR",
		NotifyURL:       "https://yourstore.com/webhook/addpay/notify",
		Description:     "Monthly subscription debit",
	}

	// Make the request
	ctx := context.Background()
	response, err := client.DebitCheck(ctx, debitReq)
	if err != nil {
		t.Fatalf("DebitCheck failed: %v", err)
	}

	// Verify response
	if response.MandateID == "" {
		t.Error("Expected MandateID to be set")
	}

	if response.MandateStatus != "PENDING" {
		t.Errorf("Expected MandateStatus to be 'PENDING', got '%s'", response.MandateStatus)
	}

	t.Logf("Debit check successful, MandateID: %s, Status: %s",
		response.MandateID, response.MandateStatus)
}

func TestClientConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  types.Config
		wantErr bool
		errMsg  string
	}{
		// Note: Can't test nil config anymore since we use value types
		{
			name: "missing app_id",
			config: types.Config{
				GatewayURL:         "https://api.addpay.com",
				MerchantPrivateKey: []byte(testMerchantPrivateKey),
				GatewayPublicKey:   []byte(testGatewayPublicKey),
			},
			wantErr: true,
			errMsg:  "app_id is required",
		},
		{
			name: "missing gateway_url",
			config: types.Config{
				AppID:              "test-app-id",
				MerchantPrivateKey: []byte(testMerchantPrivateKey),
				GatewayPublicKey:   []byte(testGatewayPublicKey),
			},
			wantErr: true,
			errMsg:  "gateway_url is required",
		},
		{
			name: "missing merchant private key",
			config: types.Config{
				AppID:            "test-app-id",
				GatewayURL:       "https://api.addpay.com",
				GatewayPublicKey: []byte(testGatewayPublicKey),
			},
			wantErr: true,
			errMsg:  "merchant_private_key is required",
		},
		{
			name: "missing gateway public key",
			config: types.Config{
				AppID:              "test-app-id",
				GatewayURL:         "https://api.addpay.com",
				MerchantPrivateKey: []byte(testMerchantPrivateKey),
			},
			wantErr: true,
			errMsg:  "gateway_public_key is required",
		},
		{
			name: "valid config",
			config: types.Config{
				AppID:              "test-app-id",
				GatewayURL:         "https://api.addpay.com",
				MerchantPrivateKey: []byte(testMerchantPrivateKey),
				GatewayPublicKey:   []byte(testGatewayPublicKey),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := addpay.NewClient(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewClient() expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("NewClient() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("NewClient() unexpected error = %v", err)
				}
			}
		})
	}
}
