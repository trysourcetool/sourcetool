package table

type Value struct {
	Selection *Selection
}

type Selection struct {
	Row  int
	Rows []int
}

type SelectionBehavior string

const (
	SelectionBehaviorIgnore SelectionBehavior = "ignore"
	SelectionBehaviorRerun  SelectionBehavior = "rerun"
)

func (b SelectionBehavior) String() string {
	return string(b)
}

type SelectionMode string

const (
	SelectionModeSingle   SelectionMode = "single"
	SelectionModeMultiple SelectionMode = "multiple"
)

func (m SelectionMode) String() string {
	return string(m)
}
