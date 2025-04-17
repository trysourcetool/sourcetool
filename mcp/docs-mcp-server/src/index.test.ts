import { describe, test, expect, vi, afterEach, beforeEach } from 'vitest';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { InMemoryTransport } from '@modelcontextprotocol/sdk/inMemory.js';
import { server } from './index.js';

describe('getDocs', () => {
  const result = [
    {
      path: 'docs/index',
      title: 'Docs',
      content: 'Docs',
    },
  ];
  const want = [
    {
      type: 'text',
      path: 'docs/index',
      title: 'Docs',
      text: 'Docs',
    },
  ];

  beforeEach(async () => {
    vi.spyOn(global, 'fetch').mockImplementation(
      async () => new Response(JSON.stringify(result), { status: 200 }),
    );
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  test('Test whether the document can be retrieved correctly', async () => {
    const client = new Client({
      name: 'test client',
      version: '0.1.0',
    });

    const [clientTransport, serverTransport] =
      InMemoryTransport.createLinkedPair();

    await Promise.all([
      client.connect(clientTransport),
      server.connect(serverTransport),
    ]);

    const result = await client.callTool({
      name: 'getDocs',
      arguments: {},
    });

    expect(result.content).toEqual(want);
  });
});
