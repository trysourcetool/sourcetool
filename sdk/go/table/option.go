package table

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.TableOptions)
}

type headerOption string

func (h headerOption) Apply(opts *options.TableOptions) {
	opts.Header = string(h)
}

func Header(header string) Option {
	return headerOption(header)
}

type descriptionOption string

func (d descriptionOption) Apply(opts *options.TableOptions) {
	opts.Description = string(d)
}

func Description(description string) Option {
	return descriptionOption(description)
}

type heightOption int32

func (h heightOption) Apply(opts *options.TableOptions) {
	opts.Height = (*int32)(&h)
}

func Height(height int32) Option {
	return heightOption(height)
}

type columnOrderOption []string

func (c columnOrderOption) Apply(opts *options.TableOptions) {
	opts.ColumnOrder = []string(c)
}

func ColumnOrder(order ...string) Option {
	return columnOrderOption(order)
}

type onSelectOption SelectionBehavior

func (o onSelectOption) Apply(opts *options.TableOptions) {
	opts.OnSelect = SelectionBehavior(o).String()
}

func OnSelect(behavior SelectionBehavior) Option {
	return onSelectOption(behavior)
}

type rowSelectionOption SelectionMode

func (r rowSelectionOption) Apply(opts *options.TableOptions) {
	opts.RowSelection = SelectionMode(r).String()
}

func RowSelection(mode SelectionMode) Option {
	return rowSelectionOption(mode)
}
