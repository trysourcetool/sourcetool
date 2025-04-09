import { expect, test, describe } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { ColumnsState } from './internal/session/state/columns';
import { ColumnItemState } from './internal/session/state/columnitem';
import {
  convertColumnsProtoToState,
  convertStateToColumnsProto,
  convertColumnItemProtoToState,
  convertStateToColumnItemProto,
} from './columns';
import { createSessionManager, newSession } from './internal/session';
import { UIBuilder } from './uibuilder';
import { Page, PageManager } from './internal/page';
import { Runtime } from './runtime';
import { MockClient } from './internal/websocket/mock/websocket';

test('convertStateToColumnsProto', () => {
  const id = uuidv4();
  const cols = 3;
  const state = new ColumnsState(id, cols);
  const proto = convertStateToColumnsProto(state);

  expect(proto.columns).toBe(cols);
});

test('convertColumnsProtoToState', () => {
  const id = uuidv4();
  const cols = 3;
  const tempState = new ColumnsState(id, cols);
  const proto = convertStateToColumnsProto(tempState);
  const state = convertColumnsProtoToState(id, proto);

  if (!state) {
    throw new Error('Columns state not found');
  }

  expect(state.id).toBe(id);
  expect(state.columns).toBe(cols);
});

test('convertStateToColumnItemProto', () => {
  const id = uuidv4();
  const weight = 0.5;
  const state = new ColumnItemState(id, weight);
  const proto = convertStateToColumnItemProto(state);

  expect(proto.weight).toBe(weight);
});

test('convertColumnItemProtoToState', () => {
  const id = uuidv4();
  const weight = 0.5;
  const tempState = new ColumnItemState(id, weight);
  const proto = convertStateToColumnItemProto(tempState);
  const state = convertColumnItemProtoToState(id, proto);

  if (!state) {
    throw new Error('ColumnItem state not found');
  }

  expect(state.id).toBe(id);
  expect(state.weight).toBe(weight);
});

test('columns', () => {
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

  const cols = 3;

  builder.columns(cols);

  const messages = mockWS.getMessages();

  const expectedMessages = cols + 1;

  if (messages.length !== expectedMessages) {
    throw new Error(
      `Websocket messages count = ${messages.length}, want ${expectedMessages}`,
    );
  }

  const widgetId = builder.generatePageID('columns', [0]);
  const state = session.state.get(widgetId) as ColumnsState;

  if (!state) {
    throw new Error('Columns not found');
  }
  expect(state.columns).toBe(cols);

  for (let i = 0; i < cols; i++) {
    const columnPath = [0, i];
    const columnId = builder.generatePageID('columnItem', columnPath);

    const columnState = session.state.get(columnId) as ColumnItemState;

    if (!columnState) {
      throw new Error('Column not found');
    }

    const expectedWeight = 1 / cols;

    expect(columnState.weight).toBeCloseTo(expectedWeight);
  }
});

test('columns with weight', () => {
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

  const cols = 3;
  const options = { weight: [2, 1, 1] };
  const totalWeight = 4;

  builder.columns(cols, options);

  for (let i = 0; i < cols; i++) {
    const columnPath = [0, i];
    const columnId = builder.generatePageID('columnItem', columnPath);

    const columnState = session.state.get(columnId) as ColumnItemState;

    if (!columnState) {
      throw new Error('Column not found');
    }

    const expectedWeight = options.weight[i] / totalWeight;

    expect(columnState.weight).toBeCloseTo(expectedWeight);
  }
});

describe('columns invalid input', () => {
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

  const testInputs = [
    ['Zero columns', 0, undefined],
    ['Negative columns', -1, undefined],
    ['Invalid weights length', 3, { weight: [1, 1] }],
    ['Zero weights', 3, { weight: [0, 0, 0] }],
    ['Negative weights', 3, { weight: [-1, 1, 1] }],
  ] as [string, number, { weight: number[] } | undefined][];

  for (let i = 0; i < testInputs.length; i++) {
    const [name, cols, options] = testInputs[i];

    test(name, () => {
      builder.columns(cols, options);
      const columns = session.state.get(builder.generatePageID('columns', [i]));

      expect(columns).toBeUndefined();
    });
  }
});
