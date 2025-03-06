---
sidebar_position: 15
---

# Form

The Form widget provides a container for organizing form elements with built-in submission handling.

## States

| State | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | `bool` | `False` | Current submission state of the form |

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `button_label` | `string` | `"Submit"` | Text displayed on the form's submit button |
| `button_disabled` | `bool` | `False` | Whether the submit button is disabled |
| `clear_on_submit` | `bool` | `False` | Whether to clear form inputs after submission |

## Examples

### Basic Form

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/form"
    "github.com/trysourcetool/sourcetool-go/textinput"
)

func main() {
    func page(ui sourcetool.UIBuilder) error {
        // Create a basic form
        form, submitted := ui.Form("Update", form.ClearOnSubmit(true))

        // Add form elements
        nameInput := form.TextInput("Name", textinput.Placeholder("Enter your name"), textinput.DefaultValue(defaultName), textinput.Required(true))
        emailInput := form.TextInput("Email", textinput.Placeholder("Enter your email"), textinput.DefaultValue(defaultEmail))

        if submitted {
            user := User{
                Name:   formName,
                Email:  formEmail
            }
            if err := createUser(&user); err != nil {
                return err
            }
        }
    }
}
```

## Related Components

- [TextInput](./text-input) - For single-line text input
- [TextArea](./textarea) - For multi-line text input
- [Checkbox](./checkbox) - For boolean input
- [Select](./select) - For selecting from predefined options
- [DateInput](./date-input) - For date selection
