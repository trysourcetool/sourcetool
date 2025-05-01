---
sidebar_position: 1
---

# Pages

Pages are the fundamental building blocks of your Sourcetool application. They define the different screens and views that users can navigate to and interact with.

## What is a Page?

In Sourcetool, a page is:

1. A unique URL route (e.g., `/users`, `/dashboard`, `/settings`)
2. A handler function that builds the UI for that route
3. A set of UI components that users interact with

Pages are defined in your Go code and rendered by Sourcetool Cloud in the user's browser.

## Defining Pages

Pages are defined using the `Page` method of the Sourcetool instance:

```go
st := sourcetool.New("your-api-key")

// Register a page
st.Page("/users", "Users", usersPageHandler)
```

The `Page` method takes three parameters:

1. **Route**: The URL path for the page (e.g., `/users`)
2. **Title**: The title displayed in the browser and navigation
3. **Handler**: A function that builds the UI for the page

## Page Handler

The page handler is a function that takes a `UIBuilder` and returns an error:

```go
func usersPageHandler(ui sourcetool.UIBuilder) error {
    // Build the UI for the page
    ui.Markdown("# Users")
    
    // Add components, handle user input, etc.
    
    return nil
}
```

The handler function is called:

1. When a user first navigates to the page
2. When a user interacts with a component on the page (e.g., clicks a button)
3. When the page is refreshed

## Page Lifecycle

Pages in Sourcetool follow a specific lifecycle:

1. **Initialization**: When a user navigates to a page, Sourcetool initializes the page state
2. **Rendering**: The page handler builds the UI using the `UIBuilder`
3. **Interaction**: The user interacts with the UI components
4. **Re-rendering**: When a user interacts with a component, the page handler is called again with the updated state
5. **Cleanup**: When a user navigates away from the page, Sourcetool cleans up the page state

## Page Routing

Sourcetool handles routing automatically based on the routes you define:

```go
// Basic routes
st.Page("/home", "Home", homePageHandler)
st.Page("/users", "Users", usersPageHandler)
st.Page("/settings", "Settings", settingsPageHandler)

// Routes with parameters
st.Page("/users/:id", "User Details", userDetailsPageHandler)
```

### Route Parameters

You can define dynamic route parameters using the `:param` syntax:

```go
st.Page("/users/:id", "User Details", userDetailsPageHandler)
```

These parameters are accessible in your page handler:

```go
func userDetailsPageHandler(ui sourcetool.UIBuilder) error {
    // Get the route parameters
    params := sourcetool.RouteParams(ui.Context())
    userID := params["id"]
    
    // Fetch user data
    user, err := fetchUser(userID)
    if err != nil {
        return err
    }
    
    // Build the UI
    ui.Markdown(fmt.Sprintf("# User: %s", user.Name))
    
    return nil
}
```

## Navigation Between Pages

Sourcetool provides several ways to navigate between pages:

### Link Component

```go
ui.Link("Go to Users", "/users")
```

### Button with Navigation

```go
ui.Button("View Users", button.OnClick(func() {
    sourcetool.Navigate(ui.Context(), "/users")
}))
```

### Programmatic Navigation

```go
func handleAction(ui sourcetool.UIBuilder) error {
    // Perform some action
    
    // Navigate to another page
    sourcetool.Navigate(ui.Context(), "/users")
    
    return nil
}
```

## Page State Management

Each page maintains its own state, which is persisted across re-renders:

### Local State

Components maintain their state automatically:

```go
// The value of this input is preserved across re-renders
name := ui.TextInput("Name")
```

### Shared State

You can share state between components using Go variables:

```go
func userPageHandler(ui sourcetool.UIBuilder) error {
    // Define shared state
    var selectedUserID string
    
    // User list with selection
    users, err := fetchUsers()
    if err != nil {
        return err
    }
    
    selection := ui.Table(users, 
        table.OnSelect(func(row int) {
            selectedUserID = users[row].ID
        }),
    )
    
    // If a user is selected, show details
    if selectedUserID != "" {
        user, err := fetchUserDetails(selectedUserID)
        if err != nil {
            return err
        }
        
        ui.Markdown(fmt.Sprintf("## User Details\n\nName: %s\nEmail: %s", 
            user.Name, user.Email))
    }
    
    return nil
}
```

## Page Layout

Pages can have different layouts:

### Single Column Layout

```go
func singleColumnPage(ui sourcetool.UIBuilder) error {
    ui.Markdown("# Single Column")
    ui.TextInput("Input 1")
    ui.TextInput("Input 2")
    return nil
}
```

### Multi-Column Layout

```go
func multiColumnPage(ui sourcetool.UIBuilder) error {
    ui.Markdown("# Multi-Column Layout")
    
    cols := ui.Columns(2) // Create a 2-column layout
    
    // First column
    cols[0].Markdown("## Column 1")
    cols[0].TextInput("Input 1")
    
    // Second column
    cols[1].Markdown("## Column 2")
    cols[1].TextInput("Input 2")
    
    return nil
}
```

## Page Access Control

You can control access to pages using groups:

```go
st.AccessGroups("admin")
usersGroup := st.Group("/users")
{
  usersGroup.AccessGroups("user_group")
  usersGroup.Page("/", "Users", usersPageHandler)
}
```

## Page Organization

For larger applications, you can organize pages into sections:

```go
// Main pages
st.Page("/home", "Home", homePageHandler)
st.Page("/dashboard", "Dashboard", dashboardPageHandler)

// User management pages
st.Page("/users", "Users", usersPageHandler)
st.Page("/users/:id", "User Details", userDetailsPageHandler)
st.Page("/users/new", "New User", newUserPageHandler)

// Settings pages
st.Page("/settings", "Settings", settingsPageHandler)
st.Page("/settings/profile", "Profile Settings", profileSettingsPageHandler)
st.Page("/settings/security", "Security Settings", securitySettingsPageHandler)
```

## Best Practices

### Keep Pages Focused

Each page should have a single responsibility. For example:

- A list page for displaying multiple items
- A detail page for viewing a single item
- A form page for creating or editing an item

### Handle Errors Gracefully

Always handle errors in your page handlers:

```go
func userPageHandler(ui sourcetool.UIBuilder) error {
    users, err := fetchUsers()
    if err != nil {
        // Display an error message to the user
        ui.Markdown("# Error")
        ui.Markdown(fmt.Sprintf("Failed to load users: %s", err.Error()))
        return err
    }
    
    ui.Markdown("# Users")
    ui.Table(users)
    
    return nil
}
```

### Organize Complex Pages

For complex pages, break down the UI building into smaller functions:

```go
func dashboardPageHandler(ui sourcetool.UIBuilder) error {
    ui.Markdown("# Dashboard")
    
    cols := ui.Columns(2)
    
    buildUsersSummary(cols[0])
    buildActivityFeed(cols[1])
    
    return nil
}

func buildUsersSummary(ui sourcetool.UIBuilder) {
    ui.Markdown("## Users Summary")
    // Build users summary UI
}

func buildActivityFeed(ui sourcetool.UIBuilder) {
    ui.Markdown("## Recent Activity")
    // Build activity feed UI
}
```

## Next Steps

Now that you understand pages, learn about:

- [Environments](./environments) for deploying your application
- [Organizations](./organizations) for managing users and access control
- [Components](../reference/components) for building rich UIs
