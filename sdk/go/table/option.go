package table

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.TableOptions)
}

type headerOption string

func (h headerOption) Apply(opts *options.TableOptions) {
	opts.Header = string(h)
}

func WithHeader(header string) Option {
	return headerOption(header)
}

type descriptionOption string

func (d descriptionOption) Apply(opts *options.TableOptions) {
	opts.Description = string(d)
}

func WithDescription(description string) Option {
	return descriptionOption(description)
}

type heightOption int32

func (h heightOption) Apply(opts *options.TableOptions) {
	opts.Height = (*int32)(&h)
}

func WithHeight(height int32) Option {
	return heightOption(height)
}

type columnOrderOption []string

func (c columnOrderOption) Apply(opts *options.TableOptions) {
	opts.ColumnOrder = []string(c)
}

func WithColumnOrder(order ...string) Option {
	return columnOrderOption(order)
}

type onSelectOption OnSelect

func (o onSelectOption) Apply(opts *options.TableOptions) {
	opts.OnSelect = OnSelect(o).String()
}

func WithOnSelect(behavior OnSelect) Option {
	return onSelectOption(behavior)
}

type rowSelectionOption RowSelection

func (r rowSelectionOption) Apply(opts *options.TableOptions) {
	opts.RowSelection = RowSelection(r).String()
}

func WithRowSelection(mode RowSelection) Option {
	return rowSelectionOption(mode)
}
