package options

import "time"

type TimeInputOptions struct {
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Disabled     bool
	Location     *time.Location
}
