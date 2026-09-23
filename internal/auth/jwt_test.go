package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := []byte("test-secret")
	userID := "42"

	token, err := GenerateToken(userID, secret)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	resultUserID, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if resultUserID != userID {
		t.Errorf("expected userID %s got %s", userID, resultUserID)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	secret := []byte("test-secret")

	claims := jwt.MapClaims{
		"sub": "42", 
		"exp": time.Now().Add(-1 * time.Minute).Unix(),
		"iat": time.Now().Add(-16 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("unexpected error signing token: %v", err)
	}

	_, err = ValidateToken(tokenString, secret)
	if err == nil {
		t.Error("expected error validating expired token, got nil")
	}
}