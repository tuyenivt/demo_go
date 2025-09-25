package token

import (
	"basic-fullstack/internal/logger"
	"basic-fullstack/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecret(logger logger.Logger) string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-for-dev"
		logger.Info("JWT_SECRET not set, using default development secret")
	} else {
		logger.Info("Using JWT_SECRET from environment")
	}
	return jwtSecret
}

func CreateJWT(user models.User, logger logger.Logger) string {
	jwtSecret := GetJWTSecret(logger)

	// Create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.Error("Failed to sign JWT", err)
		return ""
	}

	return tokenString
}

func ValidateJWT(tokenString string, logger logger.Logger) (*jwt.Token, error) {
	jwtSecret := GetJWTSecret(logger)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("Unexpected signing method", nil)
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		logger.Error("Failed to validate JWT", err)
		return nil, err
	}

	if !token.Valid {
		logger.Error("Invalid JWT token", nil)
		return nil, jwt.ErrTokenInvalidId
	}

	return token, nil
}
