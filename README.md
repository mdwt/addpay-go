# AddPay Go Client Library

[![Go Reference](https://pkg.go.dev/badge/github.com/mdwt/addpay-go.svg)](https://pkg.go.dev/github.com/mdwt/addpay-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdwt/addpay-go)](https://goreportcard.com/report/github.com/mdwt/addpay-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go client library for the AddPay payment processing API. Supports hosted checkout, tokenized payments, and debit check functionality with RSA authentication.

## Installation

```bash
go get github.com/mdwt/addpay-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/mdwt/addpay-go"
    "github.com/mdwt/addpay-go/types"
)

func main() {
    config := types.Config{
        AppID:              "your-app-id",
        GatewayURL:         "https://api.paycloud.africa",
        MerchantPrivateKey: []byte(merchantPrivateKeyPEM),
        GatewayPublicKey:   []byte(gatewayPublicKeyPEM),
    }

    client, err := addpay.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    checkoutReq := types.CheckoutRequest{
        MerchantNo:      "MERCHANT001",
        StoreNo:         "STORE001",
        MerchantOrderNo: "ORDER-123",
        PriceCurrency:   "USD",
        OrderAmount:     99.99,
        NotifyURL:       "https://yoursite.com/webhook",
        ReturnURL:       "https://yoursite.com/success",
    }

    response, err := client.HostedCheckout(context.Background(), checkoutReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Checkout URL: %s\n", response.PayURL)
}
```

## API Methods

### Hosted Checkout
```go
response, err := client.HostedCheckout(ctx, types.CheckoutRequest{...})
```

### Query Token
```go
response, err := client.QueryToken(ctx, types.QueryTokenRequest{Token: "tok_123"})
```

### Tokenized Payment
```go
response, err := client.TokenizedPay(ctx, types.TokenizedPayRequest{...})
```

### Debit Check
```go
response, err := client.DebitCheck(ctx, types.DebitCheckRequest{...})
```

## Authentication

AddPay uses RSA key pairs. You need:
- **Merchant Private Key**: Your RSA private key (PEM format)
- **Gateway Public Key**: AddPay's public key (PEM format)  
- **App ID**: Your application identifier

## Configuration

```go
config := types.Config{
    AppID:              "your-app-id",           // Required
    GatewayURL:         "https://api.paycloud.africa", // Required
    MerchantPrivateKey: privateKeyPEM,          // Required
    GatewayPublicKey:   publicKeyPEM,           // Required
    Timeout:            30 * time.Second,       // Optional (default: 30s)
    Logger:             customLogger,           // Optional (default: JSON logger)
}
```

## Custom Logging

Implement the simple `Logger` interface:

```go
type Logger interface {
    Debug(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
}

// Built-in options
config.Logger = addpay.NewDefaultLogger() // JSON logger
config.Logger = addpay.NewNoOpLogger()    // Silent logger
```

The SDK automatically redacts sensitive data like tokens and account numbers in logs.

## Testing

```bash
go test ./...
```

## Examples

Complete examples are in the `examples/` directory:
- `hosted_checkout.go` - Hosted checkout flow
- `tokenized_payment.go` - Recurring payments
- `debit_check.go` - South African debit mandates

## License

MIT License