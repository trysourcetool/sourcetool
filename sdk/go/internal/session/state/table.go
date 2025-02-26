package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeTable WidgetType = "table"

type TableState struct {
	ID           uuid.UUID
	Data         any
	Value        TableStateValue
	Header       string
	Description  string
	Height       *int32
	ColumnOrder  []string
	OnSelect     string
	RowSelection string
}

type TableStateValue struct {
	Selection *TableStateValueSelection
}

type TableStateValueSelection struct {
	Row  int32
	Rows []int32
}

func (s *TableState) IsWidgetState()      {}
func (s *TableState) GetType() WidgetType { return WidgetTypeTable }
