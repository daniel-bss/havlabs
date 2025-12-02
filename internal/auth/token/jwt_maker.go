package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"reflect"
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
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &JWTMaker{&privateKey.PublicKey, privateKey}, nil
}

// type MyClaims struct {
// 	UserID string `json:"uid"`
// 	jwt.RegisteredClaims
// }

func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration, tokenType)
	if err != nil {
		return "", payload, err
	}
	// payload := MyClaims{
	// 	UserID: "123",
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		Subject:   "hehe",
	// 		Issuer:    "auth-service",
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
	// 		IssuedAt:  jwt.NewNumericDate(time.Now()),
	// 	},
	// }

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token, err := jwtToken.SignedString(maker.privateKey)
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	// func (maker *JWTMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			fmt.Println("masuk")
			fmt.Println(token.Method)
			fmt.Println(reflect.TypeOf(token.Method))
			return nil, ErrInvalidToken
		}
		return maker.publicKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		fmt.Println(">>!", err.Error())
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// payload, ok := jwtToken.Claims.(*Payload)
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok || !jwtToken.Valid {
		fmt.Println("KACAU")
		return nil, ErrInvalidToken
	}

	err = payload.Valid(tokenType)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
