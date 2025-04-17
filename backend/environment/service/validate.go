package service

import "regexp"

func validateSlug(s string) bool {
	pattern := `^[a-zA-Z0-9\-_]+$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}
