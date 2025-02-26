package columns

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.ColumnsOptions)
}

type weightOption []int

func (w weightOption) Apply(opts *options.ColumnsOptions) {
	opts.Weight = []int(w)
}

func Weight(weight ...int) Option {
	return weightOption(weight)
}
