package options

type NumberInputOptions struct {
	Label        string
	Placeholder  string
	DefaultValue *float64
	Required     bool
	Disabled     bool
	MaxValue     *float64
	MinValue     *float64
}
