package multiselect

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.MultiSelectOptions)
}

type optionsOption []string

func (o optionsOption) Apply(opts *options.MultiSelectOptions) {
	opts.Options = []string(o)
}

func WithOptions(options ...string) Option {
	return optionsOption(options)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.MultiSelectOptions) {
	opts.Placeholder = string(p)
}

func WithPlaceholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption []string

func (d defaultValueOption) Apply(opts *options.MultiSelectOptions) {
	opts.DefaultValue = []string(d)
}

func WithDefaultValue(defaultValue ...string) Option {
	return defaultValueOption(defaultValue)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.MultiSelectOptions) {
	opts.Required = bool(r)
}

func WithRequired(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.MultiSelectOptions) {
	opts.Disabled = bool(d)
}

func WithDisabled(disabled bool) Option {
	return disabledOption(disabled)
}

type formatFuncOption func(string, int) string

func (f formatFuncOption) Apply(opts *options.MultiSelectOptions) {
	opts.FormatFunc = f
}

func WithFormatFunc(formatFunc func(string, int) string) Option {
	return formatFuncOption(formatFunc)
}
