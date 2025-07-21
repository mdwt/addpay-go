package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mdwt/addpay-go"
	"github.com/mdwt/addpay-go/types"
)

func main() {
	// Example RSA keys (replace with your actual keys)
	merchantPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
your-merchant-private-key-here
-----END RSA PRIVATE KEY-----`

	gatewayPublicKey := `-----BEGIN PUBLIC KEY-----
your-gateway-public-key-here
-----END PUBLIC KEY-----`

	// Create client configuration
	config := types.Config{
		AppID:              "your-app-id",
		GatewayURL:         "https://api.paycloud.africa",
		MerchantPrivateKey: []byte(merchantPrivateKey),
		GatewayPublicKey:   []byte(gatewayPublicKey),
		Timeout:            30 * time.Second,
		Logger:             addpay.NewDefaultLogger(),
	}

	// Create the AddPay client
	client, err := addpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create AddPay client: %v", err)
	}

	// First, let's query a token to make sure it's valid
	token := "tok_1234567890abcdef" // This would be a token from a previous transaction

	fmt.Printf("🔍 Querying token: %s\n", token)
	tokenQuery := types.QueryTokenRequest{
		Token: token,
	}

	ctx := context.Background()
	tokenInfo, err := client.QueryToken(ctx, tokenQuery)
	if err != nil {
		log.Fatalf("Token query failed: %v", err)
	}

	fmt.Printf("✅ Token is valid!\n")
	fmt.Printf("💳 Card: %s (%s)\n", tokenInfo.TokenInfo.CardNumber, tokenInfo.TokenInfo.CardType)
	fmt.Printf("📅 Expires: %s\n", tokenInfo.TokenInfo.ExpiryDate)
	fmt.Printf("🟢 Status: %s\n", tokenInfo.TokenStatus)

	// Now process a tokenized payment
	fmt.Printf("\n💳 Processing tokenized payment...\n")

	paymentReq := types.TokenizedPayRequest{
		MerchantNo:      "MERCHANT001",                                 // Your merchant number
		StoreNo:         "STORE001",                                    // Your store number
		MerchantOrderNo: generatePaymentOrderNumber(),                  // Unique order number
		Token:           token,                                         // Payment token from previous transaction
		PriceCurrency:   "USD",                                         // Currency code
		OrderAmount:     29.99,                                         // Amount to charge
		NotifyURL:       "https://yourstore.com/webhook/addpay/notify", // Webhook URL
		Description:     "Monthly subscription renewal",                // Optional description
	}

	// Make the tokenized payment request
	paymentResp, err := client.TokenizedPay(ctx, paymentReq)
	if err != nil {
		log.Fatalf("Tokenized payment failed: %v", err)
	}

	// Display the payment result
	fmt.Printf("✅ Payment processed successfully!\n")
	fmt.Printf("🆔 Transaction ID: %s\n", paymentResp.TransactionID)
	fmt.Printf("📊 Status: %s\n", paymentResp.TransactionStatus)
	fmt.Printf("📋 Order Number: %s\n", paymentReq.MerchantOrderNo)
	fmt.Printf("💰 Amount: %.2f %s\n", paymentReq.OrderAmount, paymentReq.PriceCurrency)

	if paymentResp.TransactionStatus == "SUCCESS" {
		fmt.Printf("\n🎉 Payment completed successfully! You can now fulfill the order.\n")
	} else {
		fmt.Printf("\n⚠️  Payment status: %s - Check your webhook for updates.\n", paymentResp.TransactionStatus)
	}
}

func generatePaymentOrderNumber() string {
	return fmt.Sprintf("ORDER-%d", time.Now().Unix())
}
