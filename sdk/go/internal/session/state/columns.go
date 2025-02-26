package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeColumns WidgetType = "columns"

type ColumnsState struct {
	ID      uuid.UUID
	Columns int
}

func (s *ColumnsState) IsWidgetState()      {}
func (s *ColumnsState) GetType() WidgetType { return WidgetTypeColumns }
