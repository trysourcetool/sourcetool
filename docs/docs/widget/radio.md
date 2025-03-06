---
sidebar_position: 14
---

# Radio

The Radio widget provides a group of mutually exclusive options where users can select only one option at a time.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `int32` | `None` | Current selected option index |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the radio group |
| `options` | `[]string` | `[]` | Array of option labels for each radio button |
| `default_value` | `int32` | `None` | Initial selected option index |
| `required` | `bool` | `False` | Whether a selection is required |
| `disabled` | `bool` | `False` | Whether the entire radio group is disabled |

## Examples

### Basic Radio Group

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/radio"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic radio group
        radio := ui.Radio("Gender", radio.Options("male", "female", "Non-binary", "Prefer not to say"))
    }
}
```

### Radio Group with Default Selection

```go
// Create a radio group with default selection
radio := ui.Radio("Gender", radio.Options("male", "female", "Non-binary", "Prefer not to say"), radio.DefaultValue("male"))
```

### Required Radio Group

```go
// Create a required radio group
radio := ui.Radio("Gender", radio.Options("male", "female", "Non-binary", "Prefer not to say"), radio.Required(true))
```

### Disabled Radio Group

```go
// Create a disabled radio group
radio := ui.Radio("Gender", radio.Options("male", "female", "Non-binary", "Prefer not to say"), radio.DefaultValue("male"), radio.Disabled(true))
```

## Related Components

- [Select](./select) - Alternative for selecting a single option from a dropdown
- [CheckboxGroup](./checkbox-group) - For selecting multiple options from a group
