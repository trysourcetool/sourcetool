---
sidebar_position: 1
---

# Text Input

The Text Input widget provides a single-line text input field that allows users to enter and edit text.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `string` | `None` | Current text content of the input |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when the input is empty |
| `default_value` | `string` | `None` | Initial value of the input |
| `required` | `bool` | `False` | Whether the input is required |
| `disabled` | `bool` | `False` | Whether the input is disabled |
| `max_length` | `int32` | `None` | Maximum number of characters allowed |
| `min_length` | `int32` | `None` | Minimum number of characters allowed |

## Examples

### Basic Text Input

```go
package main

import (
    "github.com/sourcetool/sourcetool-go"
    "github.com/sourcetool/sourcetool-go/textinput"
)

func main() {
    // Create a basic text input
    textInput := ui.TextInput("Name", textinput.Placeholder("Enter your name"))
}
```

### Disabled Text Input

```go
// Create a disabled text input
textInput := ui.TextInput("Name", textinput.Placeholder("Enter your name"), textinput.Disabled(true))
```

### Text Input with Length Constraints

```go
// Create a text input with length constraints
passwordInput := ui.TextInput("Password", textinput.Placeholder("Enter a secure password"), textinput.Required(true), textinput.MinLength(8), textinput.MaxLength(64))
```

## Related Components

- [TextArea](./textarea) - For multi-line text input
- [NumberInput](./number-input) - For numerical input
