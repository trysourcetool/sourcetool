import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { MultiSelectState } from '../session/state/multiselect';
import {
  convertMultiSelectProtoToState,
  convertStateToMultiSelectProto,
  multiSelect,
} from '../uibuilder/widgets/multiselect';
import { createSessionManager, newSession } from '../session';
import { MockClient } from '../websocket/mock/websocket';
import { Cursor, uiBuilderGeneratePageID } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';

test('convertStateToMultiSelectProto', () => {
  const id = uuidv4();
  const label = 'Test MultiSelect';
  const value = [0, 2];
  const options = ['Option 1', 'Option 2', 'Option 3'];
  const placeholder = 'Select options';
  const defaultValue = [0];
  const required = true;
  const disabled = false;

  const state = new MultiSelectState(
    id,
    value,
    label,
    options,
    placeholder,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToMultiSelectProto(state);

  expect(proto.label).toBe(label);
  expect(proto.value).toEqual(value);
  expect(proto.options).toEqual(options);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toEqual(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
});

test('convertMultiSelectProtoToState', () => {
  const id = uuidv4();
  const label = 'Test MultiSelect';
  const value = [0, 2];
  const options = ['Option 1', 'Option 2', 'Option 3'];
  const placeholder = 'Select options';
  const defaultValue = [0];
  const required = true;
  const disabled = false;

  const tempState = new MultiSelectState(
    id,
    value,
    label,
    options,
    placeholder,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToMultiSelectProto(tempState);

  const state = convertMultiSelectProtoToState(id, proto);

  if (!state) {
    throw new Error('MultiSelectState not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value).toEqual(value);
  expect(state.options).toEqual(options);
  expect(state.placeholder).toBe(placeholder);
  expect(state.defaultValue).toEqual(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
});

test('multiSelect', () => {
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

  const label = 'Test MultiSelect';
  const initialOptions = ['Option 1', 'Option 2', 'Option 3'];
  const options = {
    options: initialOptions,
    defaultValue: ['Option 1', 'Option 3'],
    placeholder: 'Select options',
    required: true,
    disabled: false,
    formatFunc: (val: string) => `Option: ${val}`,
  };

  const cursor = new Cursor();

  multiSelect({ runtime, session, page, cursor }, label, options);

  const widgetId = uiBuilderGeneratePageID(page.id, 'multiselect', [0]);
  const state = session.state.getMultiSelect(widgetId);

  const expectedDefaultIndexes = options.defaultValue
    .map((dv) => initialOptions.indexOf(dv))
    .filter((index) => index !== -1);
  const expectedFormattedOptions = initialOptions.map(options.formatFunc);

  if (!state) {
    throw new Error('MultiSelectState not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value).toEqual(expectedDefaultIndexes);
  expect(state.options).toEqual(expectedFormattedOptions);
  expect(state.placeholder).toBe(options.placeholder);
  expect(state.defaultValue).toEqual(expectedDefaultIndexes);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
});
