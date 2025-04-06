package table

type Value struct {
	Selection *Selection
}

type Selection struct {
	Row  int
	Rows []int
}

type OnSelect string

const (
	OnSelectIgnore OnSelect = "ignore"
	OnSelectRerun  OnSelect = "rerun"
)

func (o OnSelect) String() string {
	return string(o)
}

type RowSelection string

const (
	RowSelectionSingle   RowSelection = "single"
	RowSelectionMultiple RowSelection = "multiple"
)

func (r RowSelection) String() string {
	return string(r)
}
