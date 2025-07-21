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

	// Create a debit check request
	fmt.Printf("🏦 Creating debit check mandate...\n")

	debitReq := types.DebitCheckRequest{
		MerchantNo:      "MERCHANT001",                                 // Your merchant number
		StoreNo:         "STORE001",                                    // Your store number
		MerchantOrderNo: generateMandateNumber(),                       // Unique mandate reference
		AccountNumber:   "1234567890",                                  // Customer's bank account number
		BankCode:        "ABSA",                                        // Bank code (e.g., ABSA, FNB, Standard Bank)
		Amount:          199.99,                                        // Amount to debit
		Currency:        "ZAR",                                         // Currency (typically ZAR for South African banks)
		NotifyURL:       "https://yourstore.com/webhook/addpay/notify", // Webhook URL
		Description:     "Monthly insurance premium debit",             // Description for the mandate
	}

	// Make the debit check request
	ctx := context.Background()
	response, err := client.DebitCheck(ctx, debitReq)
	if err != nil {
		log.Fatalf("Debit check failed: %v", err)
	}

	// Display the mandate result
	fmt.Printf("✅ Debit check mandate created successfully!\n")
	fmt.Printf("🆔 Mandate ID: %s\n", response.MandateID)
	fmt.Printf("📊 Status: %s\n", response.MandateStatus)
	fmt.Printf("📋 Reference: %s\n", debitReq.MerchantOrderNo)
	fmt.Printf("🏦 Account: %s (%s)\n", debitReq.AccountNumber, debitReq.BankCode)
	fmt.Printf("💰 Amount: %.2f %s\n", debitReq.Amount, debitReq.Currency)
	fmt.Printf("📝 Description: %s\n", debitReq.Description)

	// Explain the next steps based on mandate status
	switch response.MandateStatus {
	case "PENDING":
		fmt.Printf("\n⏳ Mandate Status: PENDING\n")
		fmt.Printf("📞 The customer will receive a call or SMS to confirm the debit mandate.\n")
		fmt.Printf("🔔 You'll receive a webhook notification when the mandate is confirmed or rejected.\n")
	case "ACTIVE":
		fmt.Printf("\n✅ Mandate Status: ACTIVE\n")
		fmt.Printf("🎉 The mandate is active and ready for debiting!\n")
		fmt.Printf("💳 You can now process debit transactions using this mandate.\n")
	case "REJECTED":
		fmt.Printf("\n❌ Mandate Status: REJECTED\n")
		fmt.Printf("😞 The customer has rejected the debit mandate.\n")
		fmt.Printf("🔄 You may need to try a different payment method.\n")
	default:
		fmt.Printf("\n❓ Mandate Status: %s\n", response.MandateStatus)
		fmt.Printf("🔔 Check your webhook for status updates.\n")
	}

	fmt.Printf("\n📘 Note: Debit check mandates are commonly used in South Africa for recurring payments\n")
	fmt.Printf("    such as insurance premiums, subscriptions, and loan repayments.\n")
}

func generateMandateNumber() string {
	return fmt.Sprintf("MANDATE-%d", time.Now().Unix())
}
