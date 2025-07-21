package types

import "time"

// Logger is a simple logging interface that can be implemented by any logger
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// Config represents the configuration for AddPay client
type Config struct {
	AppID              string
	GatewayURL         string
	MerchantPrivateKey []byte
	GatewayPublicKey   []byte
	Timeout            time.Duration
	Logger             Logger // Optional: uses default slog logger if nil
}


// CheckoutRequest represents a hosted checkout request
type CheckoutRequest struct {
	MerchantNo      string  `json:"merchant_no"`
	StoreNo         string  `json:"store_no"`
	MerchantOrderNo string  `json:"merchant_order_no"`
	PriceCurrency   string  `json:"price_currency"`
	OrderAmount     float64 `json:"order_amount"`
	Expires         int64   `json:"expires"`
	NotifyURL       string  `json:"notify_url"`
	ReturnURL       string  `json:"return_url"`
	Description     string  `json:"description,omitempty"`
	Geolocation     string  `json:"geolocation,omitempty"`
}

// CheckoutResponse represents the response from hosted checkout
type CheckoutResponse struct {
	PayURL string `json:"pay_url"`
}

// QueryTokenRequest represents a query token request
type QueryTokenRequest struct {
	Token string `json:"token"`
}

// QueryTokenResponse represents the response from query token
type QueryTokenResponse struct {
	TokenStatus string `json:"token_status"`
	TokenInfo   struct {
		CardNumber string `json:"card_number"`
		ExpiryDate string `json:"expiry_date"`
		CardType   string `json:"card_type"`
	} `json:"token_info"`
}

// TokenizedPayRequest represents a tokenized payment request
type TokenizedPayRequest struct {
	MerchantNo      string  `json:"merchant_no"`
	StoreNo         string  `json:"store_no"`
	MerchantOrderNo string  `json:"merchant_order_no"`
	Token           string  `json:"token"`
	PriceCurrency   string  `json:"price_currency"`
	OrderAmount     float64 `json:"order_amount"`
	NotifyURL       string  `json:"notify_url"`
	Description     string  `json:"description,omitempty"`
}

// TokenizedPayResponse represents the response from tokenized payment
type TokenizedPayResponse struct {
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
}

// DebitCheckRequest represents a debit check request
type DebitCheckRequest struct {
	MerchantNo      string  `json:"merchant_no"`
	StoreNo         string  `json:"store_no"`
	MerchantOrderNo string  `json:"merchant_order_no"`
	AccountNumber   string  `json:"account_number"`
	BankCode        string  `json:"bank_code"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	NotifyURL       string  `json:"notify_url"`
	Description     string  `json:"description,omitempty"`
}

// DebitCheckResponse represents the response from debit check
type DebitCheckResponse struct {
	MandateID     string `json:"mandate_id"`
	MandateStatus string `json:"mandate_status"`
}

// APIError represents an API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   APIError    `json:"error,omitempty"`
}
