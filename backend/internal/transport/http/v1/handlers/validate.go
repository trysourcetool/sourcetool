package handlers

import (
	"regexp"

	"github.com/go-playground/validator/v10"

	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func validateRequest(p any) error {
	v := validator.New()
	v.RegisterValidation("password", validatePassword)

	if err := v.Struct(p); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check for at least one letter
	hasLetter := false
	for _, c := range password {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return false
	}

	// Check for at least one digit
	hasDigit := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	// Check for valid characters only
	validChars := regexp.MustCompile(`^[a-zA-Z0-9!?_+*'"\` + "`" + `#$%&\-^\\@;:,./=~|[\](){}<>]+$`)
	return validChars.MatchString(password)
}
