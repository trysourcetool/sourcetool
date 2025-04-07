package textarea

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.TextAreaOptions)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.TextAreaOptions) {
	opts.Placeholder = string(p)
}

func WithPlaceholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption string

func (d defaultValueOption) Apply(opts *options.TextAreaOptions) {
	opts.DefaultValue = (*string)(&d)
}

func WithDefaultValue(value string) Option {
	return defaultValueOption(value)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.TextAreaOptions) {
	opts.Required = bool(r)
}

func WithRequired(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.TextAreaOptions) {
	opts.Disabled = bool(d)
}

func WithDisabled(disabled bool) Option {
	return disabledOption(disabled)
}

type maxLengthOption int32

func (m maxLengthOption) Apply(opts *options.TextAreaOptions) {
	opts.MaxLength = (*int32)(&m)
}

func WithMaxLength(length int32) Option {
	return maxLengthOption(length)
}

type minLengthOption int32

func (m minLengthOption) Apply(opts *options.TextAreaOptions) {
	opts.MinLength = (*int32)(&m)
}

func WithMinLength(length int32) Option {
	return minLengthOption(length)
}

type maxLinesOption int32

func (m maxLinesOption) Apply(opts *options.TextAreaOptions) {
	opts.MaxLines = (*int32)(&m)
}

func WithMaxLines(lines int32) Option {
	return maxLinesOption(lines)
}

type minLinesOption int32

func (m minLinesOption) Apply(opts *options.TextAreaOptions) {
	opts.MinLines = (*int32)(&m)
}

func WithMinLines(lines int32) Option {
	return minLinesOption(lines)
}

type autoResizeOption bool

func (a autoResizeOption) Apply(opts *options.TextAreaOptions) {
	opts.AutoResize = bool(a)
}

func WithAutoResize(autoResize bool) Option {
	return autoResizeOption(autoResize)
}
