package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const cost = 12

// MaxLength is bcrypt's input limit; longer passwords are rejected instead of silently truncated.
const MaxLength = 72

var ErrTooLong = errors.New("password exceeds 72 bytes")

func Hash(plain string) (string, error) {
	if len(plain) > MaxLength {
		return "", ErrTooLong
	}
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Compare(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
