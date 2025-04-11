import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { DateInputState } from '../internal/session/state/dateinput';
import {
  convertDateInputProtoToState,
  convertStateToDateInputProto,
  dateInput,
} from '../dateinput';
import { createSessionManager, newSession } from '../internal/session';
import { UIBuilder } from '../uibuilder';
import { Page, PageManager } from '../internal/page';
import { Runtime } from '../runtime';
import { MockClient } from '../internal/websocket/mock/websocket';

// Helper function to format date to YYYY-MM-DD
const formatDate = (date: Date | null): string | null => {
  if (!date) {
    return null;
  }
  // Use UTC methods to avoid timezone shifts when formatting
  const year = date.getUTCFullYear();
  const month = (date.getUTCMonth() + 1).toString().padStart(2, '0');
  const day = date.getUTCDate().toString().padStart(2, '0');
  return `${year}-${month}-${day}`;
};

test('convertStateToDateInputProto', () => {
  const now = new Date(Date.UTC(2024, 0, 0));

  const id = uuidv4();
  const label = 'Test DateInput';
  const value = now;
  const placeholder = 'Select date';
  const defaultValue = now;
  const required = true;
  const disabled = false;
  const format = 'YYYY/MM/DD';
  const maxValue = new Date(Date.UTC(now.getFullYear() + 1, 0, 0));
  const minValue = new Date(Date.UTC(now.getFullYear() - 1, 0, 0));
  const location = 'local';

  const state = new DateInputState(
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
  const proto = convertStateToDateInputProto(state);

  expect(proto.value).toBe(formatDate(value));
  expect(proto.label).toBe(label);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toBe(formatDate(defaultValue));
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
  expect(proto.format).toBe(format);
  expect(proto.maxValue).toBe(formatDate(maxValue));
  expect(proto.minValue).toBe(formatDate(minValue));
});

test('convertDateInputProtoToState', () => {
  const now = new Date(Date.UTC(2024, 0, 0));

  const id = uuidv4();
  const label = 'Test DateInput';
  const valueDate = now;
  const defaultValueDate = now;
  const maxValueDate = new Date(Date.UTC(now.getFullYear() + 1, 0, 0));
  const minValueDate = new Date(Date.UTC(now.getFullYear() - 1, 0, 0));
  const placeholder = 'Select date';
  const required = true;
  const disabled = false;
  const format = 'YYYY/MM/DD';
  const location = 'local';

  const tempState = new DateInputState(
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
  const proto = convertStateToDateInputProto(tempState);

  const state = convertDateInputProtoToState(id, proto, location);

  if (!state) {
    throw new Error('DateInput state not found');
  }

  expect(state.id).toBe(id);
  expect(state.label).toBe(label);
  expect(state.value?.getTime()).toBe(valueDate.getTime());
  expect(state.placeholder).toBe(placeholder);
  expect(state.defaultValue?.getTime()).toBe(defaultValueDate.getTime());
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
  expect(state.format).toBe(format);
  expect(state.maxValue?.getTime()).toBe(maxValueDate.getTime());
  expect(state.minValue?.getTime()).toBe(minValueDate.getTime());
  expect(state.location).toBe(location);
});

test('dateInput', () => {
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

  const label = 'Test DateInput';
  const now = new Date(Date.UTC(2024, 0, 0));
  const options = {
    placeholder: 'Select date',
    defaultValue: now,
    required: true,
    disabled: false,
    maxValue: new Date(Date.UTC(now.getFullYear() + 1, 0, 0)),
    minValue: new Date(Date.UTC(now.getFullYear() - 1, 0, 0)),
    location: 'UTC',
    format: 'YYYY-MM-DD',
  };

  dateInput(builder, label, options);

  const widgetId = builder.generatePageID('dateInput', [0]);
  const state = session.state.getDateInput(widgetId);

  if (!state) {
    throw new Error('DateInput not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value?.getTime()).toBe(options.defaultValue.getTime());
  expect(state.placeholder).toBe(options.placeholder);
  expect(state.defaultValue?.getTime()).toBe(options.defaultValue.getTime());
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
  expect(state.format).toBe(options.format);
  expect(state.maxValue?.getTime()).toBe(options.maxValue.getTime());
  expect(state.minValue?.getTime()).toBe(options.minValue.getTime());
  expect(state.location).toBe(options.location);
});
