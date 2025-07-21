package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// RSAAuth handles RSA key operations for AddPay authentication
type RSAAuth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSAAuth creates a new RSA authentication handler
func NewRSAAuth(privateKeyPEM, publicKeyPEM []byte) (RSAAuth, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return RSAAuth{}, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return RSAAuth{}, fmt.Errorf("failed to parse public key: %w", err)
	}

	return RSAAuth{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// Sign signs data using the private key with SHA256WithRSA (matches Java SDK)
func (r RSAAuth) Sign(data []byte) (string, error) {
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// SignParameters signs request parameters using the Java SDK approach
func (r RSAAuth) SignParameters(params map[string]interface{}) (string, error) {
	// Filter out empty values and existing sign parameter
	filtered := filterParameters(params)

	// Create sorted parameter string for signing
	signString := createSignString(filtered)

	// Sign the string
	return r.Sign([]byte(signString))
}

// filterParameters removes empty values and 'sign' parameter (matches Java SDK paraFilter)
func filterParameters(params map[string]interface{}) map[string]string {
	filtered := make(map[string]string)

	for key, value := range params {
		// Skip sign parameter and empty values
		if key == "sign" || value == nil {
			continue
		}

		// Convert to string and check if not empty
		strValue := fmt.Sprintf("%v", value)
		if strValue != "" && strValue != "0" {
			filtered[key] = strValue
		}
	}

	return filtered
}

// createSignString creates a sorted parameter string for signing (matches Java SDK createLinkString)
func createSignString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	// Sort keys
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build URL-encoded parameter string
	values := url.Values{}
	for _, key := range keys {
		values.Add(key, params[key])
	}

	return values.Encode()
}

// Verify verifies a signature using the public key
func (r RSAAuth) Verify(data []byte, signature string) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	hash := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA256, hash[:], sig)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}
	return nil
}

// Encrypt encrypts data using the public key
func (r RSAAuth) Encrypt(data []byte) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, data)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts data using the private key
func (r RSAAuth) Decrypt(encryptedData string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	return decrypted, nil
}

// parsePrivateKey parses a private key (base64 PKCS8 or PEM-encoded)
func parsePrivateKey(privateKeyData []byte) (*rsa.PrivateKey, error) {
	if len(privateKeyData) == 0 {
		return nil, fmt.Errorf("private key is empty")
	}

	keyStr := strings.TrimSpace(string(privateKeyData))

	// Check if it's PEM format (starts with -----BEGIN)
	if strings.HasPrefix(keyStr, "-----BEGIN") {
		return parsePEMPrivateKey(privateKeyData)
	}

	// Otherwise, treat as base64-encoded PKCS8 (Java SDK format)
	return parseBase64PKCS8PrivateKey(keyStr)
}

// parsePEMPrivateKey parses a PEM-encoded private key
func parsePEMPrivateKey(privateKeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try PKCS8 format first
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not an RSA private key")
		}
		return rsaKey, nil
	}

	// Try PKCS1 format
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return privateKey, nil
}

// parseBase64PKCS8PrivateKey parses a base64-encoded private key (PKCS1 or PKCS8)
func parseBase64PKCS8PrivateKey(base64Key string) (*rsa.PrivateKey, error) {
	// Decode base64
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %w", err)
	}

	// Try PKCS1 format first (common format)
	privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err == nil {
		return privateKey, nil
	}

	// Try PKCS8 format
	key, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key (tried both PKCS1 and PKCS8): %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA private key")
	}
	return rsaKey, nil
}

// parsePublicKey parses a public key (base64 X.509 or PEM-encoded)
func parsePublicKey(publicKeyData []byte) (*rsa.PublicKey, error) {
	if len(publicKeyData) == 0 {
		return nil, fmt.Errorf("public key is empty")
	}

	keyStr := strings.TrimSpace(string(publicKeyData))

	// Check if it's PEM format (starts with -----BEGIN)
	if strings.HasPrefix(keyStr, "-----BEGIN") {
		return parsePEMPublicKey(publicKeyData)
	}

	// Otherwise, treat as base64-encoded X.509 (Java SDK format)
	return parseBase64X509PublicKey(keyStr)
}

// parsePEMPublicKey parses a PEM-encoded public key
func parsePEMPublicKey(publicKeyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA public key")
	}
	return rsaKey, nil
}

// parseBase64X509PublicKey parses a base64-encoded X.509 public key (Java SDK format)
func parseBase64X509PublicKey(base64Key string) (*rsa.PublicKey, error) {
	// Decode base64
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 public key: %w", err)
	}

	// Parse as X.509
	key, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X.509 public key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA public key")
	}
	return rsaKey, nil
}

// ensurePEMFormat adds PEM headers if they're missing
func ensurePEMFormat(keyData []byte, keyType string) []byte {
	keyStr := string(keyData)

	// If it already has PEM headers, return as-is
	if strings.HasPrefix(keyStr, "-----BEGIN") {
		return keyData
	}

	// Remove any whitespace and non-base64 characters
	keyStr = strings.TrimSpace(keyStr)
	// Remove any characters that aren't valid base64
	cleaned := ""
	for _, r := range keyStr {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') || r == '+' || r == '/' || r == '=' {
			cleaned += string(r)
		}
	}

	header := fmt.Sprintf("-----BEGIN %s-----", keyType)
	footer := fmt.Sprintf("-----END %s-----", keyType)

	// Insert line breaks every 64 characters for proper PEM format
	var formatted strings.Builder
	formatted.WriteString(header + "\n")

	for i := 0; i < len(cleaned); i += 64 {
		end := i + 64
		if end > len(cleaned) {
			end = len(cleaned)
		}
		formatted.WriteString(cleaned[i:end] + "\n")
	}

	formatted.WriteString(footer)
	return []byte(formatted.String())
}
