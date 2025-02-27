# Hello World Example

This is a simple "Hello World" example demonstrating the basic usage of the Sourcetool Go SDK.

## Features

- Simple UI with text input and button
- Basic interaction handling
- Minimal server setup

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

3. The server will start at http://localhost:8082/

## Structure

- `main.go`: Sets up a minimal HTTP server and Sourcetool UI with a simple page

## Components Demonstrated

- Markdown rendering
- Text input with placeholder
- Button with click handling
- Conditional UI rendering based on input
