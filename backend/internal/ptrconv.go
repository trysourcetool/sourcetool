package internal

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

// BoolValue returns the value of a bool pointer.
// If the pointer is nil, returns false.
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// BoolPtr returns a pointer to the bool value.
// If the bool is false, returns nil.
func BoolPtr(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

// IntValue returns the value of an int pointer.
// If the pointer is nil, returns 0.
func IntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// IntPtr returns a pointer to the int value.
// If the int is 0, returns nil.
func IntPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
