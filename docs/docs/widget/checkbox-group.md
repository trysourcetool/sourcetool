---
sidebar_position: 13
---

# Checkbox Group

The Checkbox Group widget provides a collection of related checkboxes that allow users to select multiple options from a predefined set.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `[]int32` | `[]` | Array of selected option indices |
| `label` | `string` | `""` | Label text displayed above the checkbox group |
| `options` | `[]string` | `[]` | Array of option labels for each checkbox |
| `default_value` | `[]int32` | `[]` | Initial selected option indices |
| `required` | `bool` | `False` | Whether at least one option must be selected |
| `disabled` | `bool` | `False` | Whether the entire checkbox group is disabled |

## Event Handling

The Checkbox Group widget emits events when selections change:

```go
// Define a checkbox group with change handler
checkboxGroup := widget.NewCheckboxGroup(ctx, widget.CheckboxGroupOptions{
    Label: "Select Interests",
    Options: []string{
        "Technology",
        "Science",
        "Art",
        "Sports",
        "Music",
    },
    OnChange: func(values []int) {
        // Handle checkbox group selection changes
        fmt.Printf("Selected interests: %v\n", values)
    },
})
```

## Examples

### Basic Checkbox Group

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic checkbox group
    interestsGroup := widget.NewCheckboxGroup(ctx, widget.CheckboxGroupOptions{
        Label: "Select your interests",
        Options: []string{
            "Technology",
            "Science",
            "Art",
            "Sports",
            "Music",
        },
        OnChange: func(values []int) {
            fmt.Printf("Selected interests: %v\n", values)
        },
    })
    
    // Add the checkbox group to your UI
    container.Add(interestsGroup)
}
```

### Checkbox Group with Default Selections

```go
// Create a checkbox group with default selections
notificationsGroup := widget.NewCheckboxGroup(ctx, widget.CheckboxGroupOptions{
    Label: "Notification Preferences",
    Options: []string{
        "Email notifications",
        "Push notifications",
        "SMS notifications",
        "In-app notifications",
    },
    DefaultValue: []int{0, 3}, // Email and in-app notifications selected by default
    OnChange: func(values []int) {
        fmt.Printf("Selected notification types: %v\n", values)
    },
})
```

### Required Checkbox Group

```go
// Create a required checkbox group
termsGroup := widget.NewCheckboxGroup(ctx, widget.CheckboxGroupOptions{
    Label: "Terms and Conditions",
    Options: []string{
        "I agree to the Terms of Service",
        "I agree to the Privacy Policy",
        "I agree to receive marketing emails (optional)",
    },
    Required: true, // At least one option must be selected
    OnChange: func(values []int) {
        if len(values) == 0 {
            fmt.Println("You must agree to at least one term")
        } else {
            fmt.Printf("Agreed to terms: %v\n", values)
        }
    },
})
```

### Disabled Checkbox Group

```go
// Create a disabled checkbox group
planFeaturesGroup := widget.NewCheckboxGroup(ctx, widget.CheckboxGroupOptions{
    Label: "Premium Plan Features (Unavailable)",
    Options: []string{
        "Advanced analytics",
        "Priority support",
        "Custom branding",
        "API access",
    },
    DefaultValue: []int{0, 1, 2, 3}, // All selected
    Disabled: true, // Entire group is disabled
})
```

## Best Practices

1. Use clear, concise labels for both the group and individual options
2. Group related options together in a logical order
3. Consider using a vertical layout for more than 3-4 options
4. Provide a clear group label that describes the collection of options
5. Use default selections when there are common or recommended choices
6. Consider the `required` property when at least one selection is necessary
7. Provide visual feedback when options are selected
8. Ensure sufficient spacing between options for touch interfaces
9. Use consistent styling for all checkbox groups in your application
10. Consider using a Select or MultiSelect widget when there are many options

## Related Components

- [Checkbox](./checkbox) - For a single checkbox option
- [Radio](./radio) - For selecting a single option from a group (mutually exclusive)
- [MultiSelect](./multi-select) - Alternative for selecting multiple options from a dropdown
