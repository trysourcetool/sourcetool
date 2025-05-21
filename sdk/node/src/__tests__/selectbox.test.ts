import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { SelectboxState } from '../session/state/selectbox';
import {
  convertSelectboxProtoToState,
  convertStateToSelectboxProto,
  selectbox,
} from '../uibuilder/widgets/selectbox';
import { createSessionManager, newSession } from '../session';
import { MockClient } from '../websocket/mock/websocket';
import { Cursor, generateWidgetId } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';

test('convertStateToSelectboxProto', () => {
  const id = uuidv4();
  const label = 'Test Selectbox';
  const value = 1;
  const options = ['Option 1', 'Option 2'];
  const placeholder = 'Select an option';
  const defaultValue = 0;
  const required = true;
  const disabled = false;

  const state = new SelectboxState(
    id,
    value,
    label,
    options,
    placeholder,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToSelectboxProto(state);

  expect(proto.label).toBe(label);
  expect(proto.value).toBe(value);
  expect(proto.options).toEqual(options);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toBe(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
});

test('convertSelectboxProtoToState', () => {
  const id = uuidv4();
  const label = 'Test Selectbox';
  const value = 1;
  const options = ['Option 1', 'Option 2'];
  const placeholder = 'Select an option';
  const defaultValue = 0;
  const required = true;
  const disabled = false;

  const tempState = new SelectboxState(
    id,
    value,
    label,
    options,
    placeholder,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToSelectboxProto(tempState);

  const state = convertSelectboxProtoToState(id, proto);

  if (!state) {
    throw new Error('SelectboxState not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value).toBe(value);
  expect(state.options).toEqual(options);
  expect(state.placeholder).toBe(placeholder);
  expect(state.defaultValue).toBe(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
});

test('selectbox', () => {
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

  const label = 'Test Selectbox';
  const initialOptions = ['Option 1', 'Option 2'];
  const options = {
    options: initialOptions,
    defaultValue: 'Option 1',
    placeholder: 'Select an option',
    required: true,
    disabled: false,
    formatFunc: (val: string) => val.toLowerCase(),
  };

  const cursor = new Cursor();

  selectbox({ runtime, session, page, cursor }, label, options);

  const widgetId = generateWidgetId(page.id, 'selectbox', [0]);
  const state = session.state.getSelectbox(widgetId);

  const expectedDefaultIndex = initialOptions.indexOf(options.defaultValue);
  const expectedFormattedOptions = initialOptions.map(options.formatFunc);

  if (!state) {
    throw new Error('SelectboxState not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value).toBe(expectedDefaultIndex);
  expect(state.options).toEqual(expectedFormattedOptions);
  expect(state.placeholder).toBe(options.placeholder);
  expect(state.defaultValue).toBe(expectedDefaultIndex);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
});
