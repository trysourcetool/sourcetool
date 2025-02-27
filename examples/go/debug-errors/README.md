# Debug Errors Example

This example demonstrates how to return errors from page functions in the Sourcetool Go SDK and how these errors are handled.

## Features

- Intentionally returning errors from page functions
- Custom error types and error handling patterns
- Panic recovery demonstrations
- Error middleware for HTTP handlers

## Error Types Demonstrated

1. **Validation Errors**: Form validation errors with field-specific messages that are returned from the page function
2. **Database Errors**: Simulated database operation failures returned as errors
3. **Runtime Panics**: Various panic scenarios including:
   - Explicit panics
   - Nil pointer dereference
   - Index out of range
   - Division by zero

## How Errors Are Handled

When a page function returns an error:
1. Sourcetool catches the error
2. The error is logged
3. An appropriate error message is displayed to the user
4. The page is re-rendered with the error information

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

3. The server will start at http://localhost:8083/

## Usage

The example provides a UI with several sections:

1. **Trigger Errors**: Select an error type from the dropdown and click the button to trigger it
2. **Form with Validation Errors**: Submit the form with invalid data to see validation errors
3. **Error Recovery Demo**: Test panic recovery functionality

## Debugging Tips

- Check the server logs for detailed error information
- Use the middleware's panic recovery to prevent server crashes
- Examine how different error types are handled and displayed in the UI

## Structure

- `main.go`: Contains the error demonstration code and server setup
