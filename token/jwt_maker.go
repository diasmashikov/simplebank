package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

    // Define the JWT claims
    claims := jwt.MapClaims{}
    claims["authorized"] = true
    claims["user"] = payload.Username
    claims["exp"] = payload.ExpiredAt

    // Create the token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with the secret key
    return token.SignedString([]byte(maker.secretKey))
}
	
// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*jwt.MapClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, keyFunc)

	if err != nil {
        return nil, ErrInvalidToken
    }

	payload, ok := jwtToken.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}