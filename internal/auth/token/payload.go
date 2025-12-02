package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid hehe")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenType byte

const (
	TokenTypeAccessToken  = 1
	TokenTypeRefreshToken = 2
	ISSUER                = "havlabs"
)

// Payload contains the payload data of the token
type Payload struct {
	jwt.RegisteredClaims

	ID uuid.UUID `json:"jti"`
	// KeyID     string    `json:"kid"`
	Issuer    string    `json:"iss"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
	Type      TokenType `json:"token_type"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, role string, duration time.Duration, tokenType TokenType) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Type:      tokenType,
		Username:  username,
		Role:      role,
		Issuer:    ISSUER,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid(tokenType TokenType) error {
	if payload.Type != tokenType {
		return ErrInvalidToken
	}
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.ExpiredAt,
	}, nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetIssuer() (string, error) {
	return payload.Issuer, nil
}

func (payload *Payload) GetSubject() (string, error) {
	return "", nil
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
