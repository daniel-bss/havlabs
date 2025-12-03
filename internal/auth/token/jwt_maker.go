package token

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJWTMaker() (*JWTMaker, error) {
	privateKeyData, err := os.ReadFile("./private.pem")
	if err != nil {
		// TODO:
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &JWTMaker{&privateKey.PublicKey, privateKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration, tokenType)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS384, payload)
	token, err := jwtToken.SignedString(maker.privateKey)
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, ErrInvalidToken
		}
		return maker.publicKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	err = payload.Valid(tokenType)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (maker *JWTMaker) PublicKey() *rsa.PublicKey {
	return maker.publicKey
}
