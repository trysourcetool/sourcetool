---
sidebar_position: 5
---

# Checkbox

The Checkbox widget provides a toggleable input control that allows users to select or deselect an option.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `bool` | `False` | Whether the checkbox is checked |
| `label` | `string` | `""` | Text displayed next to the checkbox |
| `default_value` | `bool` | `False` | Initial checked state of the checkbox |
| `required` | `bool` | `False` | Whether the checkbox must be checked |
| `disabled` | `bool` | `False` | Whether the checkbox is disabled |

## Event Handling

The Checkbox widget emits events when its state changes:

```go
// Define a checkbox with change handler
checkbox := widget.NewCheckbox(ctx, widget.CheckboxOptions{
    Label: "I agree to the terms and conditions",
    Checked: false,
    OnChange: func(checked bool) {
        // Handle checkbox state change
        if checked {
            fmt.Println("User agreed to terms")
        } else {
            fmt.Println("User disagreed with terms")
        }
    },
})
```

## Examples

### Basic Checkbox

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic checkbox
    checkbox := widget.NewCheckbox(ctx, widget.CheckboxOptions{
        Label: "Subscribe to newsletter",
        Checked: false,
        OnChange: func(checked bool) {
            fmt.Printf("Newsletter subscription: %v\n", checked)
        },
    })
    
    // Add the checkbox to your UI
    container.Add(checkbox)
}
```

### Required Checkbox with Description

```go
// Create a required checkbox with description
termsCheckbox := widget.NewCheckbox(ctx, widget.CheckboxOptions{
    Label: "I agree to the terms and conditions",
    Description: "You must accept the terms to continue",
    Required: true,
    OnChange: func(checked bool) {
        if checked {
            fmt.Println("Terms accepted")
        } else {
            fmt.Println("Terms must be accepted to proceed")
        }
    },
})
```

### Disabled Checkbox

```go
// Create a disabled checkbox
disabledCheckbox := widget.NewCheckbox(ctx, widget.CheckboxOptions{
    Label: "Premium features (unavailable)",
    Checked: false,
    Disabled: true,
})
```

### Indeterminate Checkbox

```go
// Create an indeterminate checkbox (useful for parent checkboxes in a tree)
parentCheckbox := widget.NewCheckbox(ctx, widget.CheckboxOptions{
    Label: "Select all items",
    Indeterminate: true, // Some but not all child items are selected
    OnChange: func(checked bool) {
        fmt.Println("User selected all items:", checked)
        // Update child checkboxes accordingly
    },
})
```

## Best Practices

1. Use clear, concise labels that describe the option being toggled
2. Place checkboxes in logical groups when presenting multiple related options
3. Use the `required` property for checkboxes that must be checked (e.g., terms acceptance)
4. Provide descriptive text to explain the implications of checking/unchecking
5. Use the indeterminate state appropriately for parent-child checkbox relationships
6. Ensure checkboxes have sufficient spacing for touch interfaces
7. Consider using checkbox groups for related options

## Related Components

- [CheckboxGroup](./checkbox-group) - For managing multiple related checkboxes
- `Switch` - Alternative toggle control with a different visual appearance
- [Radio](./radio) - For mutually exclusive options (unlike checkboxes)
