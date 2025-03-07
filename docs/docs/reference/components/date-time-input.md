---
sidebar_position: 9
---

# Date Time Input

The Date Time Input widget provides a specialized input field for selecting both date and time with calendar and time picker interfaces.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `string` | `None` | Current selected date and time value |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when no date/time is selected |
| `default_value` | `string` | `None` | Initial date and time value |
| `required` | `bool` | `False` | Whether a date and time selection is required |
| `disabled` | `bool` | `False` | Whether the input is disabled |
| `format` | `string` | `"YYYY-MM-DD HH:mm:ss"` | Format string for date and time display |
| `max_value` | `string` | `None` | Maximum selectable date and time |
| `min_value` | `string` | `None` | Minimum selectable date and time |

## Examples

### Basic Date Time Input

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/datetimeinput"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic date time input
        dateTimeInput := ui.DateTimeInput("Date Time", datetimeinput.Placeholder("Select date and time"))
    }
}
```

### Date Time Input with Custom Format

```go
// Create a date time input with custom format
dateTimeInput := ui.DateTimeInput("Date Time", datetimeinput.Placeholder("Select date and time"), datetimeinput.Format("YYYY-MM-DD HH:mm"))
```

### Date Time Input with Range Constraints

```go
// Create a date time input with range constraints
dateTimeInput := ui.DateTimeInput("Date Time", datetimeinput.Placeholder("Select date and time"), datetimeinput.MinValue(time.Now()), datetimeinput.MaxValue(time.Now().AddDate(0, 0, 30)))
```

### Disabled Date Time Input with Default Value

```go
// Create a disabled date time input with a default value
dateTimeInput := ui.DateTimeInput("Date Time", datetimeinput.Placeholder("Select date and time"), datetimeinput.DefaultValue(time.Now()))
```

## Related Components

- [DateInput](./date-input) - For selecting only a date without time
- [TimeInput](./time-input) - For selecting only time without a date