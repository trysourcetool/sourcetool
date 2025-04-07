package numberinput

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.NumberInputOptions)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.NumberInputOptions) {
	opts.Placeholder = string(p)
}

func WithPlaceholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption float64

func (d defaultValueOption) Apply(opts *options.NumberInputOptions) {
	opts.DefaultValue = (*float64)(&d)
}

func WithDefaultValue(value float64) Option {
	return defaultValueOption(value)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.NumberInputOptions) {
	opts.Required = bool(r)
}

func WithRequired(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.NumberInputOptions) {
	opts.Disabled = bool(d)
}

func WithDisabled(disabled bool) Option {
	return disabledOption(disabled)
}

type maxValueOption float64

func (m maxValueOption) Apply(opts *options.NumberInputOptions) {
	opts.MaxValue = (*float64)(&m)
}

func WithMaxValue(value float64) Option {
	return maxValueOption(value)
}

type minValueOption float64

func (m minValueOption) Apply(opts *options.NumberInputOptions) {
	opts.MinValue = (*float64)(&m)
}

func WithMinValue(value float64) Option {
	return minValueOption(value)
}
