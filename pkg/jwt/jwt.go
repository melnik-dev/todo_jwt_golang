package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Data struct {
	UserId   int
	TokenTTL time.Duration
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data Data) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": data.UserId,
		"exp":     time.Now().Add(data.TokenTTL).Unix(),
	})

	s, err := claims.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(token string) (bool, *Data) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}
	rawID, ok := claims["user_id"]
	if !ok {
		return false, nil
	}
	userIDFloat, ok := rawID.(float64)
	if !ok {
		return false, nil
	}
	return true, &Data{
		UserId: int(userIDFloat),
	}
}
