package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeCheckboxGroup WidgetType = "checkboxGroup"

type CheckboxGroupState struct {
	ID           uuid.UUID
	Label        string
	Value        []int32
	Options      []string
	DefaultValue []int32
	Required     bool
	Disabled     bool
}

func (s *CheckboxGroupState) IsWidgetState()      {}
func (s *CheckboxGroupState) GetType() WidgetType { return WidgetTypeCheckboxGroup }
