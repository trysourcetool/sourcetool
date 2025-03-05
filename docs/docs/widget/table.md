---
sidebar_position: 11
---

# Table

The Table widget provides a way to display and interact with tabular data, supporting features like sorting, pagination, and row selection.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `data` | `bytes` | `""` | Table data as base64 encoded JSON string |
| `value` | `TableValue` | `None` | Current selection state of the table |
| `header` | `string` | `""` | Optional header text displayed above the table |
| `description` | `string` | `""` | Optional description text displayed below the header |
| `height` | `int32` | `None` | Number of rows to display per page |
| `column_order` | `[]string` | `[]` | Custom order for columns |
| `on_select` | `string` | `""` | Event handler for selection |
| `row_selection` | `string` | `""` | Row selection mode: "single", "multiple", or empty string |

## Event Handling

The Table widget emits events for various interactions:

```go
// Define a table with selection handler
table := widget.NewTable(ctx, widget.TableOptions{
    Header: "Employee Directory",
    Data: employeeData,
    Columns: []widget.Column{
        {Key: "id", Label: "ID", Sortable: true},
        {Key: "name", Label: "Name", Sortable: true},
        {Key: "department", Label: "Department"},
        {Key: "email", Label: "Email"},
    },
    RowSelection: "single",
    OnRowSelect: func(selection map[string]interface{}) {
        // Handle row selection
        fmt.Printf("Selected employee: %v\n", selection)
    },
})
```

## Examples

### Basic Table

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic table
    data := [][]interface{}{
        {"1", "John Doe", "Engineering", "john@example.com"},
        {"2", "Jane Smith", "Marketing", "jane@example.com"},
        {"3", "Bob Johnson", "Finance", "bob@example.com"},
    }
    
    table := widget.NewTable(ctx, widget.TableOptions{
        Header: "Employee Directory",
        Description: "List of company employees",
        Data: data,
        Columns: []widget.Column{
            {Key: "id", Label: "ID"},
            {Key: "name", Label: "Name"},
            {Key: "department", Label: "Department"},
            {Key: "email", Label: "Email"},
        },
        Striped: true,
        HoverEffect: true,
    })
    
    // Add the table to your UI
    container.Add(table)
}
```

### Sortable Table with Pagination

```go
// Create a sortable table with pagination
productTable := widget.NewTable(ctx, widget.TableOptions{
    Header: "Product Inventory",
    Data: productData,
    Columns: []widget.Column{
        {Key: "id", Label: "ID", Sortable: true},
        {Key: "name", Label: "Product Name", Sortable: true},
        {Key: "category", Label: "Category", Sortable: true},
        {Key: "price", Label: "Price", Sortable: true, Format: "$%.2f"},
        {Key: "stock", Label: "In Stock", Sortable: true},
    },
    Height: 10, // 10 rows per page
    DefaultSortColumn: "name",
    DefaultSortDirection: "asc",
    Bordered: true,
    OnSort: func(column string, direction string) {
        fmt.Printf("Table sorted by %s in %s order\n", column, direction)
    },
})
```

### Table with Row Selection

```go
// Create a table with multiple row selection
userTable := widget.NewTable(ctx, widget.TableOptions{
    Header: "User Management",
    Data: userData,
    Columns: []widget.Column{
        {Key: "id", Label: "ID"},
        {Key: "username", Label: "Username"},
        {Key: "role", Label: "Role"},
        {Key: "lastLogin", Label: "Last Login", Format: "2006-01-02 15:04:05"},
        {Key: "active", Label: "Active", Template: "{{if .active}}Yes{{else}}No{{end}}"},
    },
    RowSelection: "multiple",
    OnRowSelect: func(selection map[string]interface{}) {
        selectedRows := selection["rows"].([]int)
        fmt.Printf("Selected %d users\n", len(selectedRows))
    },
})
```

### Custom Column Rendering

```go
// Create a table with custom column rendering
reportTable := widget.NewTable(ctx, widget.TableOptions{
    Header: "Monthly Reports",
    Data: reportData,
    Columns: []widget.Column{
        {Key: "month", Label: "Month"},
        {Key: "revenue", Label: "Revenue", Format: "$%.2f"},
        {Key: "expenses", Label: "Expenses", Format: "$%.2f"},
        {Key: "profit", Label: "Profit", 
            Template: "{{if gt .profit 0}}<span class=\"text-green\">+${{.profit}}</span>{{else}}<span class=\"text-red\">-${{abs .profit}}</span>{{end}}"},
        {Key: "growth", Label: "Growth", 
            Template: "{{if gt .growth 0}}↑ {{.growth}}%{{else}}↓ {{abs .growth}}%{{end}}"},
    },
    Striped: true,
    Bordered: true,
})
```

## Best Practices

1. Provide clear column headers that describe the data
2. Use appropriate formatting for different data types (dates, currency, percentages)
3. Enable sorting for columns where ordering is meaningful
4. Use pagination for large datasets to improve performance
5. Consider using striped rows to improve readability
6. Implement row selection when users need to perform actions on specific rows
7. Keep tables responsive by allowing horizontal scrolling on small screens
8. Use consistent styling across all tables in your application
9. Consider adding a search or filter capability for large datasets
10. Provide appropriate empty state messaging when no data is available

## Related Components

- `DataGrid` - Advanced table with filtering, grouping, and editing capabilities
- [Form](./form) - Container for organizing form elements
- [TextInput](./text-input) - For filtering table data
