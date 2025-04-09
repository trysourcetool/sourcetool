import { expect, test } from 'vitest';
import {
  createSessionManager,
  createSessionState,
  newSession,
} from './session';
import { v4 as uuidv4 } from 'uuid';
import { RadioState } from './session/state/radio';
import { ButtonState } from './session/state/button';
import { FormState } from './session/state/form';

test('newSession', () => {
  const id = uuidv4();
  const pageID = uuidv4();

  const session = newSession(id, pageID);

  expect(session.id).toEqual(id);
  expect(session.pageId).toEqual(pageID);
});

test('session manager get set delete', () => {
  const manager = createSessionManager();
  const id = uuidv4();
  const pageID = uuidv4();
  const session = newSession(id, pageID);

  manager.setSession(session);

  const got = manager.getSession(id);
  expect(got).toEqual(session);

  manager.disconnectSession(id);
  const got2 = manager.getSession(id);
  expect(got2).toBeUndefined();
});

test('session manager concurrent access', () => {
  const manager = createSessionManager();

  for (let i = 0; i < 10; i++) {
    const id = uuidv4();
    const pageID = uuidv4();
    const session = newSession(id, pageID);

    manager.setSession(session);
    const got = manager.getSession(id);
    expect(got).toEqual(session);

    manager.disconnectSession(id);
    const got2 = manager.getSession(id);
    expect(got2).toBeUndefined();
  }
});

test('state set and get', () => {
  const state = createSessionState();
  const id = uuidv4();

  const radioState = new RadioState(id, null, 'Test Radio', [
    'Option 1',
    'Option 2',
  ]);

  state.set(id, radioState);

  const got = state.get(id);
  expect(got).toEqual(radioState);
  expect(got?.getType()).toEqual('radio');
});

test('state reset state', () => {
  const state = createSessionState();
  const id = uuidv4();
  const radioState = new RadioState(id, null, 'Test Radio', [
    'Option 1',
    'Option 2',
  ]);

  state.set(id, radioState);
  state.resetStates();
  const got = state.get(id);
  expect(got).toBeUndefined();
});

test('state reset buttons', () => {
  const state = createSessionState();
  const buttonId = uuidv4();
  const formId = uuidv4();

  const buttonState = new ButtonState(buttonId, true);
  const formState = new FormState(formId, true);

  state.set(buttonId, buttonState);
  state.set(formId, formState);

  state.resetButtons();

  const gotButton = state.getButton(buttonId);
  expect(gotButton?.value).toBeFalsy();

  const gotForm = state.getForm(formId);
  expect(gotForm?.value).toBeFalsy();
});

test('state set states', () => {
  const state = createSessionState();
  const id1 = uuidv4();
  const id2 = uuidv4();

  const states = {
    [id1]: new RadioState(id1, null, 'Radio 1', ['Option 1', 'Option 2']),
    [id2]: new RadioState(id2, null, 'Radio 2', ['Option 3', 'Option 4']),
  };

  state.setStates(new Map(Object.entries(states)));

  const got1 = state.get(id1);
  expect(got1).toEqual(states[id1]);

  const got2 = state.get(id2);
  expect(got2).toEqual(states[id2]);
});

test('state concurrent access', () => {
  const state = createSessionState();

  for (let i = 0; i < 10; i++) {
    const id = uuidv4();
    const radioState = new RadioState(id, null, 'Test Radio ' + i, [
      'Option 1',
      'Option 2',
    ]);
    state.set(id, radioState);

    const got = state.get(id);
    expect(got).toEqual(radioState);
  }
});
