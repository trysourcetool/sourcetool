---
sidebar_position: 7
---

# Multi Select

The Multi Select widget provides a dropdown menu that allows users to select multiple options from a list of choices.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `[]int32` | `[]` | Array of selected option indices |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the multi-select |
| `options` | `[]string` | `[]` | Array of option labels to display in the dropdown |
| `placeholder` | `string` | `""` | Placeholder text displayed when no options are selected |
| `default_value` | `[]int32` | `[]` | Initial selected option indices |
| `required` | `bool` | `False` | Whether at least one selection is required |
| `disabled` | `bool` | `False` | Whether the multi-select is disabled |

## Examples

### Basic Multi Select

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/multiselect"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic multi-select
        multiSelect := ui.MultiSelect("Skills", multiselect.Placeholder("Select your skills"), multiselect.Options("Go", "JavaScript", "Python", "SQL", "Docker"))
    }
}
```

### Multi Select with Default Values

```go
// Create a multi-select with default values
multiSelect := ui.MultiSelect("Skills", multiselect.Placeholder("Select your skills"), multiselect.Options("Go", "JavaScript", "Python", "SQL", "Docker"), multiselect.DefaultValue("Go", "Python"))
```

## Related Components

- [Select](./select) - For selecting a single option from a list
- [CheckboxGroup](./checkbox-group) - Alternative for selecting multiple options from a visible list
