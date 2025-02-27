# User Management Example

This example demonstrates how to use the Sourcetool Go SDK to create a simple user management application.

## Features

- List users with filtering capabilities
- Create new users
- Update existing users
- Role-based access control

## Prerequisites

- Go 1.22 or later
- Access to Sourcetool API

## Getting Started

1. Replace the API key in `main.go` with your own development API key:

```go
// Replace with your own API key for development
s := sourcetool.New("your_development_api_key")
```

2. Run the example:

```bash
go run .
```

3. The server will start at http://localhost:8081/

## Structure

- `main.go`: Sets up the HTTP server and Sourcetool UI
- `user.go`: Defines the User model and related functions

## Access Groups

The example demonstrates role-based access control with the following groups:

- `admin`: Has access to all pages
- `user_admin`: Has access to the user listing page
- `customer_support`: Has access to the user creation page

## Pages

- `/users`: Lists all users with filtering and update capabilities
- `/users/new`: Allows creating new users
