package tools

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	key             string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func NewJWTService(key string, accessDuration, refreshDuration time.Duration) *JWTService {
	return &JWTService{key: key, AccessDuration: accessDuration, RefreshDuration: refreshDuration}
}

func (s *JWTService) generateJWT(data map[string]any, expiredAt *time.Time) (string, error) {
	claims := jwt.MapClaims{
		"exp": expiredAt.Unix(),
	}

	for k, v := range data {
		if k == "id" {
			claims["sub"] = v
		}

		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.key))
	if err != nil {
		return "", fmt.Errorf("не получилось подписать токен: %w", err)
	}

	return signedToken, nil
}

func (s *JWTService) GenerateAccessToken(data map[string]any) (string, error) {
	return s.generateJWT(data, new(time.Now().Add(s.AccessDuration)))
}

func (s *JWTService) GenerateRefreshToken(data map[string]any) (string, error) {
	return s.generateJWT(data, new(time.Now().Add(s.RefreshDuration)))
}

func (s *JWTService) ParseJWT(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод подписи: %v", token.Header["alg"])
		}
		return []byte(s.key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("не получилось получить данные из токена: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("невалидные данные внутри токена: %w", err)
	}

	if expiredAt, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(expiredAt) {
			return nil, fmt.Errorf("время жизни токена исчерпано")
		}
	}

	data := make(map[string]any)
	for k, v := range claims {
		data[k] = v
	}

	return data, nil
}
