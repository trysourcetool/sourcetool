---
sidebar_position: 1
---

# Quickstart

## Introduction

Sourcetool is a powerful toolkit for building internal tools with just backend code. It provides a rich set of UI components and handles all the frontend complexities, allowing developers to focus on business logic implementation.

### Key Features

- **Backend-Only Development**: Build full-featured internal tools without writing any frontend code
- **Rich UI Components**: Comprehensive set of pre-built components (forms, tables, inputs, etc.)
- **Real-time Updates**: Built-in WebSocket support for live data synchronization
- **Type-Safe**: Fully typed API for reliable development
- **Flexible Backend**: Freedom to implement any business logic in pure Go

## Installation

### Prerequisites

1. Get your API key from [Sourcetool Dashboard](https://trysourcetool.com)
2. Install Go 1.18 or later

### Install the SDK

```bash
go get github.com/trysourcetool/sourcetool-go
```

## Creating Your First Application

Let's create a simple user management page to get started with Sourcetool:

```go
package main

import (
    "log"
    
    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/textinput"
    "github.com/trysourcetool/sourcetool-go/table"
)

// Sample user data structure
type User struct {
    ID    string
    Name  string
    Email string
    Role  string
}

// Mock function to simulate fetching users from a database
func listUsers(nameFilter string) ([]User, error) {
    // In a real application, you would query your database here
    return []User{
        {ID: "1", Name: "John Doe", Email: "john@example.com", Role: "Admin"},
        {ID: "2", Name: "Jane Smith", Email: "jane@example.com", Role: "User"},
    }, nil
}

// Page handler function for listing users
func listUsersPage(ui sourcetool.UIBuilder) error {
    ui.Markdown("## Users")

    // Search form
    name := ui.TextInput("Name", textinput.Placeholder("Enter name to search"))
    
    // Display users table
    users, err := listUsers(name)
    if err != nil {
        return err
    }
    
    ui.Table(users, table.Header("Users List"))
    
    return nil
}

func main() {
    // Initialize Sourcetool with your API key
    st := sourcetool.New("your-api-key")
    
    // Register pages
    st.Page("/users", "Users List", listUsersPage)
    
    // Start the server
    if err := st.Listen(); err != nil {
        log.Fatal(err)
    }
}
```

## Understanding the Basics

### Pages

Pages are the main building blocks of your application. Each page is defined by:

- A route (e.g., `/users`)
- A title (e.g., "Users List")
- A handler function that builds the UI

```go
st.Page("/users", "Users List", listUsersPage)
```

### UI Components

Sourcetool provides a wide range of UI components that you can use to build your pages:

#### Input Components
- TextInput: Single-line text input
- TextArea: Multi-line text input
- NumberInput: Numeric input with validation
- DateInput: Date picker
- DateTimeInput: Date and time picker
- TimeInput: Time picker

#### Selection Components
- Selectbox: Single-select dropdown
- MultiSelect: Multi-select dropdown
- Radio: Radio button group
- Checkbox: Single checkbox
- CheckboxGroup: Group of checkboxes

#### Layout Components
- Columns: Multi-column layout
- Form: Form container with submit button
- Table: Data table with sorting and selection

#### Display Components
- Markdown: Formatted text display

#### Interactive Components
- Button: Clickable button

### Component Options

Each component supports various options for customization:

```go
// TextInput with options
ui.TextInput("Username",
    textinput.Placeholder("Enter username"),
    textinput.Required(true),
    textinput.MaxLength(50),
)

// Table with options
ui.Table(data,
    table.Header("Users"),
    table.OnSelect(table.SelectionBehaviorRerun),
    table.RowSelection(table.SelectionModeSingle),
)
```

## Creating a Form

Here's an example of creating a form with multiple fields and validation:

```go
func createUserPage(ui sourcetool.UIBuilder) error {
    form, submitted := ui.Form("Create User", form.ClearOnSubmit(true))
    
    name := form.TextInput("Name", 
        textinput.Required(true),
        textinput.MinLength(2),
        textinput.MaxLength(50),
    )
    
    email := form.TextInput("Email",
        textinput.Required(true),
        textinput.Placeholder("user@example.com"),
    )
    
    role := form.Selectbox("Role",
        selectbox.Options("Admin", "User", "Guest"),
        selectbox.Required(true),
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

## Next Steps

Now that you've created your first Sourcetool application, you can:

1. Learn more about [Sourcetool concepts](../concepts/pages)
2. Explore the [component reference](../reference/components)
<!-- TODO: Add best practices documentation -->
<!-- 3. Check out [best practices](../best-practices) for building robust applications -->

For detailed documentation and examples, visit our [documentation site](https://docs.trysourcetool.com).
