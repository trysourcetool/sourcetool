---
sidebar_position: 3
---

# Button

The Button widget provides a clickable button element that can trigger actions when clicked.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `bool` | `False` | Current state of the button |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Text displayed on the button |
| `disabled` | `bool` | `False` | Whether the button is disabled |

## Examples

### Basic Button

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/button"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic button
        button := ui.Button("Click Me")
    }
}
```

### Disabled Button

```go
// Create a disabled button
button := ui.Button("Cannot Click", button.Disabled(true))
```