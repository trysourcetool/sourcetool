# Sourcetool Documentation MCP Server

Model Context Protocol (MCP) server for Sourcetool Documentation

This MCP server provides tools to access Sourcetool documentation content and enable seamless documentation search capabilities.

## Features

- **Documentation Access**: Fetch and read Sourcetool documentation content through MCP
- **Integrated Search**: Access documentation content directly from your AI assistant
- **Efficient Content Delivery**: Optimized for quick access to documentation resources

## Installation

To add this MCP server to your AI assistant configuration, add the following to your MCP config file.

```json
{
  "mcpServers": {
    "trysourcetool.sourcetool-docs-mcp-server": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "ghcr.io/trysourcetool/sourcetool-docs-mcp-server"]
    }
  }
}
```

## Tools

### getDocs

Retrieves a list of Sourcetool documentation paths and their contents.

```typescript
getDocs() -> {
  content: Array<{
    type: 'text',
    path: string,
    title: string,
    text: string
  }>
}
```

Example usage:
- "Show me the documentation about Sourcetool features"
- "Search for configuration options in Sourcetool docs"

## Development

### Prerequisites

1. Install `pnpm` from [pnpm.io](https://pnpm.io/installation)
2. Node.js 16 or newer

### Local Development Setup

1. Install dependencies:
   ```sh
   pnpm install
   ```

2. Build the project:
   ```sh
   pnpm build
   ```

3. Run the server:
   ```sh
   node dist/index.js
   ```

- Source code is in the `src/` directory
- Tests are written using Vitest
- Run tests with `pnpm test`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
