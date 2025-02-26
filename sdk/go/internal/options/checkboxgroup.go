package options

type CheckboxGroupOptions struct {
	Label        string
	Options      []string
	DefaultValue []string
	Required     bool
	Disabled     bool
	FormatFunc   func(string, int) string
}
