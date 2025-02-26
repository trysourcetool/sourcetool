package state

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const WidgetTypeDateInput WidgetType = "dateInput"

type DateInputState struct {
	ID           uuid.UUID
	Value        *time.Time
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Disabled     bool
	Format       string
	MaxValue     *time.Time
	MinValue     *time.Time
	Location     *time.Location
}

func (s *DateInputState) IsWidgetState()      {}
func (s *DateInputState) GetType() WidgetType { return WidgetTypeDateInput }
