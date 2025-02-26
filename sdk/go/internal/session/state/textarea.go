package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeTextArea WidgetType = "textArea"

type TextAreaState struct {
	ID           uuid.UUID
	Value        *string
	Label        string
	Placeholder  string
	DefaultValue *string
	Required     bool
	Disabled     bool
	MaxLength    *int32
	MinLength    *int32
	MaxLines     *int32
	MinLines     *int32
	AutoResize   bool
}

func (s *TextAreaState) IsWidgetState()      {}
func (s *TextAreaState) GetType() WidgetType { return WidgetTypeTextArea }
