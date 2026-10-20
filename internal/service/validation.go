// This file implements password and email validation
package service

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

const (
	minPasswordRunes = 8
	maxPasswordBytes = 72 //bcrypt ignores and rejects anything longer
	maxEmailLen      = 254
)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateCredentials(email, password string) error {
	if email == "" || len(email) > maxEmailLen {
		return fmt.Errorf("%w: invalid email", ErrInvalidInput)
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return fmt.Errorf("%w: invalid email", ErrInvalidInput)
	}

	if utf8.RuneCountInString(password) < minPasswordRunes {
		return fmt.Errorf("%w: password must be at least %d characters", ErrInvalidInput, minPasswordRunes)
	}

	if len(password) > maxPasswordBytes {
		return fmt.Errorf("%w: password must be at most %d bytes", ErrInvalidInput, maxPasswordBytes)
	}
	return nil
}
