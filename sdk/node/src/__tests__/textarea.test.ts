import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { TextAreaState } from '../session/state/textarea';
import {
  convertTextAreaProtoToState,
  convertStateToTextAreaProto,
  textArea,
} from '../uibuilder/widgets/textarea';
import { createSessionManager, newSession } from '../session';
import { UIBuilder } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';
import { MockClient } from '../websocket/mock/websocket';
test('convertStateToTextAreaProto', () => {
  const id = uuidv4();
  const label = 'Test TextArea';
  const value = 'test value';
  const placeholder = 'Enter text';
  const defaultValue = 'default';
  const required = true;
  const disabled = false;
  const maxLength = 1000;
  const minLength = 10;
  const maxLines = 10;
  const minLines = 3;
  const autoResize = true;

  const state = new TextAreaState(
    id,
    value,
    label,
    placeholder,
    defaultValue,
    required,
    disabled,
    maxLength,
    minLength,
    maxLines,
    minLines,
    autoResize,
  );
  const proto = convertStateToTextAreaProto(state);

  expect(proto.value).toBe(value);
  expect(proto.label).toBe(label);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toBe(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
  expect(proto.maxLength).toBe(maxLength);
  expect(proto.minLength).toBe(minLength);
  expect(proto.maxLines).toBe(maxLines);
  expect(proto.minLines).toBe(minLines);
  expect(proto.autoResize).toBe(autoResize);
});

test('convertTextAreaProtoToState', () => {
  const id = uuidv4();
  const label = 'Test Text Area';
  const value = 'test value';
  const placeholder = 'Enter text';
  const defaultValue = 'default';
  const required = true;
  const disabled = false;
  const maxLength = 1000;
  const minLength = 10;
  const maxLines = 10;
  const minLines = 3;
  const autoResize = true;

  const tempState = new TextAreaState(
    id,
    value,
    label,
    placeholder,
    defaultValue,
    required,
    disabled,
    maxLength,
    minLength,
    maxLines,
    minLines,
    autoResize,
  );
  const proto = convertStateToTextAreaProto(tempState);

  const state = convertTextAreaProtoToState(id, proto);

  if (!state) {
    throw new Error('TextAreaState not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value).toBe(value);
  expect(state.placeholder).toBe(placeholder);
  expect(state.defaultValue).toBe(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
  expect(state.maxLength).toBe(maxLength);
  expect(state.minLength).toBe(minLength);
  expect(state.maxLines).toBe(maxLines);
  expect(state.minLines).toBe(minLines);
  expect(state.autoResize).toBe(autoResize);
});

test('textArea', () => {
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

  const label = 'Test TextArea';
  const options = {
    placeholder: 'Enter text',
    defaultValue: 'default value',
    required: true,
    disabled: true,
    maxLength: 1000,
    minLength: 10,
    maxLines: 10,
    minLines: 3,
    autoResize: false,
  };

  textArea(builder, label, options);

  const widgetId = builder.generatePageID('textArea', [0]);
  const state = session.state.getTextArea(widgetId);

  if (!state) {
    throw new Error('TextAreaState not found');
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
  expect(state.maxLines).toBe(options.maxLines);
  expect(state.minLines).toBe(options.minLines);
  expect(state.autoResize).toBe(options.autoResize);
});
