# AddPay Go Client Library Development Guidelines

This document provides guidelines and information for developers working on the AddPay Go client library.

## Core Principles
- Prefer value types (structs) over pointers unless there's a specific need for pointer semantics.

### Go Struct vs Pointer Guidelines

#### Rule: Use values by default, pointers only when needed

When to use pointers:
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


## Build/Configuration Instructions

### Prerequisites

- Go 1.19 or higher
- No external dependencies (uses only Go standard library)

### Installation

```bash
go get github.com/example/addpay-go
```

### Configuration

The library requires proper configuration to connect to the AddPay API:

```go
config := &types.Config{
    AppID:              "your-app-id",           // Required
    GatewayURL:         "https://api.paycloud.africa", // Required  
    MerchantPrivateKey: merchantPrivateKeyPEM,   // Required - RSA private key in PEM format
    GatewayPublicKey:   gatewayPublicKeyPEM,     // Required - RSA public key in PEM format
    Timeout:            30 * time.Second,        // Optional (default: 30s)
    Logger:             customLogger,            // Optional (default: stdout logger)
}

client, err := addpay.NewClient(config)
if err != nil {
    // Handle error
}
```

#### RSA Key Format

Keys must be in PEM format:

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

## Testing Information

### Running Tests

The project includes comprehensive tests that can be run using standard Go testing commands:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test file
go test -v ./tests/simple_test.go
```

### Test Structure

Tests are organized in the `tests/` directory and include:

1. **Integration Tests**: Tests that verify the client's interaction with the AddPay API using mock servers.
2. **Configuration Tests**: Tests that validate the client configuration.

### Creating New Tests

When adding new tests, follow these guidelines:

1. **Test File Location**: Place test files in the `tests/` directory.
2. **Test Naming**: Use descriptive names with the `Test` prefix (e.g., `TestHostedCheckout`).
3. **Mock Servers**: Use `httptest.NewServer` to create mock API endpoints.
4. **Test Keys**: Use the test RSA keys provided in the integration tests.

### Example Test

Here's a simple test that verifies the client rejects a nil configuration:

```go
package tests

import (
	"testing"

	"github.com/example/addpay-go"
)

func TestNilConfig(t *testing.T) {
	// Test with nil config
	_, err := addpay.NewClient(nil)
	
	// Expect an error
	if err == nil {
		t.Error("NewClient() with nil config should return an error")
	}
	
	// Verify error message
	expectedErrMsg := "config cannot be nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("NewClient() error = %v, want %v", err.Error(), expectedErrMsg)
	}
}
```

## Additional Development Information

### Code Style

- Follow standard Go code style and conventions.
- Use meaningful variable and function names.
- Add comments for exported functions and types.
- Keep functions focused on a single responsibility.

### Error Handling

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

### Logging

The library supports configurable logging:

```go
// Default logger (logs to stdout)
logger := addpay.NewDefaultLogger(addpay.INFO)

// No-op logger (discards all logs)
logger := addpay.NewNoOpLogger()

// Custom logger
type CustomLogger struct {}
func (l *CustomLogger) Debug(msg string, fields ...types.Field) { /* implementation */ }
func (l *CustomLogger) Info(msg string, fields ...types.Field)  { /* implementation */ }
func (l *CustomLogger) Warn(msg string, fields ...types.Field)  { /* implementation */ }
func (l *CustomLogger) Error(msg string, fields ...types.Field) { /* implementation */ }

// Use custom logger
config.Logger = &CustomLogger{}
```

### Security Best Practices

1. **Keep Private Keys Secure**: Never commit private keys to version control.
2. **Use Environment Variables**: Store keys in environment variables or secure vaults.
3. **Validate Webhooks**: Always verify webhook signatures.
4. **Use HTTPS**: Ensure all webhook URLs use HTTPS.
5. **Rotate Keys Regularly**: Follow AddPay's key rotation recommendations.

### Project Structure

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
└── README.md                   # Project documentation
```