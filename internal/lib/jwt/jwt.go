package jwt

import (
	"idea-store-auth/internal/domain/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	secret := os.Getenv("APP_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
