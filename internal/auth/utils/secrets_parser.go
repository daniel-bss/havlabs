package utils

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const PRIVATE_KEY_PATH = "./private.pem"
const PUBLIC_KEY_PATH = "./public.pem"

func ParseRSAPrivateKeyFromPEM() (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(PRIVATE_KEY_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to read Private Key file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA Private Key: %w", err)
	}
	return privateKey, nil
}

func ParseRSAPublicKeyFromPEM() (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(PUBLIC_KEY_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to read Public Key file: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA Public Key: %w", err)
	}
	return publicKey, nil
}
