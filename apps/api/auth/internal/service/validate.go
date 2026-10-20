package service

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/pur1/to_on/packages/utils/apperror"
	"github.com/pur1/to_on/packages/utils/password"
)

var (
	emailRe    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`) // same rule as RegisterPage
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
)

const minPasswordLen = 8

func normalizeEmail(s string) string    { return strings.ToLower(strings.TrimSpace(s)) }
func normalizeUsername(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func validateRegister(in RegisterInput) error {
	if !emailRe.MatchString(in.Email) || len(in.Email) > 254 {
		return apperror.New(apperror.CodeInvalid, "invalid email format")
	}
	if !usernameRe.MatchString(in.Username) {
		return apperror.New(apperror.CodeInvalid, "username must be 3-30 letters, digits or underscores")
	}
	if utf8.RuneCountInString(in.Password) < minPasswordLen {
		return apperror.New(apperror.CodeInvalid, "password must be at least 8 characters")
	}
	if len(in.Password) > password.MaxLength {
		return apperror.New(apperror.CodeInvalid, "password must be at most 72 bytes")
	}
	return nil
}
