package dateinput

import (
	"time"

	"github.com/trysourcetool/sourcetool-go/internal/options"
)

type Option interface {
	Apply(*options.DateInputOptions)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.DateInputOptions) {
	opts.Placeholder = string(p)
}

func Placeholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption time.Time

func (d defaultValueOption) Apply(opts *options.DateInputOptions) {
	opts.DefaultValue = (*time.Time)(&d)
}

func DefaultValue(value time.Time) Option {
	return defaultValueOption(value)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.DateInputOptions) {
	opts.Required = bool(r)
}

func Required(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.DateInputOptions) {
	opts.Disabled = bool(d)
}

func Disabled(disabled bool) Option {
	return disabledOption(disabled)
}

type formatOption string

func (f formatOption) Apply(opts *options.DateInputOptions) {
	opts.Format = string(f)
}

func Format(format string) Option {
	return formatOption(format)
}

type maxValueOption time.Time

func (m maxValueOption) Apply(opts *options.DateInputOptions) {
	opts.MaxValue = (*time.Time)(&m)
}

func MaxValue(value time.Time) Option {
	return maxValueOption(value)
}

type minValueOption time.Time

func (m minValueOption) Apply(opts *options.DateInputOptions) {
	opts.MinValue = (*time.Time)(&m)
}

func MinLength(value time.Time) Option {
	return minValueOption(value)
}

type locationOption time.Location

func (l locationOption) Apply(opts *options.DateInputOptions) {
	opts.Location = (*time.Location)(&l)
}

func Location(location time.Location) Option {
	return locationOption(location)
}
