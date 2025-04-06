package radio

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.RadioOptions)
}

type optionsOption []string

func (o optionsOption) Apply(opts *options.RadioOptions) {
	opts.Options = []string(o)
}

func WithOptions(options ...string) Option {
	return optionsOption(options)
}

type defaultValueOption string

func (d defaultValueOption) Apply(opts *options.RadioOptions) {
	opts.DefaultValue = (*string)(&d)
}

func WithDefaultValue(defaultValue string) Option {
	return defaultValueOption(defaultValue)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.RadioOptions) {
	opts.Required = bool(r)
}

func WithRequired(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.RadioOptions) {
	opts.Disabled = bool(d)
}

func WithDisabled(disabled bool) Option {
	return disabledOption(disabled)
}

type formatFuncOption func(string, int) string

func (f formatFuncOption) Apply(opts *options.RadioOptions) {
	opts.FormatFunc = func(s string, i int) string {
		return f(s, i)
	}
}

func WithFormatFunc(formatFunc func(string, int) string) Option {
	return formatFuncOption(formatFunc)
}
