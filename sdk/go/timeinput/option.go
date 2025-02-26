package timeinput

import (
	"time"

	"github.com/trysourcetool/sourcetool-go/internal/options"
)

type Option interface {
	Apply(*options.TimeInputOptions)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.TimeInputOptions) {
	opts.Placeholder = string(p)
}

func Placeholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption time.Time

func (d defaultValueOption) Apply(opts *options.TimeInputOptions) {
	opts.DefaultValue = (*time.Time)(&d)
}

func DefaultValue(value time.Time) Option {
	return defaultValueOption(value)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.TimeInputOptions) {
	opts.Required = bool(r)
}

func Required(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.TimeInputOptions) {
	opts.Disabled = bool(d)
}

func Disabled(disabled bool) Option {
	return disabledOption(disabled)
}

type locationOption time.Location

func (l locationOption) Apply(opts *options.TimeInputOptions) {
	opts.Location = (*time.Location)(&l)
}

func Location(location time.Location) Option {
	return locationOption(location)
}
