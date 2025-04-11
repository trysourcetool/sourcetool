import { expect, test } from 'vitest';
import { newPageManager, Page } from '../page';
import { v4 as uuidv4 } from 'uuid';
import { MockClient } from '../websocket/mock/websocket';
import { createSessionManager, newSession } from '../session';
import { Runtime } from '../runtime';
import { create } from '@bufbuild/protobuf';
import {
  CloseSessionSchema,
  InitializeClientSchema,
  RerunPageSchema,
} from '../pb/websocket/v1/message_pb';
test('initialize client', () => {
  const pages: { [pageId: string]: Page } = {};
  const pageId = uuidv4();

  let handlerCalled = false;
  const testPage = new Page(
    pageId,
    'Test Page',
    '/test',
    [],
    async () => {
      handlerCalled = true;
    },
    [],
  );

  pages[pageId] = testPage;

  const mockWS = new MockClient();
  const sessionManager = createSessionManager();
  const pageManager = newPageManager(pages);
  const runtime = new Runtime(mockWS, sessionManager, pageManager);

  // Create test message
  const sessionId = uuidv4();
  const initClient = create(InitializeClientSchema, {
    sessionId,
    pageId,
  });

  // Register handler
  mockWS.registerHandler((message) => {
    if (message.type.case === 'initializeClient') {
      return runtime.handleInitializeClient(message.type.value);
    }
  });

  // Send message
  mockWS.enqueue(uuidv4(), initClient);

  // Verify that session was created
  const session = runtime.sessionManager.getSession(sessionId);
  if (!session) {
    throw new Error('Session not created');
  }

  expect(handlerCalled).toBe(true);
});

test('rerun page', () => {
  const pages: { [pageId: string]: Page } = {};
  const pageId = uuidv4();
  const sessionId = uuidv4();

  let handlerCallCount = 0;

  const testPage = new Page(
    pageId,
    'Test Page',
    '/test',
    [],
    async () => {
      handlerCallCount++;
    },
    [],
  );

  pages[pageId] = testPage;

  const mockWS = new MockClient();
  const sessionManager = createSessionManager();
  const pageManager = newPageManager(pages);
  const runtime = new Runtime(mockWS, sessionManager, pageManager);

  // Initialize session
  const session = newSession(sessionId, pageId);
  sessionManager.setSession(session);

  // Create test message
  const rerunPage = create(RerunPageSchema, {
    sessionId,
    pageId,
  });

  // Register handler
  mockWS.registerHandler((message) => {
    if (message.type.case === 'rerunPage') {
      return runtime.handleRerunPage(message.type.value);
    }
  });

  // Send message
  mockWS.enqueue(uuidv4(), rerunPage);

  // Verify that page was rerun
  expect(handlerCallCount).toBe(1);
});

test('close session', () => {
  const pageId = uuidv4();
  const sessionId = uuidv4();

  const mockWS = new MockClient();
  const sessionManager = createSessionManager();
  const pageManager = newPageManager();
  const runtime = new Runtime(mockWS, sessionManager, pageManager);

  // Initialize session
  const session = newSession(sessionId, pageId);
  runtime.sessionManager.setSession(session);

  // Create test message
  const closeSession = create(CloseSessionSchema, {
    sessionId,
  });

  // Register handler
  mockWS.registerHandler((message) => {
    if (message.type.case === 'closeSession') {
      return runtime.handleCloseSession(message.type.value);
    }
  });

  // Send message
  mockWS.enqueue(uuidv4(), closeSession);

  // Verify that session was closed
  expect(runtime.sessionManager.getSession(sessionId)).toBeUndefined();
});
