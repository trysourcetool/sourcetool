package checkbox

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.CheckboxOptions)
}

type defaultValueOption bool

func (d defaultValueOption) Apply(opts *options.CheckboxOptions) {
	opts.DefaultValue = bool(d)
}

func DefaultValue(defaultValue bool) Option {
	return defaultValueOption(defaultValue)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.CheckboxOptions) {
	opts.Required = bool(r)
}

func Required(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.CheckboxOptions) {
	opts.Disabled = bool(d)
}

func Disabled(disabled bool) Option {
	return disabledOption(disabled)
}
