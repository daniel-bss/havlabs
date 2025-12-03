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

	ID        uuid.UUID `json:"jti"`
	KeyID     string    `json:"kid"`
	Issuer    string    `json:"iss"`
	IssuedAt  int       `json:"iat"`
	ExpiredAt int       `json:"exp"`
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
		KeyID:     "oneandonly-jwk-until-further-implementation",
		Type:      tokenType,
		Username:  username,
		Role:      role,
		Issuer:    ISSUER,
		IssuedAt:  int(time.Now().Unix()),
		ExpiredAt: int(time.Now().Add(duration).Unix()),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid(tokenType TokenType) error {
	if payload.Type != tokenType {
		return ErrInvalidToken
	}

	expiredAt := time.Unix(payload.ExpiresAt.Unix(), 0)
	if time.Now().After(expiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: time.Unix(int64(payload.ExpiredAt), 0),
	}, nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: time.Unix(int64(payload.IssuedAt), 0),
	}, nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: time.Unix(int64(payload.IssuedAt), 0),
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
