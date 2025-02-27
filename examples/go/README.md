# Sourcetool Go Examples

This directory contains various examples demonstrating how to use the Sourcetool Go SDK.

## Examples

### [Hello World](./hello-world)

A simple "Hello World" example demonstrating the basic usage of the Sourcetool Go SDK:
- Markdown rendering
- Text input with placeholder
- Button with click handling
- Conditional UI rendering based on input

### [User Management](./user-management)

A complete example of a user management application that demonstrates:
- Creating UI components (forms, tables, input fields)
- Implementing role-based access control
- Setting up page routing
- Handling user interactions

## Adding New Examples

To add a new example:

1. Create a new directory for your example (e.g., `data-visualization`, `authentication`, etc.)
2. Set up a proper Go module structure with `go.mod` file
3. Include a comprehensive README.md explaining the example
4. Implement the example code

## Running Examples

Each example is a standalone Go module. To run an example:

1. Navigate to the example directory:
   ```bash
   cd user-management
   ```

2. Run the example:
   ```bash
   go run .
   ```

## Example Structure

Each example should follow this structure:

```
example-name/
├── go.mod           # Module definition
├── go.sum           # Dependency checksums
├── main.go          # Main application entry point
├── README.md        # Documentation
└── ...              # Additional files specific to the example
