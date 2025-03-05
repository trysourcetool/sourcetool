---
sidebar_position: 14
---

# Radio

The Radio widget provides a group of mutually exclusive options where users can select only one option at a time.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `int32` | `None` | Current selected option index |
| `label` | `string` | `""` | Label text displayed above the radio group |
| `options` | `[]string` | `[]` | Array of option labels for each radio button |
| `default_value` | `int32` | `None` | Initial selected option index |
| `required` | `bool` | `False` | Whether a selection is required |
| `disabled` | `bool` | `False` | Whether the entire radio group is disabled |

## Event Handling

The Radio widget emits events when a selection is made:

```go
// Define a radio group with change handler
radioGroup := widget.NewRadio(ctx, widget.RadioOptions{
    Label: "Select Payment Method",
    Options: []string{
        "Credit Card",
        "PayPal",
        "Bank Transfer",
        "Cryptocurrency",
    },
    OnChange: func(value int) {
        // Handle radio selection change
        fmt.Printf("Selected payment method: %v\n", value)
    },
})
```

## Examples

### Basic Radio Group

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic radio group
    genderRadio := widget.NewRadio(ctx, widget.RadioOptions{
        Label: "Gender",
        Options: []string{
            "Male",
            "Female",
            "Non-binary",
            "Prefer not to say",
        },
        OnChange: func(value int) {
            fmt.Printf("Selected gender option: %v\n", value)
        },
    })
    
    // Add the radio group to your UI
    container.Add(genderRadio)
}
```

### Radio Group with Default Selection

```go
// Create a radio group with default selection
experienceRadio := widget.NewRadio(ctx, widget.RadioOptions{
    Label: "Experience Level",
    Options: []string{
        "Beginner",
        "Intermediate",
        "Advanced",
        "Expert",
    },
    DefaultValue: 1, // Intermediate selected by default
    OnChange: func(value int) {
        fmt.Printf("Selected experience level: %v\n", value)
    },
})
```

### Required Radio Group

```go
// Create a required radio group
shippingRadio := widget.NewRadio(ctx, widget.RadioOptions{
    Label: "Shipping Method",
    Options: []string{
        "Standard (3-5 business days)",
        "Express (1-2 business days)",
        "Same Day",
    },
    Required: true, // A selection is required
    OnChange: func(value int) {
        fmt.Printf("Selected shipping method: %v\n", value)
    },
})
```

### Disabled Radio Group

```go
// Create a disabled radio group
planRadio := widget.NewRadio(ctx, widget.RadioOptions{
    Label: "Subscription Plan (Currently Unavailable)",
    Options: []string{
        "Basic",
        "Premium",
        "Enterprise",
    },
    DefaultValue: 0, // Basic selected by default
    Disabled: true, // Entire group is disabled
})
```

## Best Practices

1. Use clear, concise labels for both the group and individual options
2. Arrange options in a logical order (e.g., alphabetical, numerical, or by frequency of use)
3. Limit the number of options to maintain usability (consider a Select widget for many options)
4. Provide a clear group label that describes the collection of options
5. Use default selections when there is a common or recommended choice
6. Consider the `required` property when a selection is necessary
7. Ensure sufficient spacing between options for touch interfaces
8. Use consistent styling for all radio groups in your application
9. Provide immediate visual feedback when an option is selected
10. Consider using vertical alignment for more than 3-4 options

## Related Components

- [Select](./select) - Alternative for selecting a single option from a dropdown
- [CheckboxGroup](./checkbox-group) - For selecting multiple options from a group
- `SegmentedControl` - For a more compact single-selection interface
