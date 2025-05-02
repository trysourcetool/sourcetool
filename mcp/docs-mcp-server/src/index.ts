import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
export const server = new McpServer({
  name: 'sourcetool-docs-mcp',
  version: '0.1.11',
});

server.tool(
  'getDocs',
  'Get a list of sourcetool document paths and contents',
  {},
  async () => {
    const docs = await fetch(
      'https://docs.trysourcetool.com/json/docs.json',
    ).then(
      (res) =>
        res.json() as Promise<
          { path: string; title: string; content: string }[]
        >,
    );

    return {
      content: docs.map((doc) => ({
        type: 'text',
        path: doc.path,
        title: doc.title,
        text: doc.content,
      })),
    };
  },
);

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  // Output to standard error output since standard output will be interpreted as a server response
  console.error('MCP Server running on stdio');
}

main().catch((error) => {
  console.error('Fatal error in main():', error);
  process.exit(1);
});
