---
sidebar_position: 12
---

# Markdown

The Markdown widget provides a way to display formatted text content using Markdown syntax, allowing for rich text presentation without complex HTML.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `body` | `string` | `""` | Markdown content to be rendered |

## Examples

### Basic Markdown

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/markdown"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic markdown widget
        ui.Markdown(`# Welcome to Our Application
## Getting Started

This application helps you manage your tasks efficiently.

### Key Features

- **Task Management**: Create, edit, and organize tasks
- **Reminders**: Set reminders for important deadlines
- **Collaboration**: Share tasks with team members

[Learn more about our features](https://example.com/features)`)
    }
}
```

## Related Components

- [TextArea](./textarea) - For multi-line text input
- [Form](./form) - Container for organizing form elements
