import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { DateTimeInputState } from '../session/state/datetimeinput';
import {
  convertDateTimeInputProtoToState,
  convertStateToDateTimeInputProto,
  dateTimeInput,
} from '../uibuilder/widgets/datetimeinput';
import { createSessionManager, newSession } from '../session';
import { Cursor, uiBuilderGeneratePageId } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';
import { MockClient } from '../websocket/mock/websocket';

// Helper function to format date to ISO string
const formatDateTime = (date: Date | null): string | null => {
  if (!date) {
    return null;
  }
  return date.toISOString(); // ISO format: YYYY-MM-DDTHH:mm:ss.sssZ
};

test('convertStateToDateTimeInputProto', () => {
  const now = new Date(Date.UTC(2024, 0, 0));

  const id = uuidv4();
  const label = 'Test DateTimeInput';
  const value = now;
  const placeholder = 'Select date and time';
  const defaultValue = now;
  const required = true;
  const disabled = false;
  const format = 'YYYY/MM/DD HH:MM:SS';
  const maxValue = new Date(Date.UTC(now.getFullYear() + 1, 0, 0, 0, 0, 0));
  const minValue = new Date(Date.UTC(now.getFullYear() - 1, 0, 0, 0, 0, 0));
  const location = 'local';

  const state = new DateTimeInputState(
    id,
    value,
    label,
    placeholder,
    defaultValue,
    required,
    disabled,
    format,
    maxValue,
    minValue,
    location,
  );
  const proto = convertStateToDateTimeInputProto(state);

  expect(proto.value).toBe(formatDateTime(value));
  expect(proto.label).toBe(label);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toBe(formatDateTime(defaultValue));
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
  expect(proto.format).toBe(format);
  expect(proto.maxValue).toBe(formatDateTime(maxValue));
  expect(proto.minValue).toBe(formatDateTime(minValue));
});

test('convertDateTimeInputProtoToState', () => {
  const now = new Date(Date.UTC(2024, 0, 0));

  const id = uuidv4();
  const label = 'Test DateTimeInput';
  const valueDate = now;
  const defaultValueDate = now;
  const maxValueDate = new Date(Date.UTC(now.getFullYear() + 1, 0, 0, 0, 0, 0));
  const minValueDate = new Date(Date.UTC(now.getFullYear() - 1, 0, 0, 0, 0, 0));
  const placeholder = 'Select date and time';
  const required = true;
  const disabled = false;
  const format = 'YYYY/MM/DD HH:MM:SS';
  const location = 'local';

  const tempState = new DateTimeInputState(
    id,
    valueDate,
    label,
    placeholder,
    defaultValueDate,
    required,
    disabled,
    format,
    maxValueDate,
    minValueDate,
    location,
  );
  const proto = convertStateToDateTimeInputProto(tempState);

  const state = convertDateTimeInputProtoToState(id, proto, location);

  if (!state) {
    throw new Error('DateTimeInput not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value?.toISOString()).toBe(valueDate.toISOString());
  expect(state.placeholder).toBe(placeholder);
  expect(state.defaultValue?.toISOString()).toBe(
    defaultValueDate.toISOString(),
  );
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
  expect(state.format).toBe(format);
  expect(state.maxValue?.toISOString()).toBe(maxValueDate.toISOString());
  expect(state.minValue?.toISOString()).toBe(minValueDate.toISOString());
  expect(state.location).toBe(location);
});

test('dateTimeInput', () => {
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

  const label = 'Test DateTimeInput';
  const now = new Date(Date.UTC(2024, 0, 0, 0, 0, 0));
  const options = {
    placeholder: 'Select date and time',
    defaultValue: now,
    required: true,
    disabled: false,
    format: 'YYYY-MM-DD HH:mm:ss',
    maxValue: new Date(Date.UTC(now.getFullYear() + 1, 0, 0, 0, 0, 0)),
    minValue: new Date(Date.UTC(now.getFullYear() - 1, 0, 0, 0, 0, 0)),
    location: 'UTC',
  };

  const cursor = new Cursor();

  dateTimeInput({ runtime, session, page, cursor }, label, options);

  const widgetId = uiBuilderGeneratePageId(page.id, 'datetimeInput', [0]);
  const state = session.state.getDateTimeInput(widgetId);

  if (!state) {
    throw new Error('DateTimeInput not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value?.toISOString()).toBe(options.defaultValue.toISOString());
  expect(state.placeholder).toBe(options.placeholder);
  expect(state.defaultValue?.toISOString()).toBe(
    options.defaultValue.toISOString(),
  );
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
  expect(state.format).toBe(options.format);
  expect(state.maxValue?.toISOString()).toBe(options.maxValue.toISOString());
  expect(state.minValue?.toISOString()).toBe(options.minValue.toISOString());
  expect(state.location).toBe(options.location);
});
