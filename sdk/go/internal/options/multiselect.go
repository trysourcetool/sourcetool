package options

type MultiSelectOptions struct {
	Label        string
	Options      []string
	Placeholder  string
	DefaultValue []string
	Required     bool
	Disabled     bool
	FormatFunc   func(string, int) string
}
