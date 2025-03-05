---
sidebar_position: 2
---

# Button

The Button widget provides a clickable button element that can trigger actions when clicked.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `bool` | `False` | Current state of the button |
| `label` | `string` | `""` | Text displayed on the button |
| `disabled` | `bool` | `False` | Whether the button is disabled |

## Event Handling

The Button widget emits events when clicked:

```go
// Define a button with a click handler
button := widget.NewButton(ctx, widget.ButtonOptions{
    Label: "Submit",
    OnClick: func() {
        // Handle button click
        fmt.Println("Button clicked!")
        // Perform actions like form submission, navigation, etc.
    },
})
```

## Examples

### Basic Button

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic button
    button := widget.NewButton(ctx, widget.ButtonOptions{
        Label: "Click Me",
        OnClick: func() {
            fmt.Println("Button clicked!")
        },
    })
    
    // Add the button to your UI
    container.Add(button)
}
```

### Disabled Button

```go
// Create a disabled button
disabledButton := widget.NewButton(ctx, widget.ButtonOptions{
    Label: "Cannot Click",
    Disabled: true,
    OnClick: func() {
        // This won't be called while the button is disabled
        fmt.Println("Button clicked!")
    },
})
```

### Button with Icon

```go
// Create a button with an icon
iconButton := widget.NewButton(ctx, widget.ButtonOptions{
    Label: "Save",
    Icon: "save",
    IconPosition: "left",
    Variant: "outline",
    OnClick: func() {
        fmt.Println("Saving data...")
        // Save data logic
    },
})
```

## Best Practices

1. Use clear, action-oriented labels that describe what the button does
2. Provide visual feedback for button interactions (hover, active states)
3. Disable buttons when the action they perform is not available
4. Use consistent button styling throughout your application
5. Consider using icons to enhance the meaning of button actions
6. Ensure buttons have sufficient contrast and are easily clickable

## Related Components

- [Form](./form) - Container for form elements including buttons
- `IconButton` - Button that displays only an icon without text
- `LinkButton` - Button styled as a link
