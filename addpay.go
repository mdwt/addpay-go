// Package addpay provides a Go client library for the AddPay payment processing API.
//
// The AddPay API allows merchants to integrate payment processing functionality
// including hosted checkout, tokenized payments, and debit check services.
//
// Basic usage:
//
//	config := &types.Config{
//		AppID:               "your-app-id",
//		GatewayURL:          "https://api.addpay.com",
//		MerchantPrivateKey:  merchantPrivateKeyPEM,
//		GatewayPublicKey:    gatewayPublicKeyPEM,
//	}
//
//	client, err := addpay.NewClient(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	checkoutReq := &types.CheckoutRequest{
//		MerchantNo:      "12345",
//		StoreNo:         "001",
//		MerchantOrderNo: "ORDER-001",
//		PriceCurrency:   "USD",
//		OrderAmount:     100.00,
//		NotifyURL:       "https://yoursite.com/notify",
//		ReturnURL:       "https://yoursite.com/return",
//	}
//
//	response, err := client.HostedCheckout(context.Background(), checkoutReq)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Checkout URL: %s\n", response.PayURL)
package addpay

import (
	"github.com/mdwt/addpay-go/client"
	"github.com/mdwt/addpay-go/logger"
	"github.com/mdwt/addpay-go/types"
)

// NewClient creates a new AddPay API client
func NewClient(config types.Config) (client.Client, error) {
	return client.New(config)
}

// NewDefaultLogger creates a new default logger that outputs JSON logs
func NewDefaultLogger() types.Logger {
	return logger.NewDefaultLogger()
}

// NewNoOpLogger creates a new no-op logger that discards all log messages
func NewNoOpLogger() types.Logger {
	return logger.NewNoOpLogger()
}

