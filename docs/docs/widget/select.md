---
sidebar_position: 6
---

# Select

The Select widget provides a dropdown menu that allows users to choose one option from a list of choices.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `int32` | `None` | Current selected option index |
| `label` | `string` | `""` | Label text displayed above the select |
| `options` | `[]string` | `[]` | Array of option labels to display in the dropdown |
| `placeholder` | `string` | `""` | Placeholder text displayed when no option is selected |
| `default_value` | `int32` | `None` | Initial selected option index |
| `required` | `bool` | `False` | Whether a selection is required |
| `disabled` | `bool` | `False` | Whether the select is disabled |

## Event Handling

The Select widget emits events when a selection is made:

```go
// Define a select with change handler
select := widget.NewSelect(ctx, widget.SelectOptions{
    Label: "Country",
    Placeholder: "Select a country",
    Options: []widget.Option{
        {Value: "us", Label: "United States"},
        {Value: "ca", Label: "Canada"},
        {Value: "mx", Label: "Mexico"},
    },
    OnChange: func(value interface{}) {
        // Handle selection change
        fmt.Printf("Selected country: %v\n", value)
    },
})
```

## Examples

### Basic Select

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic select
    countrySelect := widget.NewSelect(ctx, widget.SelectOptions{
        Label: "Country",
        Placeholder: "Choose a country",
        Options: []widget.Option{
            {Value: "us", Label: "United States"},
            {Value: "ca", Label: "Canada"},
            {Value: "uk", Label: "United Kingdom"},
            {Value: "au", Label: "Australia"},
        },
        OnChange: func(value interface{}) {
            fmt.Printf("Selected country: %v\n", value)
        },
    })
    
    // Add the select to your UI
    container.Add(countrySelect)
}
```

### Select with Default Value

```go
// Create a select with a default value
languageSelect := widget.NewSelect(ctx, widget.SelectOptions{
    Label: "Programming Language",
    Options: []widget.Option{
        {Value: "go", Label: "Go"},
        {Value: "js", Label: "JavaScript"},
        {Value: "py", Label: "Python"},
        {Value: "rb", Label: "Ruby"},
    },
    DefaultValue: "go",
    OnChange: func(value interface{}) {
        fmt.Printf("Selected language: %v\n", value)
    },
})
```

### Searchable Select

```go
// Create a searchable select with many options
citySelect := widget.NewSelect(ctx, widget.SelectOptions{
    Label: "City",
    Placeholder: "Search for a city",
    Options: []widget.Option{
        {Value: "nyc", Label: "New York City"},
        {Value: "la", Label: "Los Angeles"},
        {Value: "chi", Label: "Chicago"},
        {Value: "hou", Label: "Houston"},
        // ... many more cities
    },
    Searchable: true,
    Clearable: true,
    OnChange: func(value interface{}) {
        fmt.Printf("Selected city: %v\n", value)
    },
})
```

## Best Practices

1. Always provide a clear label and placeholder text to guide users
2. Use concise, descriptive option labels
3. Order options in a logical way (alphabetical, numerical, or by frequency of use)
4. Enable the `searchable` property when there are many options
5. Consider using the `clearable` property when selection is optional
6. Group related options when the list is long
7. Provide a meaningful default value when appropriate
8. Ensure the dropdown is wide enough to accommodate the longest option label

## Related Components

- [MultiSelect](./multi-select) - For selecting multiple options from a list
- `ComboBox` - Combination of text input and select for custom entries
- [Radio](./radio) - Alternative for selecting from a small set of visible options
