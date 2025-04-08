import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { RadioState } from './internal/session/state/radio';
import { convertRadioProtoToState, convertStateToRadioProto } from './radio';
import { createSessionManager, newSession } from './internal/session';
import { MockClient } from './internal/websocket/mock/websocket';
import { UIBuilder } from './uibuilder';
import { Page, PageManager } from './internal/page';
import { Runtime } from './runtime';

test('convertStateToRadioProto', () => {
  const id = uuidv4();
  const label = 'Test Radio';
  const value = 1;
  const options = ['Option 1', 'Option 2'];
  const defaultValue = 0;
  const required = true;
  const disabled = false;

  const state = new RadioState(
    id,
    value,
    label,
    options,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToRadioProto(state);

  expect(proto.label).toBe(label);
  expect(proto.value).toBe(value);
  expect(proto.options).toEqual(options);
  expect(proto.defaultValue).toBe(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
});

test('convertRadioProtoToState', () => {
  const id = uuidv4();
  const label = 'Test Radio';
  const value = 1;
  const options = ['Option 1', 'Option 2'];
  const defaultValue = 0;
  const required = true;
  const disabled = false;

  const tempState = new RadioState(
    id,
    value,
    label,
    options,
    defaultValue,
    required,
    disabled,
  );
  const proto = convertStateToRadioProto(tempState);

  const state = convertRadioProtoToState(id, proto);

  if (!state) {
    throw new Error('RadioState not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value).toBe(value);
  expect(state.options).toEqual(options);
  expect(state.defaultValue).toBe(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
});

test('radio', () => {
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

  const label = 'Test Radio';
  const initialOptions = ['Option 1', 'Option 2'];
  const options = {
    options: initialOptions,
    defaultValue: 'Option 1',
    required: true,
    disabled: false,
    formatFunc: (val: string) => val.toUpperCase(),
  };

  builder.radio(label, options);

  const widgetId = builder.generatePageID('radio', [0]);
  const state = session.state.getRadio(widgetId);

  const expectedDefaultIndex = initialOptions.indexOf(options.defaultValue);
  const expectedFormattedOptions = initialOptions.map(options.formatFunc);

  if (!state) {
    throw new Error('RadioState not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value).toBe(expectedDefaultIndex);
  expect(state.options).toEqual(expectedFormattedOptions);
  expect(state.defaultValue).toBe(expectedDefaultIndex);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
});
