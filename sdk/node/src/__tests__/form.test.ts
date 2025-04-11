import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import { FormState } from '../session/state/form';
import {
  convertFormProtoToState,
  convertStateToFormProto,
} from '../uibuilder/widgets/form';
import { createSessionManager, newSession } from '../session';
import { UIBuilder } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';
import { MockClient } from '../websocket/mock/websocket';

test('convertStateToFormProto', () => {
  const id = uuidv4();
  const value = true;
  const buttonLabel = 'Submit';
  const buttonDisabled = true;
  const clearOnSubmit = true;

  const state = new FormState(
    id,
    value,
    buttonLabel,
    buttonDisabled,
    clearOnSubmit,
  );
  const proto = convertStateToFormProto(state);

  expect(proto.value).toBe(value);
  expect(proto.buttonLabel).toBe(buttonLabel);
  expect(proto.buttonDisabled).toBe(buttonDisabled);
  expect(proto.clearOnSubmit).toBe(clearOnSubmit);
});

test('convertFormProtoToState', () => {
  const id = uuidv4();
  const value = true;
  const buttonLabel = 'Submit';
  const buttonDisabled = true;
  const clearOnSubmit = true;

  const tempState = new FormState(
    id,
    value,
    buttonLabel,
    buttonDisabled,
    clearOnSubmit,
  );
  const proto = convertStateToFormProto(tempState);

  const state = convertFormProtoToState(id, proto);

  if (!state) {
    throw new Error('FormState not found');
  }

  expect(state.id).toBe(id);
  expect(state.value).toBe(value);
  expect(state.buttonLabel).toBe(buttonLabel);
  expect(state.buttonDisabled).toBe(buttonDisabled);
  expect(state.clearOnSubmit).toBe(clearOnSubmit);
});

test('form', () => {
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

  const buttonLabel = 'Submit';
  const formOptions = {
    buttonDisabled: false,
    clearOnSubmit: false,
  };

  const [formBuilder, submitted] = builder.form(buttonLabel, formOptions);

  expect(submitted).toBe(false);

  const formWidgetId = builder.generatePageID('form', [0]);
  let formState = session.state.getForm(formWidgetId);

  if (!formState) {
    throw new Error('FormState not found');
  }

  expect(formState.id).toBe(formWidgetId);
  expect(formState.buttonLabel).toBe(buttonLabel);
  expect(formState.buttonDisabled).toBe(formOptions.buttonDisabled);
  expect(formState.clearOnSubmit).toBe(formOptions.clearOnSubmit);
  expect(formState.value).toBe(false);

  formBuilder.textInput('Name', { defaultValue: 'John Doe' });

  const textInputId = builder.generatePageID('textInput', [0, 0]);
  const textInputState = session.state.getTextInput(textInputId);

  if (!textInputState) {
    throw new Error('TextInputState not found');
  }

  expect(textInputState.id).toBe(textInputId);
});
