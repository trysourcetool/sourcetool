package options

type TextInputOptions struct {
	Label        string
	Placeholder  string
	DefaultValue *string
	Required     bool
	Disabled     bool
	MaxLength    *int32
	MinLength    *int32
}
