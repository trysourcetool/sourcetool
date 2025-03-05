---
sidebar_position: 16
---

# Columns

The Columns widget provides a layout container that arranges child widgets horizontally in a row with equal spacing.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `columns` | `int32` | `2` | Number of columns in the layout |

## Event Handling

The Columns widget does not emit events directly, as it is a layout container. Events are handled by the child widgets placed within the columns.

```go
// Define a columns layout
columns := widget.NewColumns(ctx, widget.ColumnsOptions{
    Columns: 2, // Two columns layout
})

// Add child widgets to the columns
columns.Add(leftWidget)
columns.Add(rightWidget)
```

## Examples

### Basic Two-Column Layout

```go
package main

import (
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic two-column layout
    columns := widget.NewColumns(ctx, widget.ColumnsOptions{
        Columns: 2,
    })
    
    // Create widgets for the left column
    leftColumn := widget.NewContainer(ctx)
    leftTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
        Body: "## User Information",
    })
    nameInput := widget.NewTextInput(ctx, widget.TextInputOptions{
        Label: "Name",
        Required: true,
    })
    emailInput := widget.NewTextInput(ctx, widget.TextInputOptions{
        Label: "Email",
        Required: true,
    })
    leftColumn.Add(leftTitle)
    leftColumn.Add(nameInput)
    leftColumn.Add(emailInput)
    
    // Create widgets for the right column
    rightColumn := widget.NewContainer(ctx)
    rightTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
        Body: "## Preferences",
    })
    notificationsCheckbox := widget.NewCheckbox(ctx, widget.CheckboxOptions{
        Label: "Enable notifications",
    })
    themeSelect := widget.NewSelect(ctx, widget.SelectOptions{
        Label: "Theme",
        Options: []string{"Light", "Dark", "System"},
    })
    rightColumn.Add(rightTitle)
    rightColumn.Add(notificationsCheckbox)
    rightColumn.Add(themeSelect)
    
    // Add columns to the layout
    columns.Add(leftColumn)
    columns.Add(rightColumn)
    
    // Add the columns layout to your UI
    container.Add(columns)
}
```

### Three-Column Dashboard Layout

```go
// Create a three-column dashboard layout
dashboard := widget.NewColumns(ctx, widget.ColumnsOptions{
    Columns: 3,
})

// First column: Navigation
navColumn := widget.NewContainer(ctx)
navTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Navigation",
})
navLinks := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "- [Dashboard](#)\n- [Analytics](#)\n- [Settings](#)\n- [Profile](#)",
})
navColumn.Add(navTitle)
navColumn.Add(navLinks)

// Second column: Main content
contentColumn := widget.NewContainer(ctx)
contentTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Dashboard Overview",
})
statsTable := widget.NewTable(ctx, widget.TableOptions{
    Header: "Key Metrics",
    Data: statsData,
})
contentColumn.Add(contentTitle)
contentColumn.Add(statsTable)

// Third column: Notifications
notifColumn := widget.NewContainer(ctx)
notifTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Recent Notifications",
})
notifList := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "- New message from John\n- System update available\n- Meeting reminder: 3PM",
})
notifColumn.Add(notifTitle)
notifColumn.Add(notifList)

// Add columns to the dashboard
dashboard.Add(navColumn)
dashboard.Add(contentColumn)
dashboard.Add(notifColumn)
```

### Responsive Form Layout

```go
// Create a responsive form layout
formLayout := widget.NewColumns(ctx, widget.ColumnsOptions{
    Columns: 2,
})

// Left column: Personal information
personalInfoColumn := widget.NewContainer(ctx)
personalInfoTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Personal Information",
})
firstNameInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "First Name",
    Required: true,
})
lastNameInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Last Name",
    Required: true,
})
emailInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Email",
    Required: true,
})
personalInfoColumn.Add(personalInfoTitle)
personalInfoColumn.Add(firstNameInput)
personalInfoColumn.Add(lastNameInput)
personalInfoColumn.Add(emailInput)

// Right column: Address information
addressColumn := widget.NewContainer(ctx)
addressTitle := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Body: "## Address Information",
})
streetInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Street Address",
    Required: true,
})
cityInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "City",
    Required: true,
})
stateInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "State/Province",
    Required: true,
})
zipInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "ZIP/Postal Code",
    Required: true,
})
addressColumn.Add(addressTitle)
addressColumn.Add(streetInput)
addressColumn.Add(cityInput)
addressColumn.Add(stateInput)
addressColumn.Add(zipInput)

// Add columns to the form layout
formLayout.Add(personalInfoColumn)
formLayout.Add(addressColumn)

// Add a submit button below the columns
submitButton := widget.NewButton(ctx, widget.ButtonOptions{
    Label: "Submit",
    OnClick: func() {
        // Handle form submission
    },
})

// Add the form layout and submit button to your UI
container.Add(formLayout)
container.Add(submitButton)
```

## Best Practices

1. Use columns to create balanced, organized layouts
2. Consider the appropriate number of columns based on the content and screen size
3. Group related content within each column
4. Ensure each column has a similar amount of content to maintain visual balance
5. Consider how the layout will respond on different screen sizes
6. Use consistent spacing and alignment within and between columns
7. Provide clear visual separation between columns when necessary
8. Consider using column headers or titles to describe the content of each column
9. Avoid overcrowding columns with too many widgets
10. Test the layout on different screen sizes to ensure usability

## Related Components

- [ColumnItem](./column-item) - For defining individual items within a Columns layout
- [Form](./form) - Container for organizing form elements
- `Grid` - Alternative for more complex grid-based layouts
