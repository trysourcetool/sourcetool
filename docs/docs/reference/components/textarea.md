---
sidebar_position: 2
---

# TextArea

The TextArea widget provides a multi-line text input field for entering longer text content.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `string` | `None` | Current text content of the textarea |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `label` | `string` | `""` | Label text displayed above the textarea |
| `placeholder` | `string` | `""` | Placeholder text displayed when the textarea is empty |
| `default_value` | `string` | `None` | Initial value of the textarea |
| `required` | `bool` | `False` | Whether the textarea input is required |
| `disabled` | `bool` | `False` | Whether the textarea is disabled |
| `max_length` | `int32` | `None` | Maximum number of characters allowed |
| `min_length` | `int32` | `None` | Minimum number of characters allowed |
| `max_lines` | `int32` | `None` | Maximum number of lines allowed |
| `min_lines` | `int32` | `None` | Minimum number of lines allowed |
| `auto_resize` | `bool` | `False` | Whether the textarea should automatically resize based on content |

## Examples

### Basic TextArea

```go
package main

import (
    "github.com/sourcetool/sourcetool-go"
    "github.com/sourcetool/sourcetool-go/textarea"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic textarea
        textarea := ui.TextArea("Comments", textarea.Placeholder("Enter your comments here"))
    }
}
```

### Disabled TextArea

```go
// Create a disabled textarea
disabledTextarea := ui.TextArea("System Log", textarea.Placeholder("System log content that cannot be edited"), textarea.Disabled(true))
```

## Related Components

- [TextInput](./text-input) - For single-line text input
- [Markdown](./markdown) - For displaying formatted text content
