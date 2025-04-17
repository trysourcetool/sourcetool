# mcp Directory

This directory contains modules related to the Model Context Protocol (MCP) integration for the Sourcetool project. Each subdirectory is intended to provide a specific MCP server or tool, making it easy to extend and add new MCP-related features in the future.

## Current structure

- `docs-mcp-server/`: An MCP server that provides access to Sourcetool documentation. It exposes a `getDocs` tool, which returns a list of document paths and their contents from the Sourcetool documentation site.

## Future plans

Additional MCP servers and tools will be added to this directory as the project evolves. Each new feature or integration should be placed in its own subdirectory, following the structure established by `docs-mcp-server`.

## Contribution

If you wish to add a new MCP server or tool, create a new subdirectory within `mcp/` and follow the conventions used in the existing modules. Please ensure to include appropriate documentation and tests for your additions.