package state

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const WidgetTypeDateTimeInput WidgetType = "dateTimeInput"

type DateTimeInputState struct {
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

func (s *DateTimeInputState) IsWidgetState()      {}
func (s *DateTimeInputState) GetType() WidgetType { return WidgetTypeDateTimeInput }
