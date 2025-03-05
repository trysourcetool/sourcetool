---
sidebar_position: 3
---

# TextArea

The TextArea widget provides a multi-line text input field for entering longer text content.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `string` | `None` | Current text content of the textarea |
| `label` | `string` | `""` | Label text displayed above the textarea |
| `placeholder` | `string` | `""` | Placeholder text displayed when the textarea is empty |
| `default_value` | `string` | `None` | Initial value of the textarea |
| `required` | `bool` | `False` | Whether the textarea input is required |
| `disabled` | `bool` | `False` | Whether the textarea is disabled |
| `max_length` | `int32` | `None` | Maximum number of characters allowed |
| `min_length` | `int32` | `None` | Minimum number of characters allowed |
| `max_lines` | `int32` | `None` | Maximum number of lines allowed |
| `min_lines` | `int32` | `None` | Minimum number of lines allowed |
| `auto_resize` | `bool` | `False` | Whether the textarea should automatically resize based on content |

## Event Handling

The TextArea widget emits events when its value changes:

```go
// Define a textarea with change handler
textarea := widget.NewTextArea(ctx, widget.TextAreaOptions{
    Label: "Description",
    Placeholder: "Enter a detailed description",
    DefaultValue: "",
    OnChange: func(value string) {
        // Handle textarea value change
        fmt.Println("TextArea value changed:", value)
    },
})
```

## Examples

### Basic TextArea

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic textarea
    textarea := widget.NewTextArea(ctx, widget.TextAreaOptions{
        Label: "Comments",
        Placeholder: "Enter your comments here",
        Rows: 4,
    })
    
    // Add the textarea to your UI
    container.Add(textarea)
}
```

### Required TextArea with Validation

```go
// Create a required textarea with validation
textarea := widget.NewTextArea(ctx, widget.TextAreaOptions{
    Label: "Feedback",
    Placeholder: "Please provide your feedback",
    Required: true,
    MinLength: 10,
    MaxLength: 500,
    OnChange: func(value string) {
        // Validate the input
        if len(value) < 10 {
            fmt.Println("Feedback must be at least 10 characters")
        } else {
            fmt.Println("Valid feedback received")
        }
    },
})
```

### Disabled TextArea

```go
// Create a disabled textarea
disabledTextarea := widget.NewTextArea(ctx, widget.TextAreaOptions{
    Label: "System Log",
    DefaultValue: "System log content that cannot be edited",
    Disabled: true,
    Rows: 6,
    Resize: "none",
})
```

## Best Practices

1. Always provide a meaningful label and placeholder text to guide users
2. Use the `required` property when the field must not be empty
3. Set appropriate `min_length` and `max_length` for data validation
4. Consider using `disabled` state for read-only information
5. Choose an appropriate number of `rows` based on the expected content length
6. Provide clear validation feedback for user input
7. Consider the resize behavior that makes sense for your UI layout

## Related Components

- [TextInput](./text-input) - For single-line text input
- [Markdown](./markdown) - For displaying formatted text content
- `CodeEditor` - For editing code with syntax highlighting
