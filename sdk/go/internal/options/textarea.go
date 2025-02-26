package options

type TextAreaOptions struct {
	Label        string
	Placeholder  string
	DefaultValue *string
	Required     bool
	Disabled     bool
	MaxLength    *int32
	MinLength    *int32
	MaxLines     *int32
	MinLines     *int32
	AutoResize   bool
}
