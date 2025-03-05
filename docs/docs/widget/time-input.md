---
sidebar_position: 10
---

# Time Input

The Time Input widget provides a specialized input field for selecting time values with a time picker interface.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `string` | `None` | Current selected time value |
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when no time is selected |
| `default_value` | `string` | `None` | Initial time value |
| `required` | `bool` | `False` | Whether a time selection is required |
| `disabled` | `bool` | `False` | Whether the time input is disabled |

## Event Handling

The Time Input widget emits events when a time is selected:

```go
// Define a time input with change handler
timeInput := widget.NewTimeInput(ctx, widget.TimeInputOptions{
    Label: "Meeting Time",
    Format: "HH:mm",
    OnChange: func(time string) {
        // Handle time selection
        fmt.Printf("Selected time: %s\n", time)
    },
})
```

## Examples

### Basic Time Input

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic time input
    timeInput := widget.NewTimeInput(ctx, widget.TimeInputOptions{
        Label: "Start Time",
        Placeholder: "Select time",
        OnChange: func(time string) {
            fmt.Printf("Start time: %s\n", time)
        },
    })
    
    // Add the time input to your UI
    container.Add(timeInput)
}
```

### Time Input with Custom Format

```go
// Create a time input with 12-hour format
appointmentTime := widget.NewTimeInput(ctx, widget.TimeInputOptions{
    Label: "Appointment Time",
    Placeholder: "hh:mm AM/PM",
    Format: "hh:mm A",
    Use24HourFormat: false,
    ShowSeconds: false,
    OnChange: func(time string) {
        fmt.Printf("Appointment time: %s\n", time)
    },
})
```

### Time Input with Range Constraints

```go
// Create a time input with range constraints
businessHours := widget.NewTimeInput(ctx, widget.TimeInputOptions{
    Label: "Business Hours",
    Placeholder: "Select time",
    // Only allow times between 9 AM and 5 PM
    MinTime: "09:00:00",
    MaxTime: "17:00:00",
    Step: 1800, // 30-minute increments (1800 seconds)
    OnChange: func(time string) {
        fmt.Printf("Selected business hour: %s\n", time)
    },
})
```

### Disabled Time Input with Default Value

```go
// Create a disabled time input with a default value
closingTime := widget.NewTimeInput(ctx, widget.TimeInputOptions{
    Label: "Closing Time",
    DefaultValue: "18:00:00", // 6 PM
    Disabled: true,
    Format: "HH:mm",
})
```

## Best Practices

1. Always provide a clear label and placeholder text to guide users
2. Use a time format that is familiar to your target audience
3. Set appropriate `min_time` and `max_time` constraints when applicable
4. Choose a suitable `step` value based on the required precision
5. Consider whether seconds are necessary for your use case
6. Choose between 12-hour and 24-hour time format based on user preferences
7. Provide a time picker interface for easy selection
8. Allow keyboard input for users who prefer typing times
9. Validate time inputs to ensure they match the expected format
10. Consider localization requirements for international time formats

## Related Components

- [DateInput](./date-input) - For selecting only a date without time
- [DateTimeInput](./date-time-input) - For selecting both date and time
- `DurationInput` - For selecting a time duration or interval
