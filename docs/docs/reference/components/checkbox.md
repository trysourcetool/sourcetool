---
sidebar_position: 5
---

# Checkbox

The Checkbox widget provides a toggleable input control that allows users to select or deselect an option.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `bool` | `False` | Whether the checkbox is checked |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Text displayed next to the checkbox |
| `default_value` | `bool` | `False` | Initial checked state of the checkbox |
| `required` | `bool` | `False` | Whether the checkbox must be checked |
| `disabled` | `bool` | `False` | Whether the checkbox is disabled |

## Examples

### Basic Checkbox

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/checkbox"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic checkbox
        checkbox := ui.Checkbox("Subscribe to newsletter")
    }
}
```

### Default Checkbox

```go
// Create a default checkbox
checkbox := ui.Checkbox("Subscribe to newsletter", checkbox.DefaultValue(true))
```

### Required Checkbox

```go
// Create a required checkbox with description
checkbox := ui.Checkbox("Subscribe to newsletter", checkbox.Required(true))
```

### Disabled Checkbox

```go
// Create a disabled checkbox
checkbox := ui.Checkbox("Subscribe to newsletter", checkbox.Disabled(true))
```

## Related Components

- [CheckboxGroup](./checkbox-group) - For managing multiple related checkboxes
- [Radio](./radio) - For mutually exclusive options (unlike checkboxes)
