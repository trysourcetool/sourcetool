package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeNumberInput WidgetType = "numberInput"

type NumberInputState struct {
	ID           uuid.UUID
	Value        *float64
	Label        string
	Placeholder  string
	DefaultValue *float64
	Required     bool
	Disabled     bool
	MaxValue     *float64
	MinValue     *float64
}

func (s *NumberInputState) IsWidgetState()      {}
func (s *NumberInputState) GetType() WidgetType { return WidgetTypeNumberInput }
