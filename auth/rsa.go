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
)

// RSAAuth handles RSA key operations for AddPay authentication
type RSAAuth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSAAuth creates a new RSA authentication handler
func NewRSAAuth(privateKeyPEM, publicKeyPEM []byte) (*RSAAuth, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &RSAAuth{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// Sign signs data using the private key
func (r *RSAAuth) Sign(data []byte) (string, error) {
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify verifies a signature using the public key
func (r *RSAAuth) Verify(data []byte, signature string) error {
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
func (r *RSAAuth) Encrypt(data []byte) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, data)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts data using the private key
func (r *RSAAuth) Decrypt(encryptedData string) ([]byte, error) {
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

// parsePrivateKey parses a PEM-encoded private key
func parsePrivateKey(privateKeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not an RSA private key")
		}
		return rsaKey, nil
	}
	return privateKey, nil
}

// parsePublicKey parses a PEM-encoded public key
func parsePublicKey(publicKeyPEM []byte) (*rsa.PublicKey, error) {
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
