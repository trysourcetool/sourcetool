---
sidebar_position: 16
---

# Columns

The Columns widget provides a layout container that arranges child widgets horizontally in a row with customizable spacing.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `columns` | `int` | - | Number of columns in the layout |
| `weight` | `[]int` | `[]` | Weight of each column, allowing for proportional sizing |

## Examples

### Basic Two-Column Layout

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/columns"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic two-column layout
        cols := ui.Columns(2)
        nameInput := cols[0].TextInput("Name", textinput.Placeholder("Enter name"))
        emailInput := cols[1].TextInput("Email", textinput.Placeholder("Enter email"))
    }
}
```

### Weighted Column Layout

```go
baseCols := ui.Columns(2, columns.Weight(3, 1))
```

## Related Components

- [Form](./form) - Container for organizing form elements
