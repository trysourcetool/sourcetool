package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeButton WidgetType = "button"

type ButtonState struct {
	ID       uuid.UUID
	Value    bool
	Label    string
	Disabled bool
}

func (s *ButtonState) IsWidgetState()      {}
func (s *ButtonState) GetType() WidgetType { return WidgetTypeButton }
