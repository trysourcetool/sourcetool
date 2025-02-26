package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeTextInput WidgetType = "textInput"

type TextInputState struct {
	ID           uuid.UUID
	Value        *string
	Label        string
	Placeholder  string
	DefaultValue *string
	Required     bool
	Disabled     bool
	MaxLength    *int32
	MinLength    *int32
}

func (s *TextInputState) IsWidgetState()      {}
func (s *TextInputState) GetType() WidgetType { return WidgetTypeTextInput }
