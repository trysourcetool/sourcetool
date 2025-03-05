---
sidebar_position: 17
---

# Column Item

The Column Item widget is used within a Columns layout to define individual column items with specific weight or sizing properties.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `weight` | `double` | `1.0` | Relative width of the column compared to other columns |

## Event Handling

The Column Item widget does not emit events directly, as it is a layout container. Events are handled by the child widgets placed within the column item.

```go
// Define a column item with specific weight
columnItem := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 2.0, // This column will be twice as wide as a column with weight 1.0
})

// Add child widgets to the column item
columnItem.Add(childWidget)
```

## Examples

### Basic Column Items with Different Weights

```go
package main

import (
    "github.com/sourcetool/widget"
)

func main() {
    // Create a columns layout
    columns := widget.NewColumns(ctx, widget.ColumnsOptions{})
    
    // Create a narrow column (1/3 width)
    narrowColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
        Weight: 1.0,
    })
    
    // Create a wide column (2/3 width)
    wideColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
        Weight: 2.0,
    })
    
    // Add content to the narrow column
    sidebarTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
        Body: "## Navigation",
    })
    sidebarLinks := widget.NewMarkdown(ctx, widget.MarkdownOptions{
        Body: "- [Home](#)\n- [Products](#)\n- [About](#)\n- [Contact](#)",
    })
    narrowColumn.Add(sidebarTitle)
    narrowColumn.Add(sidebarLinks)
    
    // Add content to the wide column
    contentTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
        Body: "## Welcome to Our Website",
    })
    contentText := widget.NewMarkdown(ctx, widget.MarkdownOptions{
        Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam euismod, nisl eget aliquam ultricies, nunc nisl aliquet nunc, quis aliquam nisl nunc eu nisl. Nullam euismod, nisl eget aliquam ultricies, nunc nisl aliquet nunc, quis aliquam nisl nunc eu nisl.",
    })
    wideColumn.Add(contentTitle)
    wideColumn.Add(contentText)
    
    // Add columns to the layout
    columns.Add(narrowColumn)
    columns.Add(wideColumn)
    
    // Add the columns layout to your UI
    container.Add(columns)
}
```

### Three-Column Layout with Custom Weights

```go
// Create a three-column layout with custom weights
columns := widget.NewColumns(ctx, widget.ColumnsOptions{})

// Create a narrow sidebar column (1/6 width)
sidebarColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 1.0,
})

// Create a main content column (3/6 width)
mainColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 3.0,
})

// Create a supplementary column (2/6 width)
supplementaryColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 2.0,
})

// Add content to the sidebar column
sidebarColumn.Add(widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Navigation\n- [Dashboard](#)\n- [Reports](#)\n- [Settings](#)",
}))

// Add content to the main column
mainColumn.Add(widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Main Content\nThis is the primary content area of the application.",
}))
mainColumn.Add(widget.NewTextArea(ctx, widget.TextAreaOptions{
    Label: "Notes",
    Placeholder: "Enter your notes here",
}))

// Add content to the supplementary column
supplementaryColumn.Add(widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Additional Information\nThis section contains supplementary details and resources.",
}))
supplementaryColumn.Add(widget.NewTable(ctx, widget.TableOptions{
    Header: "Recent Activity",
    Data: activityData,
}))

// Add columns to the layout
columns.Add(sidebarColumn)
columns.Add(mainColumn)
columns.Add(supplementaryColumn)
```

### Responsive Dashboard Layout

```go
// Create a responsive dashboard layout
dashboard := widget.NewColumns(ctx, widget.ColumnsOptions{})

// Create a narrow sidebar column (1/5 width)
sidebarColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 1.0,
})

// Create a main content column (4/5 width)
contentColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 4.0,
})

// Add content to the sidebar column
profileSection := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## User Profile\n**John Doe**\nAdmin",
})
navigationSection := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Navigation\n- [Dashboard](#)\n- [Analytics](#)\n- [Reports](#)\n- [Settings](#)\n- [Help](#)",
})
sidebarColumn.Add(profileSection)
sidebarColumn.Add(navigationSection)

// Create a nested columns layout for the content area
contentColumns := widget.NewColumns(ctx, widget.ColumnsOptions{})

// Create top row widgets (full width)
statsRow := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 1.0,
})
statsTable := widget.NewTable(ctx, widget.TableOptions{
    Header: "Key Performance Indicators",
    Data: kpiData,
})
statsRow.Add(statsTable)

// Create bottom row with two equal columns
chartColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 1.0,
})
chartColumn.Add(widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Sales Chart\n[Chart visualization would go here]",
}))

activityColumn := widget.NewColumnItem(ctx, widget.ColumnItemOptions{
    Weight: 1.0,
})
activityColumn.Add(widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Recent Activity\n- User signup: Jane Smith\n- New order: #12345\n- Payment received: $1,299.99",
}))

// Add rows and columns to the content area
contentColumns.Add(statsRow)
contentColumns.Add(chartColumn)
contentColumns.Add(activityColumn)

// Add the nested columns to the main content column
contentColumn.Add(contentColumns)

// Add columns to the dashboard layout
dashboard.Add(sidebarColumn)
dashboard.Add(contentColumn)
```

## Best Practices

1. Use weights to create proportional column widths that make sense for your content
2. Consider how the layout will respond on different screen sizes
3. Use larger weights for content-heavy columns and smaller weights for navigation or supplementary content
4. Maintain a reasonable total sum of weights (e.g., if you have three columns with weights 1, 2, and 1, the total is 4, and each column takes up 1/4, 2/4, and 1/4 of the space)
5. Avoid extreme weight differences that might make some columns too narrow to be usable
6. Test your layout on different screen sizes to ensure all columns remain usable
7. Consider using nested columns for more complex layouts
8. Provide clear visual separation between columns
9. Ensure content within each column is properly aligned and formatted
10. Use consistent spacing within and between columns

## Related Components

- [Columns](./columns) - Container for organizing multiple ColumnItems in a row
- [Form](./form) - Container for organizing form elements
- `Grid` - Alternative for more complex grid-based layouts
