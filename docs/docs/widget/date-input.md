---
sidebar_position: 8
---

# Date Input

The Date Input widget provides a specialized input field for selecting dates with a calendar picker interface.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `string` | `None` | Current selected date value |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when no date is selected |
| `default_value` | `string` | `None` | Initial date value |
| `required` | `bool` | `False` | Whether a date selection is required |
| `disabled` | `bool` | `False` | Whether the date input is disabled |
| `format` | `string` | `"YYYY-MM-DD"` | Format string for date display |
| `max_value` | `string` | `None` | Maximum selectable date |
| `min_value` | `string` | `None` | Minimum selectable date |

## Examples

### Basic Date Input

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/dateinput"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic date input
        dateInput := ui.DateInput("Date", dateinput.Placeholder("Select a date"))
    }
}
```

### Date Input with Custom Format

```go
// Create a date input with custom format
dateInput := ui.DateInput("Date", dateinput.Placeholder("Select a date"), dateinput.Format("MM/DD/YYYY"))
```

### Date Input with Range Constraints

```go
// Create a date input with range constraints
dateInput := ui.DateInput("Date", dateinput.Placeholder("Select a date"), dateinput.MinValue(time.Now()), dateinput.MaxValue(time.Now().AddDate(0, 0, 30)))
```

### Disabled Date Input with Default Value

```go
// Create a disabled date input with a default value
dateInput := ui.DateInput("Date", dateinput.Placeholder("Select a date"), dateinput.Disabled(true))
```

## Related Components

- [DateTimeInput](./date-time-input) - For selecting both date and time
- [TimeInput](./time-input) - For selecting only time without a date
