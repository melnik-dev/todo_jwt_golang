package crypto_test

import (
	"github.com/melnik-dev/go_todo_jwt/pkg/crypto"
	"testing"
)

func Test_HashPassword(t *testing.T) {
	password := "123456"
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hash == "" {
		t.Error("hash is empty")
	}

	if hash == password {
		t.Error("hash equals password")
	}
}

func Test_ComparePasswords(t *testing.T) {
	password := "123456"
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !crypto.ComparePasswords(password, hash) {
		t.Error("valid password did not match hash")
	}

	if crypto.ComparePasswords("invalid_password", hash) {
		t.Error("invalid password matched hash")
	}
}
