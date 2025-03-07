---
sidebar_position: 6
---

# Select

The Select widget provides a dropdown menu that allows users to choose one option from a list of choices.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `int32` | `None` | Current selected option index |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the select |
| `options` | `[]string` | `[]` | Array of option labels to display in the dropdown |
| `placeholder` | `string` | `""` | Placeholder text displayed when no option is selected |
| `default_value` | `int32` | `None` | Initial selected option index |
| `required` | `bool` | `False` | Whether a selection is required |
| `disabled` | `bool` | `False` | Whether the select is disabled |

## Examples

### Basic Select

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/selectbox"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic select
        selectbox := ui.Selectbox("Country", selectbox.Placeholder("Choose a country"), selectbox.Options("United States", "Canada", "United Kingdom", "Australia"))
    }
}
```

### Select with Default Value

```go
// Create a select with a default value
selectbox := ui.Selectbox("Country", selectbox.Options("United States", "Canada", "United Kingdom", "Australia"), selectbox.DefaultValue("United States"))
```

## Related Components

- [MultiSelect](./multi-select) - For selecting multiple options from a list
- [Radio](./radio) - Alternative for selecting from a small set of visible options
