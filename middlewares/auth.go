package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecretKey = "amat-sangat-rahasia"

func GenerateJWT(userID, email string) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"sub":   userID,                                // Subject (user ID)
		"email": email,                                 // User email
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expiration (24 hours)
		"iat":   time.Now().Unix(),                     // Issued at
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString([]byte(jwtSecretKey))
}
