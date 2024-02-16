package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your-secret-key") // Замените на ваш секретный ключ

// Claims - структура для хранения пользовательских данных в токене
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken - генерация JWT-токена
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен будет действителен 24 часа

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return claims, nil
}
