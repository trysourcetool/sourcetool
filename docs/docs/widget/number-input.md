---
sidebar_position: 4
---

# Number Input

The Number Input widget provides a specialized input field for numerical values with optional validation and formatting.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `double` | `None` | Current numeric value of the input |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when the input is empty |
| `default_value` | `double` | `None` | Initial value of the input |
| `required` | `bool` | `False` | Whether the input is required |
| `disabled` | `bool` | `False` | Whether the input is disabled |
| `max_value` | `double` | `None` | Maximum allowed value |
| `min_value` | `double` | `None` | Minimum allowed value |

## Examples

### Basic Number Input

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/numberinput"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic number input
        numberInput := ui.NumberInput("Age")
    }
}
```

### Disabled Number Input

```go
// Create a disabled number input
disabledInput := ui.NumberInput("Score (Read Only)", numberinput.Disabled(true))
```

## Related Components

- [TextInput](./text-input) - For general text input
