# Sourcetool Documentation MCP Server

Model Context Protocol (MCP) server for Sourcetool Documentation

This MCP server provides tools to access Sourcetool documentation content.

## Features

- **Documentation Access**: Fetch Sourcetool documentation content through MCP
- **Integrated Search**: Access documentation content directly from your AI assistant

## Usage

You can run this MCP server in two ways:

### 1. Using Docker

```sh
docker run -i --rm ghcr.io/trysourcetool/sourcetool-docs-mcp-server
```

### 2. Running Locally

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

## Development

- Source code is in the `src/` directory
- Tests are written using Vitest
- Run tests with `pnpm test`
