import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { MarkdownState } from './internal/session/state/markdown';
import {
  convertMarkdownProtoToState,
  convertStateToMarkdownProto,
  markdown,
} from './markdown'; // convertStateToMarkdownProto は後でエクスポートする必要あり
import { createSessionManager, newSession } from './internal/session';
import { UIBuilder } from './uibuilder';
import { Page, PageManager } from './internal/page';
import { Runtime } from './runtime';
import { MockClient } from './internal/websocket/mock/websocket';

test('convertStateToMarkdownProto', () => {
  const id = uuidv4();
  const body = '# Test Markdown';

  const state = new MarkdownState(id, body);
  const proto = convertStateToMarkdownProto(state);

  expect(proto.body).toBe(body);
});

test('convertMarkdownProtoToState', () => {
  const id = uuidv4();
  const body = '# Test Markdown';

  const tempState = new MarkdownState(id, body);
  const proto = convertStateToMarkdownProto(tempState);

  const state = convertMarkdownProtoToState(id, proto);

  if (!state) {
    throw new Error('MarkdownState not found');
  }

  expect(state.id).toBe(id);
  expect(state.body).toBe(body);
});

test('markdown', () => {
  const sessionId = uuidv4();
  const pageId = uuidv4();
  const session = newSession(sessionId, pageId);

  const pageManager = new PageManager({
    [pageId]: new Page(
      pageId,
      'Test Page',
      '/test',
      [1, 2, 3],
      async () => {},
      ['test'],
    ),
  });

  const sessionManager = createSessionManager();
  const mockWS = new MockClient();
  const runtime = new Runtime(mockWS, sessionManager, pageManager);

  const page = pageManager.getPage(pageId);

  if (!page) {
    throw new Error('Page not found');
  }

  const builder = new UIBuilder(runtime, session, page);

  const bodyContent = '# Test Markdown';

  markdown(builder, bodyContent);

  const widgetId = builder.generatePageID('markdown', [0]);
  const state = session.state.getMarkdown(widgetId);

  if (!state) {
    throw new Error('MarkdownState not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.body).toBe(bodyContent);
});
