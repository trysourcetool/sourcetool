---
sidebar_position: 7
---

# Multi Select

The Multi Select widget provides a dropdown menu that allows users to select multiple options from a list of choices.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `[]int32` | `[]` | Array of selected option indices |
| `label` | `string` | `""` | Label text displayed above the multi-select |
| `options` | `[]string` | `[]` | Array of option labels to display in the dropdown |
| `placeholder` | `string` | `""` | Placeholder text displayed when no options are selected |
| `default_value` | `[]int32` | `[]` | Initial selected option indices |
| `required` | `bool` | `False` | Whether at least one selection is required |
| `disabled` | `bool` | `False` | Whether the multi-select is disabled |

## Event Handling

The Multi Select widget emits events when selections change:

```go
// Define a multi-select with change handler
multiSelect := widget.NewMultiSelect(ctx, widget.MultiSelectOptions{
    Label: "Skills",
    Placeholder: "Select your skills",
    Options: []widget.Option{
        {Value: 1, Label: "Go"},
        {Value: 2, Label: "JavaScript"},
        {Value: 3, Label: "Python"},
        {Value: 4, Label: "SQL"},
    },
    OnChange: func(values []interface{}) {
        // Handle selection changes
        fmt.Printf("Selected skills: %v\n", values)
    },
})
```

## Examples

### Basic Multi Select

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic multi-select
    skillsSelect := widget.NewMultiSelect(ctx, widget.MultiSelectOptions{
        Label: "Skills",
        Placeholder: "Select your skills",
        Options: []widget.Option{
            {Value: 1, Label: "Go"},
            {Value: 2, Label: "JavaScript"},
            {Value: 3, Label: "Python"},
            {Value: 4, Label: "SQL"},
            {Value: 5, Label: "Docker"},
        },
        OnChange: func(values []interface{}) {
            fmt.Printf("Selected skills: %v\n", values)
        },
    })
    
    // Add the multi-select to your UI
    container.Add(skillsSelect)
}
```

### Multi Select with Default Values

```go
// Create a multi-select with default values
interestsSelect := widget.NewMultiSelect(ctx, widget.MultiSelectOptions{
    Label: "Interests",
    Options: []widget.Option{
        {Value: "tech", Label: "Technology"},
        {Value: "sci", Label: "Science"},
        {Value: "art", Label: "Art"},
        {Value: "music", Label: "Music"},
        {Value: "sports", Label: "Sports"},
    },
    DefaultValue: []interface{}{"tech", "sci"},
    OnChange: func(values []interface{}) {
        fmt.Printf("Selected interests: %v\n", values)
    },
})
```

### Searchable Multi Select with Maximum Selection

```go
// Create a searchable multi-select with maximum selection limit
tagsSelect := widget.NewMultiSelect(ctx, widget.MultiSelectOptions{
    Label: "Tags",
    Placeholder: "Search and select tags (max 3)",
    Options: []widget.Option{
        {Value: "go", Label: "golang"},
        {Value: "web", Label: "web-development"},
        {Value: "api", Label: "api"},
        {Value: "db", Label: "database"},
        {Value: "ui", Label: "user-interface"},
        {Value: "cloud", Label: "cloud-computing"},
        // ... many more tags
    },
    Searchable: true,
    MaxSelected: 3,
    OnChange: func(values []interface{}) {
        fmt.Printf("Selected tags: %v\n", values)
        if len(values) >= 3 {
            fmt.Println("Maximum tags selected")
        }
    },
})
```

## Best Practices

1. Always provide a clear label and placeholder text to guide users
2. Use concise, descriptive option labels
3. Order options in a logical way (alphabetical, numerical, or by frequency of use)
4. Enable the `searchable` property when there are many options
5. Consider setting a `max_selected` limit when appropriate
6. Show selected items in a visually distinct way
7. Provide a way to easily remove individual selections
8. Consider grouping related options when the list is long
9. Ensure the dropdown is wide enough to accommodate the longest option label

## Related Components

- [Select](./select) - For selecting a single option from a list
- [CheckboxGroup](./checkbox-group) - Alternative for selecting multiple options from a visible list
- `TagInput` - For free-form entry of multiple tags or values
