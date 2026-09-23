package auth

import "testing"

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
