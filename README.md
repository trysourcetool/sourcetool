# Sourcetool

**Build Internal Tools with Just Backend Code**

Sourcetool is an open-source internal tool builder that enables you to create full-featured applications without writing any frontend code.

## 🌟 About Sourcetool

We develop Sourcetool, an open-source internal tool builder that handles frontend complexities automatically, allowing developers to focus on implementing business logic in backend code only.

### Backend-First Development
Focus on your business logic while we handle the UI. Build complete internal tools using only Go code. No frontend expertise required.

### Type-Safe & Flexible
Built with Go's type system for reliability. Create robust applications with type-safe APIs and seamless integration.

## ✨ Features

| Feature | Description |
|---------|-------------|
| 💻 Backend-only, Code-first development | Build full-featured internal tools using only backend code with type-safe APIs, Git version control, and seamless integration with development workflows |
| 🎨 Rich UI components | Pre-built components (forms, tables, inputs, etc.) |
| 🔐 Granular permissions | Manage access to your internal tools with flexible group-based permissions |
| 🌐 Multiple environment support | Easily switch between different environments (development, staging, production) |

## 🏗️ Architecture

Sourcetool uses a three-tier architecture to connect your backend code directly to web browsers:

1. **Your Backend**
   - Contains your business logic and UI definitions
   - Written in Go using Sourcetool SDK

2. **Sourcetool Server**
   - Handles authentication and authorization
   - Manages WebSocket connections
   - Secures communication between components

3. **Web Browser**
   - Automatically generates the UI
   - Provides real-time interaction
   - No frontend code required

All components communicate bidirectionally in real-time through secure WebSocket connections.

### How It Works:
1. You define UI components in your backend code
2. Sourcetool Server handles auth & permissions
3. UI is automatically rendered in browser
4. User interactions return to your backend code

## 🎯 Components

Sourcetool provides UI components you can use directly from Go code:

| Component Type | Available Components |
|---------------|---------------------|
| 📝 Input Components | TextInput, TextArea, NumberInput, DateInput, DateTimeInput, TimeInput |
| 📋 Selection Components | Selectbox, MultiSelect, Radio, Checkbox, CheckboxGroup |
| 🔳 Layout Components | Columns, Form |
| 📊 Display Components | Markdown, Table |
| 🔘 Interactive Components | Button |

## 🚀 Get Started

1. **Get your API key**
   - Sign up at [Sourcetool Dashboard](https://auth.trysourcetool.com/)

2. **Install the Sourcetool SDK**
   ```bash
   go get github.com/trysourcetool/sourcetool-go
   ```

3. **Write your first internal tool**
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
       s := sourcetool.New(&sourcetool.Config{
           APIKey:   "your_api_key",
           Endpoint: "wss://hello.trysourcetool.com"
       })
       
       // Register pages
       s.Page("/users", "Users List", listUsersPage)
       
       if err := s.Listen(); err != nil {
           log.Fatal(err)
       }
   }
   ```

## ❓ FAQ

### What is Sourcetool?
Sourcetool is an open-source internal tool builder that enables you to build full-featured internal tools without writing any frontend code. It handles all frontend complexities automatically, allowing you to focus on implementing business logic in your backend code.

### Do I need frontend skills to use Sourcetool?
No. As an internal tool builder, Sourcetool lets you create complete applications using only Go. The system automatically handles all UI rendering and interactions without requiring any frontend code.

### What types of applications can I build with Sourcetool?
Admin panels, dashboards, data management systems, monitoring tools, and any application where development speed is more important than custom UI/UX.

### Is Sourcetool secure?
Yes, Sourcetool is designed with security in mind. You deploy and run Sourcetool applications on your own infrastructure, keeping your code and sensitive data within your control.

### Is Sourcetool free to use?
Check out the [Sourcetool website](https://trysourcetool.com/) for pricing information.

## 📚 Resources

- [Documentation](https://docs.trysourcetool.com)
- [GitHub Repository](https://github.com/trysourcetool/sourcetool)
- [Security Policy](SECURITY.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

---

<div align="center">
Made with ❤️ by the Sourcetool Team
</div>
