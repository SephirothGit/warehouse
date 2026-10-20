package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func TestValidateToken_AlgNone(t *testing.T) {
	secret := []byte("test-secret")

	claims := jwt.MapClaims{
		"sub": "42",
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("unexpected error signing method: %v", err)
	}

	_, err = ValidateToken(tokenString, secret)
	if err == nil {
		t.Error("expected error validating alg none token, got nil")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken("42", []byte("correct-secret"))
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = ValidateToken(token, []byte("wrong-secret"))
	if err == nil {
		t.Error("expected error validating token with wrong secret, got nil")
	}
}

func TestValidateToken_WrongAlgorithm(t *testing.T) {
	secret := []byte("test-secret")

	claims := jwt.MapClaims{
		"sub": "42",
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("unexpected error signing token: %v", err)
	}

	_, err = ValidateToken(tokenString, secret)
	if err == nil {
		t.Error("expected error validating HS512 token, got nil")
	}
}

func TestValidateToken_NoExpiration(t *testing.T) {
	secret := []byte("test-secret")

	claims := jwt.MapClaims{"sub": "42"}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("unexpected error signing token: %v", err)
	}

	_, err = ValidateToken(tokenString, secret)
	if err == nil {
		t.Error("expected error validating token without exp, got nil")
	}
}