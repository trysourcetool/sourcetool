---
sidebar_position: 2
---

# How It Works

This guide explains the architecture and core concepts of Sourcetool to help you understand how the system works behind the scenes.

## Architecture Overview

Sourcetool follows a unique architecture that allows you to build full-featured web applications using only backend code. Here's how it works:

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│                 │      │                 │      │                 │
│  Your Go Code   │◄────►│  Sourcetool SDK │◄────►│ Sourcetool Cloud│
│                 │      │                 │      │                 │
└─────────────────┘      └─────────────────┘      └─────────────────┘
                                                          ▲
                                                          │
                                                          ▼
                                                  ┌─────────────────┐
                                                  │                 │
                                                  │    End Users    │
                                                  │                 │
                                                  └─────────────────┘
```

1. **Your Go Code**: You write business logic and UI definitions in Go
2. **Sourcetool SDK**: Translates your UI definitions into frontend components
3. **Sourcetool Cloud**: Hosts the frontend application and manages user sessions
4. **End Users**: Access your application through a web browser

## Backend to Frontend Bridge

Sourcetool eliminates the need to write frontend code by providing a bridge between your Go backend and the frontend UI:

1. **UI Builder API**: You use the `UIBuilder` interface to define UI components
2. **State Serialization**: The SDK serializes the UI state and sends it to Sourcetool Cloud
3. **Component Rendering**: Sourcetool Cloud renders the components in the user's browser
4. **Event Handling**: User interactions are sent back to your Go code for processing

## Page Lifecycle

When a user accesses a page in your application, the following sequence occurs:

1. **Request**: The user requests a page (e.g., `/users`)
2. **Handler Execution**: Sourcetool calls your page handler function
3. **UI Building**: Your handler builds the UI using the `UIBuilder` interface
4. **Rendering**: The UI is rendered in the user's browser
5. **Interaction**: When the user interacts with the UI, your handler is called again with the updated state

```go
func userListPage(ui sourcetool.UIBuilder) error {
    // 1. Define UI components
    ui.Markdown("# Users")
    
    // 2. Handle user input
    searchTerm := ui.TextInput("Search", textinput.Placeholder("Search users..."))
    
    // 3. Process data based on input
    users, err := fetchUsers(searchTerm)
    if err != nil {
        return err
    }
    
    // 4. Display results
    ui.Table(users)
    
    return nil
}
```

## Session Management

Sourcetool manages user sessions automatically:

1. **Session Creation**: When a user accesses your application, a new session is created
2. **State Persistence**: The state of UI components is persisted across page reloads
3. **Session Expiry**: Sessions expire after a period of inactivity

## Real-time Updates

Sourcetool supports real-time updates through WebSockets:

1. **WebSocket Connection**: A WebSocket connection is established between the user's browser and Sourcetool Cloud
2. **Event Streaming**: Events are streamed in real-time between the frontend and backend
3. **UI Updates**: The UI is updated automatically when the underlying data changes

## Component Rendering

When you define a UI component in your Go code, Sourcetool:

1. **Creates a Component Definition**: Converts your Go code into a component definition
2. **Assigns a Unique ID**: Each component gets a unique identifier
3. **Tracks State**: Maintains the state of the component across requests
4. **Renders the Component**: Renders the component in the user's browser

## Data Flow

Data flows through your application as follows:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │     │             │
│  User Input │────►│ Go Handler  │────►│  Data Store │────►│ UI Rendering│
│             │     │             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
      ▲                                                            │
      │                                                            │
      └────────────────────────────────────────────────────────────┘
```

1. **User Input**: The user interacts with the UI
2. **Go Handler**: Your handler function processes the input
3. **Data Store**: Your code interacts with your data store (database, API, etc.)
4. **UI Rendering**: The updated UI is rendered based on the new data
5. **Feedback Loop**: The user sees the updated UI and can provide new input

## Environments and Deployment

Sourcetool supports multiple environments for your applications:

1. **Development**: For local development and testing
2. **Staging**: For pre-production testing
3. **Production**: For live applications

Each environment can have its own:

- API keys
- Configuration settings
- Access controls
- Domain names

## Organizations and Access Control

Sourcetool provides organization-level access control:

1. **Organizations**: Group users and applications
2. **User Roles**: Assign different roles to users (admin, developer, member)
3. **Groups**: Create groups of users for fine-grained access control
4. **Pages**: Control which groups can access which pages

## API Keys

API keys are used to authenticate your application with Sourcetool Cloud:

1. **Development Keys**: For local development
2. **Production Keys**: For production environments
3. **Custom Keys**: Create custom keys for specific use cases or services

## Next Steps

Now that you understand how Sourcetool works, you can:

1. Learn more about [Pages](../concepts/pages) and how they structure your application
2. Explore [Environments](../concepts/environments) for deploying your application
3. Understand [Organizations](../concepts/organizations) for managing users and access control
