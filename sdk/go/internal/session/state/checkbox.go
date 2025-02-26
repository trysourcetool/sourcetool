package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeCheckbox WidgetType = "checkbox"

type CheckboxState struct {
	ID           uuid.UUID
	Label        string
	Value        bool
	DefaultValue bool
	Required     bool
	Disabled     bool
}

func (s *CheckboxState) IsWidgetState()      {}
func (s *CheckboxState) GetType() WidgetType { return WidgetTypeCheckbox }
