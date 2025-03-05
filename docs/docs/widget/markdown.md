---
sidebar_position: 12
---

# Markdown

The Markdown widget provides a way to display formatted text content using Markdown syntax, allowing for rich text presentation without complex HTML.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `body` | `string` | `""` | Markdown content to be rendered |

## Event Handling

The Markdown widget can emit events for link clicks:

```go
// Define a markdown widget with link click handler
markdown := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Content: "# Welcome\n\nClick [here](https://example.com) to learn more.",
    OnLinkClick: func(url string) {
        // Handle link click
        fmt.Printf("Link clicked: %s\n", url)
        // Open URL or navigate to internal route
    },
})
```

## Examples

### Basic Markdown

```go
package main

import (
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic markdown widget
    markdown := widget.NewMarkdown(ctx, widget.MarkdownOptions{
        Content: `# Welcome to Our Application
        
## Getting Started

This application helps you manage your tasks efficiently.

### Key Features

- **Task Management**: Create, edit, and organize tasks
- **Reminders**: Set reminders for important deadlines
- **Collaboration**: Share tasks with team members

[Learn more about our features](https://example.com/features)`,
    })
    
    // Add the markdown to your UI
    container.Add(markdown)
}
```

### Markdown with Code Highlighting

```go
// Create a markdown widget with code highlighting
codeMarkdown := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Content: `## Code Example

\`\`\`go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
\`\`\`

\`\`\`javascript
// JavaScript example
function greet(name) {
    console.log(\`Hello, \${name}!\`);
}
\`\`\``,
    Highlight: true,
    Theme: "github",
})
```

### Dynamic Markdown Content

```go
// Create a markdown widget with dynamic content
dynamicMarkdown := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Content: generateDocumentation(), // Function that returns markdown content
    AllowHTML: true,
    MaxHeight: "500px",
    OnLinkClick: func(url string) {
        if strings.HasPrefix(url, "internal://") {
            // Handle internal navigation
            route := strings.TrimPrefix(url, "internal://")
            navigateToRoute(route)
        } else {
            // Open external URL
            openExternalURL(url)
        }
    },
})
```

### Interactive Markdown

```go
// Create an interactive markdown widget
interactiveMarkdown := widget.NewMarkdown(ctx, widget.MarkdownOptions{
    Content: `# Interactive Demo

Click on the buttons below to see different sections:

- [Show Section 1](#section1)
- [Show Section 2](#section2)
- [Show Section 3](#section3)

<div id="section1" style="display:none">
## Section 1 Content
This is the content for section 1.
</div>

<div id="section2" style="display:none">
## Section 2 Content
This is the content for section 2.
</div>

<div id="section3" style="display:none">
## Section 3 Content
This is the content for section 3.
</div>`,
    AllowHTML: true,
    OnLinkClick: func(url string) {
        if strings.HasPrefix(url, "#section") {
            // Show the selected section and hide others
            showSection(url[1:])
        }
    },
})
```

## Best Practices

1. Use markdown for content that needs formatting but doesn't require complex interactivity
2. Keep the `sanitize` option enabled for user-generated content to prevent XSS attacks
3. Use the `allow_html` option cautiously, only when you trust the content source
4. Leverage code highlighting for technical documentation
5. Structure content with appropriate headings and lists for better readability
6. Use a consistent style for links, emphasis, and other formatting elements
7. Consider setting a `max_height` for long content to avoid excessive scrolling
8. Test your markdown rendering across different screen sizes
9. Use custom link handlers to create interactive documentation
10. Consider providing a way for users to copy code blocks to clipboard

## Related Components

- [TextArea](./textarea) - For multi-line text input
- [Form](./form) - Container for organizing form elements
- `HTMLViewer` - For displaying raw HTML content with more complex formatting
