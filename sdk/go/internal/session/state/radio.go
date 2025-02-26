package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeRadio WidgetType = "radio"

type RadioState struct {
	ID           uuid.UUID
	Label        string
	Value        *int32
	Options      []string
	DefaultValue *int32
	Required     bool
	Disabled     bool
}

func (s *RadioState) IsWidgetState()      {}
func (s *RadioState) GetType() WidgetType { return WidgetTypeRadio }
