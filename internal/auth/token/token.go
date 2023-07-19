package token

import (
	"errors"
	"fmt"
	"schooli-api/cmd/config"

	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 5

type Maker interface {
	Create(email string, id int64) (string, error)
	CreateRefresh(email string, id int64) (string, error)
	Verify(token string) (*JwtCustomClaims, error)
}

type tokenService struct {
	secretKey string
	cfg       *config.AppContext
}

func NewTokenMaker(secretKey string, cfg *config.AppContext) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &tokenService{secretKey: secretKey, cfg: cfg}, nil
}

// Different types of error returned by the VerifyToken function.
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	ID    int64  `json:"id"`
	jwt.RegisteredClaims
}

// CreateToken creates a new token for a specific username and duration.
func (maker *tokenService) Create(email string, id int64) (string, error) {
	exp := time.Now().Add(time.Hour * 3000)

	payload := &JwtCustomClaims{
		Email: email,
		ID:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: exp,
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not.
func (maker *tokenService) Verify(token string) (*JwtCustomClaims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtCustomClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*JwtCustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

// CreateRefreshToken creates a new token for a specific email and duration.
func (maker *tokenService) CreateRefresh(email string, id int64) (string, error) {
	exp := time.Now().Add(time.Hour * 30)

	payload := &JwtCustomClaims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: exp,
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}

	refresherToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return refresherToken.SignedString([]byte(maker.secretKey))
}
