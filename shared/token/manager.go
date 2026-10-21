package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	Issuer   = "daily-expense"
	Audience = "daily-expense-web"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret    []byte
	accessTTL time.Duration
}

func NewManager(secret string, accessTTL time.Duration) *Manager {
	return &Manager{
		secret:    []byte(secret),
		accessTTL: accessTTL,
	}
}

func (m *Manager) GenerateAccess(userID, email string) (string, int64, error) {
	now := time.Now()
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    Issuer,
			Audience:  jwt.ClaimStrings{Audience},
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, int64(m.accessTTL.Seconds()), nil
}

func (m *Manager) ParseAccess(raw string) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(Issuer),
		jwt.WithAudience(Audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	if claims.Subject == "" {
		return nil, fmt.Errorf("token tidak memiliki subject")
	}

	return claims, nil
}
