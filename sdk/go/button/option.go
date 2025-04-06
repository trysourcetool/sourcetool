package button

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.ButtonOptions)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.ButtonOptions) {
	opts.Disabled = bool(d)
}

func WithDisabled(disabled bool) Option {
	return disabledOption(disabled)
}
