package options

type TableOptions struct {
	Header       string
	Description  string
	Height       *int32
	ColumnOrder  []string
	OnSelect     string
	RowSelection string
}
