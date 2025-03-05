---
sidebar_position: 1
---

# Text Input

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `string` | `None` | Current text content of the input |
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when the input is empty |
| `default_value` | `string` | `None` | Initial value of the input |
| `required` | `bool` | `False` | Whether the input is required |
| `disabled` | `bool` | `False` | Whether the input is disabled |
| `max_length` | `int32` | `None` | Maximum number of characters allowed |
| `min_length` | `int32` | `None` | Minimum number of characters allowed |

## Event Handling

The Text Input widget emits events when its value changes:

```go
// Define a text input with change handler
textInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Username",
    Placeholder: "Enter your username",
    OnChange: func(value string) {
        // Handle text input value change
        fmt.Printf("Text input value changed: %s\n", value)
    },
})
```

## Examples

### Basic Text Input

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic text input
    nameInput := widget.NewTextInput(ctx, widget.TextInputOptions{
        Label: "Full Name",
        Placeholder: "Enter your full name",
        OnChange: func(value string) {
            fmt.Printf("Name: %s\n", value)
        },
    })
    
    // Add the text input to your UI
    container.Add(nameInput)
}
```

### Required Text Input with Validation

```go
// Create a required text input with validation
emailInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Email Address",
    Placeholder: "example@domain.com",
    Required: true,
    OnChange: func(value string) {
        // Validate email format
        if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
            fmt.Println("Please enter a valid email address")
        } else {
            fmt.Printf("Valid email: %s\n", value)
        }
    },
})
```

### Disabled Text Input

```go
// Create a disabled text input
readOnlyInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "User ID",
    DefaultValue: "USR12345",
    Disabled: true,
})
```

### Text Input with Length Constraints

```go
// Create a text input with length constraints
passwordInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Password",
    Placeholder: "Enter a secure password",
    Required: true,
    MinLength: 8,
    MaxLength: 64,
    Type: "password", // Masks the input
    OnChange: func(value string) {
        if len(value) < 8 {
            fmt.Println("Password must be at least 8 characters")
        } else {
            fmt.Println("Password length is valid")
        }
    },
})
```

## Best Practices

1. Always provide a meaningful placeholder text to guide users
2. Use the `required` property when the field must not be empty
3. Set appropriate `min_length` and `max_length` for data validation
4. Consider using `disabled` state for read-only information

## Related Components

- [TextArea](./textarea) - For multi-line text input
- [NumberInput](./number-input) - For numerical input
- `PasswordInput` - For secure password entry
