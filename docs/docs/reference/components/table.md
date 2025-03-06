---
sidebar_position: 11
---

# Table

The Table widget provides a way to display and interact with tabular data, supporting features like sorting, pagination, and row selection.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `TableValue` | `None` | Current selection state of the table |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `data` | `bytes` | `""` | Table data as base64 encoded JSON string |
| `header` | `string` | `""` | Optional header text displayed above the table |
| `description` | `string` | `""` | Optional description text displayed below the header |
| `height` | `int32` | `None` | Number of rows to display per page |
| `column_order` | `[]string` | `[]` | Custom order for columns |
| `on_select` | `string` | `""` | Event handler for selection |
| `row_selection` | `string` | `""` | Row selection mode: "single", "multiple", or empty string |

## Examples

### Basic Table

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/table"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Age       int
	Gender    string
	CreatedAt time.Time
}

baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
users := []User{
    {ID: "1", Name: "John Doe 001", Email: "john.doe+001@acme.com", Age: 25, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 0)},
    {ID: "2", Name: "John Doe 002", Email: "john.doe+002@acme.com", Age: 30, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 1)},
    {ID: "3", Name: "Jane Doe 003", Email: "jane.doe+003@acme.com", Age: 35, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 2)}
}

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic table
        table := ui.Table(
            users,
            table.Header("Users"),
            table.Height(10),
            table.ColumnOrder("ID", "Name", "Email", "Age", "Gender", "CreatedAt"),
		    table.OnSelect(table.SelectionBehaviorRerun),
	    )
	}
}
```

### Table with Row Selection

```go
// Create a table with multiple row selection
table := baseCols[0].Table(
    users,
    table.Header("Users"),
    table.Height(10),
    table.ColumnOrder("ID", "Name", "Email", "Age", "Gender", "CreatedAt"),
    table.RowSelection(table.SelectionModeMultiple),
)
```

## Related Components

- [Form](./form) - Container for organizing form elements
- [TextInput](./text-input) - For filtering table data
