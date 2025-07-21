# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the AddPay Go Client Library - a comprehensive Go SDK for integrating with the AddPay payment processing API. It provides support for hosted checkout, tokenized payments, and debit check functionality with RSA-based authentication.

## Development Commands

### Build and Test
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests specifically
go test ./tests/

# Run tests with verbose output
go test -v ./...

# Build the library
go build ./...

# Format code
go fmt ./...

# Run static analysis (if golangci-lint is available)
golangci-lint run
```

### Running Examples
```bash
cd examples
go run hosted_checkout.go
go run tokenized_payment.go
go run debit_check.go
```

## Architecture Overview

The library follows clean architecture principles with clear separation of concerns:

### Core Components

- **`addpay.go`**: Main package exports and factory functions
- **`client/`**: HTTP client implementation with request/response handling
- **`types/`**: All data structures, interfaces, and API models
- **`auth/`**: RSA cryptographic operations for request signing
- **`logger/`**: Modern structured logging using Go's log/slog package

### Key Design Patterns

1. **slog-based Logging**: Uses Go 1.21+ `log/slog` for structured logging with JSON/text handlers
2. **RSA Authentication**: All API requests are signed using RSA private keys, responses verified with public keys
3. **Context Support**: All API methods accept `context.Context` for cancellation and timeouts
4. **Error Wrapping**: Uses Go 1.13+ error wrapping with `fmt.Errorf("...: %w", err)`

### Authentication Flow

1. Client signs request body with merchant private key (RSA-SHA256)
2. Signature sent in `X-Signature` header as base64-encoded string
3. AddPay verifies signature using merchant's public key
4. Response verification uses AddPay's public key

### Request/Response Pattern

- All API methods follow pattern: `func (c Client) Method(ctx context.Context, req types.Request) (types.Response, error)`
- Requests are JSON-marshaled and signed
- Responses may be wrapped in `APIResponse` format or direct unmarshaling
- Structured error handling with `types.APIError`

## Configuration Notes

### Required Configuration
- `AppID`: Application identifier from AddPay
- `GatewayURL`: API endpoint (e.g., "https://api.paycloud.africa")
- `MerchantPrivateKey`: RSA private key in PEM format for signing
- `GatewayPublicKey`: AddPay's public key in PEM format for verification

### Optional Configuration
- `Timeout`: HTTP client timeout (default: 30 seconds)
- `Logger`: slog.Logger instance (default: JSON logger at INFO level)

## Security Considerations

- Private keys are loaded into memory and used for signing
- No key rotation mechanism built-in
- Sensitive data is automatically redacted using slog.LogValuer interface
- HTTPS is required for all API communication

## Modern Logging Features

1. **slog Integration**: Uses Go's standard `log/slog` package for structured logging
2. **Context-Aware Logging**: All logging methods use `*Context` variants for request tracing
3. **Sensitive Data Redaction**: Implements `slog.LogValuer` interface for automatic redaction
4. **Multiple Output Formats**: Supports JSON (production) and text (development) formats
5. **Performance Optimized**: Uses slog's built-in performance optimizations

## Integration Patterns

### Basic Client Creation
```go
config := types.Config{
    AppID:              "your-app-id",
    GatewayURL:         "https://api.paycloud.africa",
    MerchantPrivateKey: privateKeyPEM,
    GatewayPublicKey:   publicKeyPEM,
    Logger:             addpay.NewDefaultLogger(), // JSON logger
}

client, err := addpay.NewClient(config)
```

### Logging Configuration Options
```go
// Production: JSON structured logging
logger := addpay.NewDefaultLogger()

// Development: Human-readable text logging
logger := addpay.NewDevelopmentLogger()

// Custom level JSON logger
logger := addpay.NewJSONLogger(addpay.LevelDebug)

// Silent logger for tests
logger := addpay.NewNoOpLogger()
```

### Error Handling
```go
response, err := client.HostedCheckout(ctx, request)
if err != nil {
    if apiErr, ok := err.(*types.APIError); ok {
        // Handle API-specific errors
        log.Printf("API Error: %s (Code: %s)", apiErr.Message, apiErr.Code)
    } else {
        // Handle network/other errors
        log.Printf("Request failed: %v", err)
    }
}
```

## Testing Strategy

- Integration tests in `tests/` directory
- Examples serve as functional tests
- Tests use `addpay.NewNoOpLogger()` to suppress log output
- No unit tests for individual components (integration-focused approach)

## Structured Logging Details

### Log Output Formats

**JSON (Production)**:
```json
{"time":"2025-01-20T10:30:00Z","level":"INFO","msg":"Creating hosted checkout","merchant_order_no":"ORDER-123","order_amount":99.99,"currency":"USD"}
```

**Text (Development)**:
```
time=2025-01-20T10:30:00.000Z level=INFO msg="Creating hosted checkout" merchant_order_no=ORDER-123 order_amount=99.99 currency=USD
```

### Sensitive Data Handling

The library automatically redacts sensitive information in logs:

- `types.Token`: Payment tokens are logged as `[REDACTED_TOKEN]`
- `types.SensitiveString`: General sensitive data logged as `[REDACTED]`

### Context-Aware Logging

All client methods use context-aware logging for better tracing:

```go
c.logger.InfoContext(ctx, "Creating hosted checkout",
    slog.String("merchant_order_no", req.MerchantOrderNo),
    slog.Float64("order_amount", req.OrderAmount))
```

# Go Struct vs Pointer Guidelines

## Rule: Use values by default, pointers only when needed

## When to use pointers:
1. **Modifying the struct** - methods that change fields
2. **Large structs** - roughly >100 bytes
3. **Nil is meaningful** - "not found" return values
4. **Interface requirements** - when interface demands pointer receiver

## Examples

### ✅ Good: Value receivers for read-only
```go
func (u User) String() string { return u.Name }
func (u User) IsValid() bool { return u.Email != "" }
func ProcessUser(user User) error { /* read only */ }
```

### ✅ Good: Pointer receivers only when modifying
```go
func (u *User) SetEmail(email string) { u.Email = email }
func (u *User) Save() error { /* modifies u.ID */ }
func UpdateUser(user *User) { user.LastSeen = time.Now() }
```

### ❌ Bad: Unnecessary pointers
```go
func (u *User) GetName() string { return u.Name } // Should be value receiver
func FormatUser(user *User) string { /* read only */ } // Should be value param
```

## Decision checklist:
- Does it modify the struct? → Use pointer
- Is struct >100 bytes? → Consider pointer
- Need nil as valid value? → Use pointer
- Otherwise → Use value
