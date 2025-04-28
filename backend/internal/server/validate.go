package server

import (
	"regexp"

	"github.com/go-playground/validator/v10"

	"github.com/trysourcetool/sourcetool/backend/internal/errdefs"
)

func validateRequest(p any) error {
	v := validator.New()

	if err := v.Struct(p); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	return nil
}

func validateSlug(s string) bool {
	pattern := `^[a-zA-Z0-9\-_]+$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}
