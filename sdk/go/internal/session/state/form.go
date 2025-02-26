package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeForm WidgetType = "form"

type FormState struct {
	ID             uuid.UUID
	Value          bool
	ButtonLabel    string
	ButtonDisabled bool
	ClearOnSubmit  bool
}

func (s *FormState) IsWidgetState()      {}
func (s *FormState) GetType() WidgetType { return WidgetTypeForm }
