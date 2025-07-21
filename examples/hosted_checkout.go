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

	// Create a hosted checkout request
	checkoutReq := types.CheckoutRequest{
		MerchantNo:      "MERCHANT001",                                 // Your merchant number
		StoreNo:         "STORE001",                                    // Your store number
		MerchantOrderNo: generateOrderNumber(),                         // Unique order number
		PriceCurrency:   "USD",                                         // Currency code
		OrderAmount:     99.99,                                         // Amount to charge
		Expires:         time.Now().Add(24 * time.Hour).Unix(),         // Expiry timestamp
		NotifyURL:       "https://yourstore.com/webhook/addpay/notify", // Webhook URL
		ReturnURL:       "https://yourstore.com/checkout/success",      // Return URL after payment
		Description:     "Premium subscription purchase",               // Optional description
		Geolocation:     "US",                                          // Optional geolocation
	}

	// Make the hosted checkout request
	ctx := context.Background()
	response, err := client.HostedCheckout(ctx, checkoutReq)
	if err != nil {
		log.Fatalf("Hosted checkout failed: %v", err)
	}

	// Display the checkout URL
	fmt.Printf("✅ Hosted checkout created successfully!\n")
	fmt.Printf("🔗 Checkout URL: %s\n", response.PayURL)
	fmt.Printf("📋 Order Number: %s\n", checkoutReq.MerchantOrderNo)
	fmt.Printf("💰 Amount: %.2f %s\n", checkoutReq.OrderAmount, checkoutReq.PriceCurrency)

	// In a real application, you would redirect the user to response.PayURL
	fmt.Printf("\n🌐 Redirect your customer to the checkout URL to complete the payment.\n")
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORDER-%d", time.Now().Unix())
}
