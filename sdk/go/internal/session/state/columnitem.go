package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeColumnItem WidgetType = "columnItem"

type ColumnItemState struct {
	ID     uuid.UUID
	Weight float64
}

func (s *ColumnItemState) IsWidgetState()      {}
func (s *ColumnItemState) GetType() WidgetType { return WidgetTypeColumnItem }
