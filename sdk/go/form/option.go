package form

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.FormOptions)
}

type buttonDisabledOption bool

func (b buttonDisabledOption) Apply(opts *options.FormOptions) {
	opts.ButtonDisabled = bool(b)
}

func WithButtonDisabled(buttonDisabled bool) Option {
	return buttonDisabledOption(buttonDisabled)
}

type clearOnSubmitOption bool

func (c clearOnSubmitOption) Apply(opts *options.FormOptions) {
	opts.ClearOnSubmit = bool(c)
}

func WithClearOnSubmit(clearOnSubmit bool) Option {
	return clearOnSubmitOption(clearOnSubmit)
}
