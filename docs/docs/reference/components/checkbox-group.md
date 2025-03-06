---
sidebar_position: 13
---

# Checkbox Group

The Checkbox Group widget provides a collection of related checkboxes that allow users to select multiple options from a predefined set.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `[]int32` | `[]` | Array of selected option indices |


## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the checkbox group |
| `options` | `[]string` | `[]` | Array of option labels for each checkbox |
| `default_value` | `[]int32` | `[]` | Initial selected option indices |
| `required` | `bool` | `False` | Whether at least one option must be selected |
| `disabled` | `bool` | `False` | Whether the entire checkbox group is disabled |
| `format_func` | `func(string, int) string` | `nil` | Function to format the option label |

## Examples

### Basic Checkbox Group

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/checkboxgroup"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic checkbox group
        checkboxGroup := ui.CheckboxGroup("Select your interests", checkboxgroup.Options("Technology", "Science", "Art", "Sports", "Music"))
    }
}
```

### Checkbox Group with Default Selections

```go
// Create a checkbox group with default selections
checkboxGroup := ui.CheckboxGroup("Select your interests", checkboxgroup.Options("Technology", "Science", "Art", "Sports", "Music"), checkboxgroup.DefaultValue("Technology", "Sports"))
```

### Required Checkbox Group

```go
// Create a required checkbox group
checkboxGroup := ui.CheckboxGroup("Select your interests", checkboxgroup.Options("Technology", "Science", "Art", "Sports", "Music"), checkboxgroup.Required(true))
```

### Disabled Checkbox Group

```go
// Create a disabled checkbox group
checkboxGroup := ui.CheckboxGroup("Select your interests", checkboxgroup.Options("Technology", "Science", "Art", "Sports", "Music"), checkboxgroup.DefaultValue("Technology", "Sports"), checkboxgroup.Disabled(true))
```

## Related Components

- [Checkbox](./checkbox) - For a single checkbox option
- [Radio](./radio) - For selecting a single option from a group (mutually exclusive)
- [MultiSelect](./multi-select) - Alternative for selecting multiple options from a dropdown
