import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { NumberInputState } from '../session/state/numberinput';
import {
  convertNumberInputProtoToState,
  convertStateToNumberInputProto,
} from '../uibuilder/widgets/numberinput';
import { createSessionManager, newSession } from '../session';
import { MockClient } from '../websocket/mock/websocket';
import { UIBuilder, uiBuilderGeneratePageID } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';

test('convertStateToNumberInputProto', () => {
  const id = uuidv4();
  const label = 'Test NumberInput';
  const value = 42.5;
  const placeholder = 'Enter number';
  const defaultValue = 0.0;
  const required = true;
  const disabled = false;
  const maxValue = 100.0;
  const minValue = 0.0;

  const state = new NumberInputState(
    id,
    value,
    label,
    placeholder,
    defaultValue,
    required,
    disabled,
    maxValue,
    minValue,
  );
  const proto = convertStateToNumberInputProto(state);

  expect(proto.value).toBe(value);
  expect(proto.label).toBe(label);
  expect(proto.placeholder).toBe(placeholder);
  expect(proto.defaultValue).toBe(defaultValue);
  expect(proto.required).toBe(required);
  expect(proto.disabled).toBe(disabled);
  expect(proto.maxValue).toBe(maxValue);
  expect(proto.minValue).toBe(minValue);
});

test('convertNumberInputProtoToState', () => {
  const id = uuidv4();
  const label = 'Test NumberInput';
  const value = 42.5;
  const placeholder = 'Enter number';
  const defaultValue = 0.0;
  const required = true;
  const disabled = false;
  const maxValue = 100.0;
  const minValue = 0.0;

  const proto = convertStateToNumberInputProto(
    new NumberInputState(
      id,
      value,
      label,
      placeholder,
      defaultValue,
      required,
      disabled,
      maxValue,
      minValue,
    ),
  );
  const state = convertNumberInputProtoToState(id, proto);

  if (!state) {
    throw new Error('NumberInputState not found');
  }

  expect(state.id).toBe(id);
  expect(state.value).toBe(value);
  expect(state.label).toBe(label);
  expect(state.placeholder).toBe(placeholder);
  expect(state.defaultValue).toBe(defaultValue);
  expect(state.required).toBe(required);
  expect(state.disabled).toBe(disabled);
  expect(state.maxValue).toBe(maxValue);
  expect(state.minValue).toBe(minValue);
});

test('numberInput', () => {
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

  const label = 'Test NumberInput';
  const options = {
    placeholder: 'Enter number',
    defaultValue: 42.5,
    required: true,
    disabled: true,
    maxValue: 100.0,
    minValue: 0.0,
  };

  builder.numberInput(label, options);

  const widgetId = uiBuilderGeneratePageID(page.id, 'numberInput', [0]);
  const state = session.state.getNumberInput(widgetId);

  if (!state) {
    throw new Error('NumberInputState not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.label).toBe(label);
  expect(state.value).toBe(options.defaultValue);
  expect(state.placeholder).toBe(options.placeholder);
  expect(state.defaultValue).toBe(options.defaultValue);
  expect(state.required).toBe(options.required);
  expect(state.disabled).toBe(options.disabled);
  expect(state.maxValue).toBe(options.maxValue);
  expect(state.minValue).toBe(options.minValue);
});
