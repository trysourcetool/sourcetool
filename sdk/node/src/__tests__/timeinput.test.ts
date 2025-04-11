import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { TimeInputState } from '../session/state/timeinput';
import {
  convertTimeInputProtoToState,
  convertStateToTimeInputProto,
  timeInput,
} from '../uibuilder/widgets/timeinput';
import { createSessionManager, newSession } from '../session';
import { MockClient } from '../websocket/mock/websocket';
import { UIBuilder } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';

// Helper function to format Date object to HH:MM:SS string
const formatTime = (date: Date | null | undefined): string | null => {
  if (!date) {
    return null;
  }
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');
  return `${hours}:${minutes}:${seconds}`;
};

test('convertStateToTimeInputProto', () => {
  const now = new Date();

  const id = uuidv4();
  const label = 'Test TimeInput';
  const value = now;
  const placeholder = 'Select time';
  const defaultValue = now;
  const required = true;
  const disabled = false;
  const location = 'local';

  const state = new TimeInputState(
    id,
    value,
    label,
    placeholder,
    defaultValue,
    required,
    disabled,
    location,
  );
  const proto = convertStateToTimeInputProto(state);

  expect(proto.value).toBe(formatTime(value));
  expect(proto.label).toBe(label);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toBe(formatTime(defaultValue));
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
});

test('convertTimeInputProtoToState', () => {
  const now = new Date();

  const id = uuidv4();
  const label = 'Test TimeInput';
  const valueTime = now;
  const defaultValueTime = now;
  const placeholder = 'Select time';
  const required = true;
  const disabled = false;
  const location = 'local';

  const tempState = new TimeInputState(
    id,
    valueTime,
    label,
    placeholder,
    defaultValueTime,
    required,
    disabled,
    location,
  );
  const proto = convertStateToTimeInputProto(tempState);

  const state = convertTimeInputProtoToState(id, proto, location);

  if (!state) {
    throw new Error('TimeInputState not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(formatTime(state.value)).toBe(formatTime(valueTime));
  expect(state.placeholder).toBe(placeholder);
  expect(formatTime(state.defaultValue)).toBe(formatTime(defaultValueTime));
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
  expect(state.location).toBe(location);
});

test('timeInput', () => {
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

  const now = new Date();
  const label = 'Test TimeInput';
  const options = {
    placeholder: 'Select time',
    defaultValue: now,
    required: true,
    disabled: true,
    location: 'UTC',
  };

  timeInput(builder, label, options);

  const widgetId = builder.generatePageID('timeInput', [0]);
  const state = session.state.getTimeInput(widgetId);

  if (!state) {
    throw new Error('TimeInputState not found');
  }

  const formattedValue = formatTime(state.value);
  const formattedDefault = formatTime(state.defaultValue);
  const expectedFormattedDefault = formatTime(options.defaultValue);

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(formattedValue).toBe(expectedFormattedDefault);
  expect(state.placeholder).toBe(options.placeholder);
  expect(formattedDefault).toBe(expectedFormattedDefault);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
  expect(state.location).toBe(options.location);
});
