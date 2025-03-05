---
sidebar_position: 9
---

# Date Time Input

The Date Time Input widget provides a specialized input field for selecting both date and time with calendar and time picker interfaces.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `string` | `None` | Current selected date and time value |
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when no date/time is selected |
| `default_value` | `string` | `None` | Initial date and time value |
| `required` | `bool` | `False` | Whether a date and time selection is required |
| `disabled` | `bool` | `False` | Whether the input is disabled |
| `format` | `string` | `"YYYY-MM-DD HH:mm:ss"` | Format string for date and time display |
| `max_value` | `string` | `None` | Maximum selectable date and time |
| `min_value` | `string` | `None` | Minimum selectable date and time |

## Event Handling

The Date Time Input widget emits events when a date and time is selected:

```go
// Define a date time input with change handler
dateTimeInput := widget.NewDateTimeInput(ctx, widget.DateTimeInputOptions{
    Label: "Meeting Time",
    Format: "YYYY-MM-DD HH:mm",
    OnChange: func(dateTime string) {
        // Handle date and time selection
        fmt.Printf("Selected date and time: %s\n", dateTime)
    },
})
```

## Examples

### Basic Date Time Input

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic date time input
    dateTimeInput := widget.NewDateTimeInput(ctx, widget.DateTimeInputOptions{
        Label: "Event Start",
        Placeholder: "Select date and time",
        OnChange: func(dateTime string) {
            fmt.Printf("Event starts at: %s\n", dateTime)
        },
    })
    
    // Add the date time input to your UI
    container.Add(dateTimeInput)
}
```

### Date Time Input with Custom Format

```go
// Create a date time input with custom format
appointmentInput := widget.NewDateTimeInput(ctx, widget.DateTimeInputOptions{
    Label: "Appointment",
    Placeholder: "MM/DD/YYYY hh:mm A",
    Format: "MM/DD/YYYY hh:mm A",
    Use24HourFormat: false,
    ShowSeconds: false,
    OnChange: func(dateTime string) {
        fmt.Printf("Appointment scheduled for: %s\n", dateTime)
    },
})
```

### Date Time Input with Range Constraints

```go
// Create a date time input with range constraints
meetingScheduler := widget.NewDateTimeInput(ctx, widget.DateTimeInputOptions{
    Label: "Meeting Time",
    Placeholder: "Select meeting time",
    // Only allow dates from today to 14 days in the future
    MinDate: time.Now().Format("2006-01-02"),
    MaxDate: time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
    ShowNowButton: true,
    OnChange: func(dateTime string) {
        fmt.Printf("Meeting scheduled for: %s\n", dateTime)
    },
})
```

### Disabled Date Time Input with Default Value

```go
// Create a disabled date time input with a default value
deadlineInput := widget.NewDateTimeInput(ctx, widget.DateTimeInputOptions{
    Label: "Submission Deadline",
    DefaultValue: "2025-12-31 23:59:59", // New Year's Eve
    Disabled: true,
    Format: "MMMM D, YYYY HH:mm:ss", // Display as "December 31, 2025 23:59:59"
})
```

## Best Practices

1. Always provide a clear label and placeholder text to guide users
2. Use a date and time format that is familiar to your target audience
3. Set appropriate `min_date` and `max_date` constraints when applicable
4. Consider whether seconds are necessary for your use case
5. Choose between 12-hour and 24-hour time format based on user preferences
6. Provide calendar and time picker interfaces for easy selection
7. Allow keyboard input for users who prefer typing dates and times
8. Validate inputs to ensure they match the expected format
9. Consider localization requirements for international date and time formats

## Related Components

- [DateInput](./date-input) - For selecting only a date without time
- [TimeInput](./time-input) - For selecting only time without a date
- `DateRangeInput` - For selecting a range between two dates
