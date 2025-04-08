// src/checkbox.spec.ts
import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { CheckboxState } from './internal/session/state/checkbox';
import {
  convertCheckboxProtoToState,
  convertStateToCheckboxProto,
} from './checkbox';
import { createSessionManager, newSession } from './internal/session';
import { UIBuilder } from './uibuilder';
import { Page, PageManager } from './internal/page';
import { Runtime } from './runtime';
import { MockClient } from './internal/websocket/mock/websocket';

test('convertStateToCheckboxProto', () => {
  const id = uuidv4();
  const label = 'Test Checkbox';
  const value = true;
  const defaultValue = false;
  const required = true;
  const disabled = false;

  const state = new CheckboxState(
    id,
    label,
    value,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToCheckboxProto(state);

  expect(proto.label).toBe(label);
  expect(proto.value).toBe(value);
  expect(proto.defaultValue).toBe(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
});

test('convertCheckboxProtoToState', () => {
  const id = uuidv4();
  const label = 'Test Checkbox';
  const value = true;
  const defaultValue = false;
  const required = true;
  const disabled = false;

  const tempState = new CheckboxState(
    id,
    label,
    value,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToCheckboxProto(tempState);

  const state = convertCheckboxProtoToState(id, proto);

  if (!state) {
    throw new Error('Checkbox state not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value).toBe(value);
  expect(state.defaultValue).toBe(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
});

test('checkbox', () => {
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

  const label = 'Test Checkbox';
  const options = {
    defaultValue: false,
    required: true,
    disabled: false,
  };

  builder.checkbox(label, options);

  const widgetId = builder.generatePageID('checkbox', [0]);
  const state = session.state.getCheckbox(widgetId);

  if (!state) {
    throw new Error('Checkbox state not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value).toBe(options.defaultValue);
  expect(state.defaultValue).toBe(options.defaultValue);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
});
