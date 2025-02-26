package options

type SelectboxOptions struct {
	Label        string
	Options      []string
	Placeholder  string
	DefaultValue *string
	Required     bool
	Disabled     bool
	FormatFunc   func(string, int) string
}
