---
sidebar_position: 10
---

# Time Input

The Time Input widget provides a specialized input field for selecting time values with a time picker interface.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `string` | `None` | Current selected time value |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when no time is selected |
| `default_value` | `string` | `None` | Initial time value |
| `required` | `bool` | `False` | Whether a time selection is required |
| `disabled` | `bool` | `False` | Whether the time input is disabled |

## Examples

### Basic Time Input

```go
package main

import (
    "github.com/sourcetool/sourcetool-go"
    "github.com/sourcetool/sourcetool-go/timeinput"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic time input
        timeInput := ui.TimeInput("Start Time", timeinput.Placeholder("Select time"))
    }
}
```

### Disabled Time Input with Default Value

```go
// Create a disabled time input with a default value
closingTime := ui.TimeInput("Closing Time", timeinput.DefaultValue("18:00:00"), timeinput.Disabled(true))
```

## Related Components

- [DateInput](./date-input) - For selecting only a date without time
- [DateTimeInput](./date-time-input) - For selecting both date and time
