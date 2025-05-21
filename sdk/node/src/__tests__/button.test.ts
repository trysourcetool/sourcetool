import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { ButtonState } from '../session/state/button';
import {
  convertButtonProtoToState,
  convertStateToButtonProto,
} from '../uibuilder/widgets/button';
import { createSessionManager, newSession } from '../session';
import { UIBuilder, uiBuilderGeneratePageID } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';
import { MockClient } from '../websocket/mock/websocket';

test('convertStateToButtonProto', () => {
  const id = uuidv4();
  const label = 'Test Button';
  const value = true;
  const disabled = false;
  const state = new ButtonState(id, value, label, disabled);
  const proto = convertStateToButtonProto(state);

  expect(proto.value).toBe(value);
  expect(proto.label).toBe(label);
  expect(proto.disabled).toBe(disabled);
});

test('convertButtonProtoToState', () => {
  const id = uuidv4();
  const label = 'Test Button';
  const value = true;
  const disabled = false;
  const proto = convertStateToButtonProto(
    new ButtonState(id, value, label, disabled),
  );
  const state = convertButtonProtoToState(id, proto);

  if (!state) {
    throw new Error('Button state not found');
  }

  expect(state.id).toBe(id);
  expect(state.value).toBe(value);
  expect(state.label).toBe(label);
  expect(state.disabled).toBe(disabled);
});

test('button', () => {
  const sessionId = uuidv4();
  const pageId = uuidv4();
  const session = newSession(sessionId, pageId);
  const page = new Page(
    pageId,
    'Test Page',
    '/test',
    [1, 2, 3],
    async () => {},
    ['test'],
  );
  const pageManager = new PageManager({
    [pageId]: page,
  });

  const sessionManager = createSessionManager();

  const mockWS = new MockClient();

  const runtime = new Runtime(mockWS, sessionManager, pageManager);

  if (!page) {
    throw new Error('Page not found');
  }

  const builder = new UIBuilder(runtime, session, page);

  builder.button('Test Button');

  const widgetId = uiBuilderGeneratePageID(page.id, 'button', [0]);

  const state = session.state.getButton(widgetId);

  if (!state) {
    throw new Error('Button state not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.value).toBe(false);
  expect(state.label).toBe('Test Button');
  expect(state.disabled).toBe(false);
});
