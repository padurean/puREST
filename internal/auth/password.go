package auth

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// Default admin credentials
const (
	DefaultAdminUser     = "admin"
	DefaultAdminPassword = "admin"
)

// HashAndSaltPassword ...
func HashAndSaltPassword(password string) (string, error) {
	hashedPasswordBytes, err :=
		bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	hashedPassword := string(hashedPasswordBytes)
	// log.Debug().Msgf("password: %s, hashed password (len = %d): %s",
	// 	password, len([]rune(hashedPassword)), hashedPassword)
	return hashedPassword, nil
}

// ComparePasswords ...
func ComparePasswords(plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		log.Error().Err(err).Msg("error comparing hashed and plain passwords")
		return false
	}
	return true
}

const minPasswordLen = 8
const maxPasswordLen = 32

// PasswordRequirementsMsg message used to inform the user about password strength requirements
var PasswordRequirementsMsg = fmt.Sprintf(
	"password must have between %d and %d characters of which at least "+
		"1 uppercase letter, 1 digit and 1 special character",
	minPasswordLen,
	maxPasswordLen,
)

// IsStrongPassword checks if the provided password meets the strength requirements
func IsStrongPassword(password string) error {
	err := errors.New(PasswordRequirementsMsg)
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		return err
	}
	var hasUpper bool
	var hasDigit bool
	var hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		default:
			return err
		}
	}
	if !hasUpper || !hasDigit || !hasSpecial {
		return err
	}
	return nil
}
