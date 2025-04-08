import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { CheckboxGroupState } from './internal/session/state/checkboxgroup';
import {
  convertCheckboxGroupProtoToState,
  convertStateToCheckboxGroupProto,
} from './checkboxgroup';
import { createSessionManager, newSession } from './internal/session';
import { UIBuilder } from './uibuilder';
import { Page, PageManager } from './internal/page';
import { Runtime } from './runtime';
import { MockClient } from './internal/websocket/mock/websocket';

test('convertStateToCheckboxGroupProto', () => {
  const id = uuidv4();
  const label = 'Test CheckboxGroup';
  const value = [0, 2];
  const options = ['Option 1', 'Option 2', 'Option 3'];
  const defaultValue = [0];
  const required = true;
  const disabled = false;

  const state = new CheckboxGroupState(
    id,
    value,
    label,
    options,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToCheckboxGroupProto(state);

  expect(proto.label).toBe(label);
  expect(proto.value).toEqual(value);
  expect(proto.options).toEqual(options);
  expect(proto.defaultValue).toEqual(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
});

test('convertCheckboxGroupProtoToState', () => {
  const id = uuidv4();
  const label = 'Test CheckboxGroup';
  const value = [0, 2];
  const options = ['Option 1', 'Option 2', 'Option 3'];
  const defaultValue = [0];
  const required = true;
  const disabled = false;

  const tempState = new CheckboxGroupState(
    id,
    value,
    label,
    options,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToCheckboxGroupProto(tempState);

  const state = convertCheckboxGroupProtoToState(id, proto);

  if (!state) {
    throw new Error('CheckboxGroup state not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value).toEqual(value);
  expect(state.options).toEqual(options);
  expect(state.defaultValue).toEqual(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
});

test('checkboxGroup', () => {
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

  const label = 'Test CheckboxGroup';
  const initialOptions = ['Option 1', 'Option 2', 'Option 3'];
  const options = {
    options: initialOptions,
    defaultValue: ['Option 1', 'Option 3'],
    required: true,
    disabled: false,
    formatFunc: (val: string, index: number) =>
      `${index + 1}. ${val.toUpperCase()}`,
  };

  builder.checkboxGroup(label, options);

  const widgetId = builder.generatePageID('checkboxGroup', [0]);
  const state = session.state.getCheckboxGroup(widgetId);

  if (!state) {
    throw new Error('CheckboxGroup not found');
  }

  const expectedDefaultIndexes = options.defaultValue
    .map((dv) => initialOptions.indexOf(dv))
    .filter((index) => index !== -1);

  const expectedFormattedOptions = initialOptions.map(options.formatFunc);

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value).toEqual(expectedDefaultIndexes);
  expect(state.options).toEqual(expectedFormattedOptions);
  expect(state.defaultValue).toEqual(expectedDefaultIndexes);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
});
