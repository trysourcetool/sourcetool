package datetimeinput

import (
	"time"

	"github.com/trysourcetool/sourcetool-go/internal/options"
)

type Option interface {
	Apply(*options.DateTimeInputOptions)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.DateTimeInputOptions) {
	opts.Placeholder = string(p)
}

func Placeholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption time.Time

func (d defaultValueOption) Apply(opts *options.DateTimeInputOptions) {
	opts.DefaultValue = (*time.Time)(&d)
}

func DefaultValue(value time.Time) Option {
	return defaultValueOption(value)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.DateTimeInputOptions) {
	opts.Required = bool(r)
}

func Required(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.DateTimeInputOptions) {
	opts.Disabled = bool(d)
}

func Disabled(disabled bool) Option {
	return disabledOption(disabled)
}

type formatOption string

func (f formatOption) Apply(opts *options.DateTimeInputOptions) {
	opts.Format = string(f)
}

func Format(format string) Option {
	return formatOption(format)
}

type maxValueOption time.Time

func (m maxValueOption) Apply(opts *options.DateTimeInputOptions) {
	opts.MaxValue = (*time.Time)(&m)
}

func MaxValue(value time.Time) Option {
	return maxValueOption(value)
}

type minValueOption time.Time

func (m minValueOption) Apply(opts *options.DateTimeInputOptions) {
	opts.MinValue = (*time.Time)(&m)
}

func MinLength(value time.Time) Option {
	return minValueOption(value)
}

type locationOption time.Location

func (l locationOption) Apply(opts *options.DateTimeInputOptions) {
	opts.Location = (*time.Location)(&l)
}

func Location(location time.Location) Option {
	return locationOption(location)
}
