# AddPay Go Client Library

[![Go Reference](https://pkg.go.dev/badge/github.com/example/addpay-go.svg)](https://pkg.go.dev/github.com/example/addpay-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/example/addpay-go)](https://goreportcard.com/report/github.com/example/addpay-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive Go client library for the AddPay payment processing API. This library provides a clean, idiomatic Go interface for integrating AddPay's payment services including hosted checkout, tokenized payments, and debit check functionality.

## Features

- 🚀 **Full API Coverage**: Supports all major AddPay API endpoints
- 🔐 **RSA Authentication**: Built-in RSA key management for secure API communication
- 🪵 **Configurable Logging**: Pluggable logger interface with default implementations
- 🧪 **Comprehensive Testing**: Extensive test suite with mock servers
- 📚 **Rich Examples**: Complete examples for all supported operations
- 🔄 **Context Support**: Full context.Context support for cancellation and timeouts
- 🏗️ **Clean Architecture**: Well-structured codebase following Go best practices

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [API Methods](#api-methods)
- [Examples](#examples)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Configuration](#configuration)
- [Error Handling](#error-handling)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/example/addpay-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/example/addpay-go"
    "github.com/example/addpay-go/types"
)

func main() {
    // Configure the client
    config := &types.Config{
        AppID:              "your-app-id",
        GatewayURL:         "https://api.paycloud.africa",
        MerchantPrivateKey: []byte(merchantPrivateKeyPEM),
        GatewayPublicKey:   []byte(gatewayPublicKeyPEM),
        Timeout:            30 * time.Second,
    }

    // Create client
    client, err := addpay.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // Create a hosted checkout
    checkoutReq := &types.CheckoutRequest{
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

## Authentication

AddPay uses RSA key pairs for authentication. You'll need:

1. **Merchant Private Key**: Your RSA private key for signing requests
2. **Gateway Public Key**: AddPay's public key for verifying responses
3. **App ID**: Your application identifier provided by AddPay

### Key Format

Keys should be in PEM format:

```
-----BEGIN RSA PRIVATE KEY-----
your-private-key-content
-----END RSA PRIVATE KEY-----
```

```
-----BEGIN PUBLIC KEY-----
gateway-public-key-content
-----END PUBLIC KEY-----
```

## API Methods

### Hosted Checkout

Create a hosted checkout session where customers complete payment on AddPay's secure pages.

```go
checkoutReq := &types.CheckoutRequest{
    MerchantNo:      "MERCHANT001",
    StoreNo:         "STORE001", 
    MerchantOrderNo: "ORDER-123",
    PriceCurrency:   "USD",
    OrderAmount:     99.99,
    Expires:         time.Now().Add(24 * time.Hour).Unix(),
    NotifyURL:       "https://yoursite.com/webhook",
    ReturnURL:       "https://yoursite.com/success",
    Description:     "Product purchase",
    Geolocation:     "US",
}

response, err := client.HostedCheckout(ctx, checkoutReq)
```

### Query Token

Query information about a payment token from previous transactions.

```go
tokenReq := &types.QueryTokenRequest{
    Token: "tok_1234567890abcdef",
}

response, err := client.QueryToken(ctx, tokenReq)
```

### Tokenized Payment

Process payments using previously stored payment tokens for recurring charges.

```go
paymentReq := &types.TokenizedPayRequest{
    MerchantNo:      "MERCHANT001",
    StoreNo:         "STORE001",
    MerchantOrderNo: "ORDER-124",
    Token:           "tok_1234567890abcdef",
    PriceCurrency:   "USD", 
    OrderAmount:     29.99,
    NotifyURL:       "https://yoursite.com/webhook",
    Description:     "Subscription renewal",
}

response, err := client.TokenizedPay(ctx, paymentReq)
```

### Debit Check

Create debit check mandates for South African bank account debits.

```go
debitReq := &types.DebitCheckRequest{
    MerchantNo:      "MERCHANT001",
    StoreNo:         "STORE001", 
    MerchantOrderNo: "MANDATE-125",
    AccountNumber:   "1234567890",
    BankCode:        "ABSA",
    Amount:          199.99,
    Currency:        "ZAR",
    NotifyURL:       "https://yoursite.com/webhook",
    Description:     "Monthly subscription",
}

response, err := client.DebitCheck(ctx, debitReq)
```

## Examples

The `examples/` directory contains complete working examples:

- **[hosted_checkout.go](examples/hosted_checkout.go)**: Complete hosted checkout flow
- **[tokenized_payment.go](examples/tokenized_payment.go)**: Token-based recurring payments  
- **[debit_check.go](examples/debit_check.go)**: South African debit check mandates

Run examples:

```bash
cd examples
go run hosted_checkout.go
```

## Project Structure

```
addpay-go/
├── addpay.go                    # Main package exports
├── client/
│   └── client.go               # HTTP client implementation
├── types/
│   └── types.go                # Type definitions and structures
├── auth/
│   └── rsa.go                  # RSA authentication handling
├── logger/
│   └── logger.go               # Logging interfaces and implementations
├── examples/
│   ├── hosted_checkout.go      # Hosted checkout example
│   ├── tokenized_payment.go    # Tokenized payment example
│   └── debit_check.go          # Debit check example
├── tests/
│   └── integration_test.go     # Integration tests
├── go.mod                      # Go module definition
├── .gitignore                  # Git ignore rules
└── README.md                   # This file
```

### Architecture Design

The library follows clean architecture principles:

- **`client/`**: Core HTTP client with request/response handling
- **`types/`**: All data structures and interfaces
- **`auth/`**: RSA cryptographic operations
- **`logger/`**: Pluggable logging system
- **`examples/`**: Usage demonstrations
- **`tests/`**: Comprehensive test suite

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests specifically
go test ./tests/

# Run tests with verbose output
go test -v ./...
```

The test suite includes:
- Unit tests for all components
- Integration tests with mock servers
- Configuration validation tests
- RSA authentication tests

## Configuration

### Basic Configuration

```go
config := &types.Config{
    AppID:              "your-app-id",           // Required
    GatewayURL:         "https://api.paycloud.africa", // Required  
    MerchantPrivateKey: merchantPrivateKeyPEM,   // Required
    GatewayPublicKey:   gatewayPublicKeyPEM,     // Required
    Timeout:            30 * time.Second,        // Optional (default: 30s)
    Logger:             customLogger,            // Optional (default: stdout logger)
}
```

### Custom Logger

Implement the `types.Logger` interface for custom logging:

```go
type CustomLogger struct {}

func (l *CustomLogger) Debug(msg string, fields ...types.Field) { /* implementation */ }
func (l *CustomLogger) Info(msg string, fields ...types.Field)  { /* implementation */ }
func (l *CustomLogger) Warn(msg string, fields ...types.Field)  { /* implementation */ }
func (l *CustomLogger) Error(msg string, fields ...types.Field) { /* implementation */ }

// Use custom logger
config.Logger = &CustomLogger{}
```

### Built-in Loggers

```go
// Default logger (logs to stdout)
logger := addpay.NewDefaultLogger(addpay.INFO)

// No-op logger (discards all logs)
logger := addpay.NewNoOpLogger()
```

## Error Handling

The library provides structured error handling:

```go
response, err := client.HostedCheckout(ctx, request)
if err != nil {
    // Check if it's an API error
    if apiErr, ok := err.(*types.APIError); ok {
        fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Message, apiErr.Code)
        if apiErr.Details != "" {
            fmt.Printf("Details: %s\n", apiErr.Details)
        }
    } else {
        // Network or other error
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}
```

### Common Error Types

- **Network Errors**: Connection issues, timeouts
- **Authentication Errors**: Invalid keys, signature verification failures  
- **API Errors**: Business logic errors returned by AddPay
- **Validation Errors**: Missing required fields, invalid values

## Security Best Practices

1. **Keep Private Keys Secure**: Never commit private keys to version control
2. **Use Environment Variables**: Store keys in environment variables or secure vaults
3. **Validate Webhooks**: Always verify webhook signatures
4. **Use HTTPS**: Ensure all webhook URLs use HTTPS
5. **Rotate Keys Regularly**: Follow AddPay's key rotation recommendations

## Dependencies

This library uses only Go standard library packages:

- `crypto/*` - RSA operations
- `net/http` - HTTP client
- `encoding/json` - JSON marshaling
- `context` - Request context
- `time` - Timeouts and timestamps

No external dependencies are required.

## Compatibility

- **Go Version**: 1.19 or higher
- **AddPay API**: Compatible with current AddPay Cloud API
- **Platforms**: All platforms supported by Go

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/example/addpay-go.git
cd addpay-go

# Run tests
go test ./...

# Run linting (if you have golangci-lint installed)
golangci-lint run

# Run examples
cd examples && go run hosted_checkout.go
```

## Support

- 📖 **Documentation**: [API Documentation](https://developers.paycloud.africa/docs/addpay/)
- 🐛 **Issues**: [GitHub Issues](https://github.com/example/addpay-go/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/example/addpay-go/discussions)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- AddPay team for providing comprehensive API documentation
- Go community for excellent tooling and libraries
- Contributors who help improve this library

---

**Note**: This library is not officially endorsed by AddPay. It's a community-maintained client library for easier integration with AddPay's payment services.