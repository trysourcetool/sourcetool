---
sidebar_position: 8
---

# Date Input

The Date Input widget provides a specialized input field for selecting dates with a calendar picker interface.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `string` | `None` | Current selected date value |
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when no date is selected |
| `default_value` | `string` | `None` | Initial date value |
| `required` | `bool` | `False` | Whether a date selection is required |
| `disabled` | `bool` | `False` | Whether the date input is disabled |
| `format` | `string` | `"YYYY-MM-DD"` | Format string for date display |
| `max_value` | `string` | `None` | Maximum selectable date |
| `min_value` | `string` | `None` | Minimum selectable date |

## Event Handling

The Date Input widget emits events when a date is selected:

```go
// Define a date input with change handler
dateInput := widget.NewDateInput(ctx, widget.DateInputOptions{
    Label: "Birth Date",
    Format: "YYYY-MM-DD",
    OnChange: func(date string) {
        // Handle date selection
        fmt.Printf("Selected date: %s\n", date)
    },
})
```

## Examples

### Basic Date Input

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic date input
    dateInput := widget.NewDateInput(ctx, widget.DateInputOptions{
        Label: "Event Date",
        Placeholder: "Select a date",
        OnChange: func(date string) {
            fmt.Printf("Event date: %s\n", date)
        },
    })
    
    // Add the date input to your UI
    container.Add(dateInput)
}
```

### Date Input with Custom Format

```go
// Create a date input with custom format
birthdayInput := widget.NewDateInput(ctx, widget.DateInputOptions{
    Label: "Birthday",
    Placeholder: "MM/DD/YYYY",
    Format: "MM/DD/YYYY",
    OnChange: func(date string) {
        fmt.Printf("Birthday: %s\n", date)
    },
})
```

### Date Input with Range Constraints

```go
// Create a date input with range constraints
appointmentDate := widget.NewDateInput(ctx, widget.DateInputOptions{
    Label: "Appointment Date",
    Placeholder: "Select a date",
    // Only allow dates from today to 30 days in the future
    MinDate: time.Now().Format("2006-01-02"),
    MaxDate: time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
    ShowTodayButton: true,
    OnChange: func(date string) {
        fmt.Printf("Appointment scheduled for: %s\n", date)
    },
})
```

### Disabled Date Input with Default Value

```go
// Create a disabled date input with a default value
holidayDate := widget.NewDateInput(ctx, widget.DateInputOptions{
    Label: "Holiday",
    DefaultValue: "2025-12-25", // Christmas
    Disabled: true,
    Format: "MMMM D, YYYY", // Display as "December 25, 2025"
})
```

## Best Practices

1. Always provide a clear label and placeholder text to guide users
2. Use a date format that is familiar to your target audience
3. Set appropriate `min_date` and `max_date` constraints when applicable
4. Consider disabling past dates for future-only selections
5. Provide a calendar picker interface for easy date selection
6. Allow keyboard input for users who prefer typing dates
7. Validate date inputs to ensure they match the expected format
8. Consider localization requirements for international date formats

## Related Components

- [DateTimeInput](./date-time-input) - For selecting both date and time
- `DateRangeInput` - For selecting a range between two dates
- [TimeInput](./time-input) - For selecting only time without a date
