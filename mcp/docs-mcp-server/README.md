# docs-mcp-server

This module provides an MCP (Model Context Protocol) server for Sourcetool documentation. It exposes a `getDocs` tool that returns a list of document paths and their contents from the Sourcetool documentation site.

## Features
- Implements an MCP server using `@modelcontextprotocol/sdk`.
- Provides the `getDocs` tool to fetch and return documentation data from https://docs.trysourcetool.com/json/docs.json.

## Setup

1. Install dependencies:
   ```sh
   pnpm install
   ```

2. Build the project:
   ```sh
   pnpm build
   ```

## Usage

You can run the server as a binary after building:

```sh
pnpm build
node build/index.js
```

Or use the provided binary:

```sh
pnpm build
./build/index.js
```

## Testing

Run tests using Vitest:

```sh
pnpm test
```

## Development
- Source code is in the `src/` directory.
- TypeScript configuration is in `tsconfig.json`.
- Tests are located in `src/index.test.ts`.