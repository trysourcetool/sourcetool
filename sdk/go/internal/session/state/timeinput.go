package state

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const WidgetTypeTimeInput WidgetType = "timeInput"

type TimeInputState struct {
	ID           uuid.UUID
	Value        *time.Time
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Disabled     bool
	Location     *time.Location
}

func (s *TimeInputState) IsWidgetState()      {}
func (s *TimeInputState) GetType() WidgetType { return WidgetTypeTimeInput }
