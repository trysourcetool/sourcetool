import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { TextInputState } from '../internal/session/state/textinput';
import {
  convertTextInputProtoToState,
  convertStateToTextInputProto,
} from '../textinput';
import { createSessionManager, newSession } from '../internal/session';
import { UIBuilder } from '../uibuilder';
import { Page, PageManager } from '../internal/page';
import { Runtime } from '../runtime';
import { MockClient } from '../internal/websocket/mock/websocket';

test('convertStateToTextInputProto', () => {
  const id = uuidv4();
  const label = 'Test TextInput';
  const value = 'test value';
  const placeholder = 'Enter text';
  const defaultValue = 'default';
  const required = true;
  const disabled = false;
  const maxLength = 100;
  const minLength = 10;

  const state = new TextInputState(
    id,
    value,
    label,
    placeholder,
    defaultValue,
    required,
    disabled,
    maxLength,
    minLength,
  );
  const proto = convertStateToTextInputProto(state);

  expect(proto.value).toBe(value);
  expect(proto.label).toBe(label);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toBe(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
  expect(proto.maxLength).toBe(maxLength);
  expect(proto.minLength).toBe(minLength);
});

test('convertTextInputProtoToState', () => {
  const id = uuidv4();
  const label = 'Test TextInput';
  const value = 'test value';
  const placeholder = 'Enter text';
  const defaultValue = 'default';
  const required = true;
  const disabled = false;
  const maxLength = 100;
  const minLength = 10;

  const proto = convertStateToTextInputProto(
    new TextInputState(
      id,
      value,
      label,
      placeholder,
      defaultValue,
      required,
      disabled,
      maxLength,
      minLength,
    ),
  );
  const state = convertTextInputProtoToState(id, proto);

  if (!state) {
    throw new Error('TextInputState not found');
  }

  expect(state.id).toBe(id);
  expect(state.value).toBe(value);
  expect(state.label).toBe(label);
  expect(state.placeholder).toBe(placeholder);
  expect(state.defaultValue).toBe(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
  expect(state.maxLength).toBe(maxLength);
  expect(state.minLength).toBe(minLength);
});

test('textInput', () => {
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

  const label = 'Test TextInput';
  const options = {
    placeholder: 'Enter text',
    defaultValue: 'default value',
    required: true,
    disabled: true,
    maxLength: 100,
    minLength: 10,
  };

  builder.textInput(label, options);

  const widgetId = builder.generatePageID('textInput', [0]);
  const state = session.state.getTextInput(widgetId);

  if (!state) {
    throw new Error('TextInput not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value).toBe(options.defaultValue);
  expect(state.placeholder).toBe(options.placeholder);
  expect(state.defaultValue).toBe(options.defaultValue);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
  expect(state.maxLength).toBe(options.maxLength);
  expect(state.minLength).toBe(options.minLength);
});
