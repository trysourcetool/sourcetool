# Sourcetool Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/trysourcetool/sourcetool-go.svg)](https://pkg.go.dev/github.com/trysourcetool/sourcetool-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/trysourcetool/sourcetool-go)](https://goreportcard.com/report/github.com/trysourcetool/sourcetool-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Sourcetool Go SDK is a powerful toolkit for building internal tools with just backend code. It provides a rich set of UI components and handles all the frontend complexities, allowing developers to focus on business logic implementation.

## Features

- 🚀 **Backend-Only Development**: Build full-featured internal tools without writing any frontend code
- 🎨 **Rich UI Components**: Comprehensive set of pre-built components (forms, tables, inputs, etc.)
- ⚡ **Real-time Updates**: Built-in WebSocket support for live data synchronization
- 🔒 **Type-Safe**: Fully typed API for reliable development
- 🛠 **Flexible Backend**: Freedom to implement any business logic in pure Go

## Installation

```bash
go get github.com/trysourcetool/sourcetool-go
```

## Quick Start

### Prerequisites

1. Get your API key from [Sourcetool Dashboard](https://trysourcetool.com)
2. Install Go 1.18 or later

### Basic Example

Here's a simple example of creating a user management page:

```go
package main

import (
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/textinput"
    "github.com/trysourcetool/sourcetool-go/table"
)

func listUsersPage(ui sourcetool.UIBuilder) error {
    ui.Markdown("## Users")

    // Search form
    name := ui.TextInput("Name", textinput.WithPlaceholder("Enter name to search"))
    
    // Display users table
    users, err := listUsers(name)
    if err != nil {
        return err
    }
    
    ui.Table(users, table.WithHeader("Users List"))
    
    return nil
}

func main() {
    s := sourcetool.New(&sourcetool.Config{
        APIKey:   "your_api_key",
        Endpoint: "wss://your-sourcetool-instance"  // Your self-hosted Sourcetool endpoint
    })
    
    // Register pages
    s.Page("/users", "Users List", listUsersPage)
    
    if err := s.Listen(); err != nil {
        log.Fatal(err)
    }
}
```

## Available Components

Sourcetool provides a wide range of UI components:

### Input Components
- TextInput: Single-line text input
- TextArea: Multi-line text input
- NumberInput: Numeric input with validation
- DateInput: Date picker
- DateTimeInput: Date and time picker
- TimeInput: Time picker

### Selection Components
- Selectbox: Single-select dropdown
- MultiSelect: Multi-select dropdown
- Radio: Radio button group
- Checkbox: Single checkbox
- CheckboxGroup: Group of checkboxes

### Layout Components
- Columns: Multi-column layout
- Form: Form container with submit button
- Table: Data table with sorting and selection

### Display Components
- Markdown: Formatted text display

### Interactive Components
- Button: Clickable button

## Component Options

Each component supports various options for customization:

```go
// TextInput with options
ui.TextInput("Username",
    textinput.WithPlaceholder("Enter username"),
    textinput.WithRequired(true),
    textinput.WithMaxLength(50),
)

// Table with options
ui.Table(data,
    table.WithHeader("Users"),
    table.WithOnSelect(table.OnSelectRerun),
    table.WithRowSelection(table.RowSelectionSingle),
)
```

## Advanced Usage

### Error Handling

Sourcetool provides robust error handling:

```go
func userProfilePage(ui sourcetool.UIBuilder) error {
    user, err := fetchUserProfile()
    if err != nil {
        // Display error message to the user
        ui.Markdown("⚠️ Error: Failed to load user profile")
        return err
    }
    
    ui.Markdown(fmt.Sprintf("## Welcome, %s!", user.Name))
    return nil
}
```

### Complex Forms

Example of a form with multiple fields and validation:

```go
func createUserPage(ui sourcetool.UIBuilder) error {
    form, submitted := ui.Form("Create User", form.WithClearOnSubmit(true))
    
    name := form.TextInput("Name", 
        textinput.WithRequired(true),
        textinput.WithMinLength(2),
        textinput.WithMaxLength(50),
    )
    
    email := form.TextInput("Email",
        textinput.WithRequired(true),
        textinput.WithPlaceholder("user@example.com"),
    )
    
    role := form.Selectbox("Role",
        selectbox.WithOptions("Admin", "User", "Guest"),
        selectbox.WithRequired(true),
    )
    
    if submitted {
        user := User{
            Name: name,
            Email: email,
            Role: role.Value,
        }
        if err := createUser(&user); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Documentation

For detailed documentation and examples, visit our [documentation site](https://docs.trysourcetool.com).

## Upgrading

When upgrading to a new version:

1. Check the [changelog](https://github.com/trysourcetool/sourcetool-go/releases) for breaking changes
2. Update your dependencies:
   ```bash
   go get -u github.com/trysourcetool/sourcetool-go
   ```
3. Run tests to ensure compatibility
4. Review deprecated features and update accordingly

## Development

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/trysourcetool/sourcetool-go.git
   cd sourcetool-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/websocket/...
```

## Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Follow standard Go formatting guidelines
4. Add tests for new features
5. Update documentation as needed
6. Commit your changes (`git commit -m 'Add some amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
