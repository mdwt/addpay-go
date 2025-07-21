package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mdwt/addpay-go"
	"github.com/mdwt/addpay-go/types"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Check environment variables
	if os.Getenv("APP_ID") == "" || os.Getenv("ENDPOINT") == "" ||
		os.Getenv("APP_RSA_PRIVATE_KEY_PKCS1") == "" || os.Getenv("GATEWAY_RSA_PUBLIC_KEY") == "" {
		log.Fatal("Required environment variables not set")
	}

	fmt.Println("Environment variables found:")
	fmt.Printf("APP_ID: %s\n", os.Getenv("APP_ID"))
	fmt.Printf("ENDPOINT: %s\n", os.Getenv("ENDPOINT"))
	fmt.Printf("MERCHANT_NO: %s\n", os.Getenv("MERCHANT_NO"))
	fmt.Printf("STORE_NO: %s\n", os.Getenv("STORE_NO"))

	// Debug key lengths
	privateKey := os.Getenv("APP_RSA_PRIVATE_KEY_PKCS1")
	publicKey := os.Getenv("GATEWAY_RSA_PUBLIC_KEY")

	fmt.Printf("Private key length: %d\n", len(privateKey))
	fmt.Printf("Public key length: %d\n", len(publicKey))
	fmt.Printf("Private key first 100 chars: %s...\n", privateKey[:min(100, len(privateKey))])
	fmt.Printf("Private key last 100 chars: ...%s\n", privateKey[max(0, len(privateKey)-100):])

	// Create client configuration with detailed logging
	// Try HTTPS endpoint
	endpoint := os.Getenv("ENDPOINT")
	if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.Replace(endpoint, "http://", "https://", 1)
		fmt.Printf("Trying HTTPS endpoint: %s\n", endpoint)
	}

	config := types.Config{
		AppID:              os.Getenv("APP_ID"),
		GatewayURL:         endpoint,
		MerchantPrivateKey: []byte(privateKey),
		GatewayPublicKey:   []byte(publicKey),
		Timeout:            30 * time.Second,
		Logger:             addpay.NewDefaultLogger(), // Use detailed logging
	}

	// Create client
	client, err := addpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Client created successfully!")

	// Create a simple checkout request
	checkoutReq := types.CheckoutRequest{
		MerchantNo:      os.Getenv("MERCHANT_NO"),
		StoreNo:         os.Getenv("STORE_NO"),
		MerchantOrderNo: "DEBUG-" + time.Now().Format("20060102150405"),
		PriceCurrency:   "ZAR",
		OrderAmount:     1.00, // Small amount for testing
		Expires:         time.Now().Add(1 * time.Hour).Unix(),
		NotifyURL:       "https://httpbin.org/post", // Test webhook URL
		ReturnURL:       "https://example.com/success",
		Description:     "Debug test payment",
		Geolocation:     "ZA",
	}

	fmt.Printf("Making checkout request with order: %s\n", checkoutReq.MerchantOrderNo)

	// Make the request
	ctx := context.Background()
	response, err := client.HostedCheckout(ctx, checkoutReq)
	if err != nil {
		log.Printf("HostedCheckout failed: %v", err)
		return
	}

	fmt.Printf("Success! PayURL: %s\n", response.PayURL)
}
