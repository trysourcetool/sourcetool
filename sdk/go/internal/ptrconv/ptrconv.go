package ptrconv

// StringValue returns the value of a string pointer.
// If the pointer is nil, returns an empty string.
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StringPtr returns a pointer to the string value.
// If the string is empty, returns nil.
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
