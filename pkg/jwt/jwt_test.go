package jwt_test

import (
	"github.com/melnik-dev/go_todo_jwt/pkg/jwt"
	"testing"
	"time"
)

func TestJWT_Create(t *testing.T) {
	const userId = 42
	jwtService := jwt.NewJWT("secret")
	token, err := jwtService.Create(jwt.Data{
		UserId:   userId,
		TokenTTL: time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is invalid")
	}
	if data.UserId != userId {
		t.Fatalf("User id %d not equal %d", data.UserId, userId)
	}
}

func TestJWT_Parse_Invalid(t *testing.T) {
	jwtService := jwt.NewJWT("secret")
	isValid, data := jwtService.Parse("invalid_token")
	if isValid || data != nil {
		t.Fatal("Expected fail for invalid token")
	}
}

func TestJWT_Parse_Expired(t *testing.T) {
	jwtService := jwt.NewJWT("secret")
	token, err := jwtService.Create(jwt.Data{
		UserId:   1,
		TokenTTL: -time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, _ := jwtService.Parse(token)
	if isValid {
		t.Fatal("Expected expired token to be invalid")
	}
}
