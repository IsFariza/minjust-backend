package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	IIN    string `json:"iin"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenStr string, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи токена")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	if claims.ExpiresAt == nil {
		return nil, errors.New("требуется истечение срока действия токена")
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("срок действия токена истек")
	}

	return claims, nil
}

func GenerateToken(userID int64, role string, iin string, secretKey string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		IIN:    iin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
