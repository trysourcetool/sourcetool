package textinput

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.TextInputOptions)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.TextInputOptions) {
	opts.Placeholder = string(p)
}

func WithPlaceholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption string

func (d defaultValueOption) Apply(opts *options.TextInputOptions) {
	opts.DefaultValue = (*string)(&d)
}

func WithDefaultValue(value string) Option {
	return defaultValueOption(value)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.TextInputOptions) {
	opts.Required = bool(r)
}

func WithRequired(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.TextInputOptions) {
	opts.Disabled = bool(d)
}

func WithDisabled(disabled bool) Option {
	return disabledOption(disabled)
}

type maxLengthOption int32

func (m maxLengthOption) Apply(opts *options.TextInputOptions) {
	opts.MaxLength = (*int32)(&m)
}

func WithMaxLength(length int32) Option {
	return maxLengthOption(length)
}

type minLengthOption int32

func (m minLengthOption) Apply(opts *options.TextInputOptions) {
	opts.MinLength = (*int32)(&m)
}

func WithMinLength(length int32) Option {
	return minLengthOption(length)
}
