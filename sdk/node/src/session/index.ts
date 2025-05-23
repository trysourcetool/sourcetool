import { ButtonState } from './state/button';
import { CheckboxState } from './state/checkbox';
import { CheckboxGroupState } from './state/checkboxgroup';
import { ColumnItemState } from './state/columnitem';
import { ColumnsState } from './state/columns';
import { DateInputState } from './state/dateinput';
import { DateTimeInputState } from './state/datetimeinput';
import { FormState } from './state/form';
import { MarkdownState } from './state/markdown';
import { MultiSelectState } from './state/multiselect';
import { NumberInputState } from './state/numberinput';
import { RadioState } from './state/radio';
import { SelectboxState } from './state/selectbox';
import { TableState } from './state/table';
import { TextAreaState } from './state/textarea';
import { TextInputState } from './state/textinput';
import { TimeInputState } from './state/timeinput';

/**
 * Constants
 */
const DISCONNECTED_SESSION_TTL = 2 * 60 * 1000; // 2 minutes in milliseconds
const MAX_DISCONNECTED_SESSIONS = 128;

/**
 * Widget state interface
 */
export type WidgetState =
  | ButtonState
  | CheckboxState
  | CheckboxGroupState
  | ColumnsState
  | ColumnItemState
  | DateInputState
  | DateTimeInputState
  | FormState
  | MarkdownState
  | MultiSelectState
  | NumberInputState
  | RadioState
  | SelectboxState
  | TableState
  | TextAreaState
  | TextInputState
  | TimeInputState;

/**
 * Session state interface
 */
export interface SessionState {
  states: Map<string, WidgetState>;
  get(id: string): WidgetState | undefined;
  set(id: string, state: WidgetState): void;
  getButton(id: string): ButtonState | undefined;
  getCheckbox(id: string): CheckboxState | undefined;
  getCheckboxGroup(id: string): CheckboxGroupState | undefined;
  getColumns(id: string): ColumnsState | undefined;
  getDateInput(id: string): DateInputState | undefined;
  getDateTimeInput(id: string): DateTimeInputState | undefined;
  getForm(id: string): FormState | undefined;
  getMarkdown(id: string): MarkdownState | undefined;
  getMultiSelect(id: string): MultiSelectState | undefined;
  getNumberInput(id: string): NumberInputState | undefined;
  getRadio(id: string): RadioState | undefined;
  getSelectbox(id: string): SelectboxState | undefined;
  getTable(id: string): TableState | undefined;
  getTextArea(id: string): TextAreaState | undefined;
  getTextInput(id: string): TextInputState | undefined;
  getTimeInput(id: string): TimeInputState | undefined;
  resetStates(): void;
  setStates(newStates: Map<string, WidgetState>): void;
  resetButtons(): void;
}

/**
 * Session interface
 */
export interface Session {
  /**
   * Session ID
   */
  id: string;

  /**
   * Page ID
   */
  pageId: string;

  /**
   * Session state
   */
  state: SessionState;
}

/**
 * Disconnected session interface
 */
interface DisconnectedSession {
  /**
   * Session
   */
  session: Session;

  /**
   * Disconnected at timestamp
   */
  disconnectedAt: Date;
}

/**
 * Session manager interface
 */
export interface SessionManager {
  /**
   * Set a session
   * @param session Session
   */
  setSession(session: Session): void;

  /**
   * Get a session
   * @param id Session ID
   * @returns Session or undefined if not found
   */
  getSession(id: string): Session | undefined;

  /**
   * Disconnect a session
   * @param id Session ID
   */
  disconnectSession(id: string): void;
}

/**
 * Create a session state
 * @returns Session state
 */
export function createSessionState(): SessionState {
  const statesMap = new Map<string, WidgetState>();

  const sessionState: SessionState = {
    states: statesMap,

    get(id: string): WidgetState | undefined {
      return statesMap.get(id);
    },

    set(id: string, state: WidgetState): void {
      statesMap.set(id, state);
    },

    getButton(id: string): ButtonState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'button' ? state : undefined;
    },

    getCheckbox(id: string): CheckboxState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'checkbox' ? state : undefined;
    },

    getCheckboxGroup(id: string): CheckboxGroupState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'checkboxGroup' ? state : undefined;
    },

    getColumns(id: string): ColumnsState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'columns' ? state : undefined;
    },

    getDateInput(id: string): DateInputState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'dateInput' ? state : undefined;
    },

    getDateTimeInput(id: string): DateTimeInputState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'datetimeInput' ? state : undefined;
    },

    getForm(id: string): FormState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'form' ? state : undefined;
    },

    getMarkdown(id: string): MarkdownState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'markdown' ? state : undefined;
    },

    getMultiSelect(id: string): MultiSelectState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'multiselect' ? state : undefined;
    },

    getNumberInput(id: string): NumberInputState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'numberInput' ? state : undefined;
    },

    getRadio(id: string): RadioState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'radio' ? state : undefined;
    },

    getSelectbox(id: string): SelectboxState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'selectbox' ? state : undefined;
    },

    getTable(id: string): TableState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'table' ? state : undefined;
    },

    getTextArea(id: string): TextAreaState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'textArea' ? state : undefined;
    },

    getTextInput(id: string): TextInputState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'textInput' ? state : undefined;
    },

    getTimeInput(id: string): TimeInputState | undefined {
      const state = statesMap.get(id);
      return state?.type === 'timeInput' ? state : undefined;
    },

    resetStates(): void {
      statesMap.clear();
    },

    setStates(newStates: Map<string, WidgetState>): void {
      for (const [id, state] of newStates.entries()) {
        statesMap.set(id, state);
      }
    },

    resetButtons(): void {
      for (const [id, state] of statesMap.entries()) {
        if (state.type === 'button') {
          state.value = false;
          statesMap.set(id, state);
        } else if (state.type === 'form') {
          state.value = false;
          statesMap.set(id, state);
        }
      }
    },
  };

  return sessionState;
}

/**
 * Create a new session
 * @param id Session ID
 * @param pageId Page ID
 * @returns Session
 */
export function newSession(id: string, pageId: string): Session {
  return {
    id,
    pageId,
    state: createSessionState(),
  };
}

/**
 * Create a session manager
 * @returns Session manager
 */
export function createSessionManager(): SessionManager {
  const activeSessions = new Map<string, Session>();
  const disconnectedSessions = new Map<string, DisconnectedSession>();

  /**
   * Remove the oldest disconnected session
   */
  function removeOldestDisconnectedSession(): void {
    if (disconnectedSessions.size === 0) {
      return;
    }

    let oldestId: string | null = null;
    let oldestTime: Date | null = null;

    for (const [id, ds] of disconnectedSessions.entries()) {
      if (oldestTime === null || ds.disconnectedAt < oldestTime) {
        oldestId = id;
        oldestTime = ds.disconnectedAt;
      }
    }

    if (oldestId !== null) {
      disconnectedSessions.delete(oldestId);
    }
  }

  return {
    /**
     * Set a session
     * @param session Session
     */
    setSession(session: Session): void {
      // Check if there's a disconnected session with the same ID
      const disconnectedSession = disconnectedSessions.get(session.id);
      if (disconnectedSession) {
        // Restore the state from the disconnected session
        session.state = disconnectedSession.session.state;
        disconnectedSessions.delete(session.id);
      }

      activeSessions.set(session.id, session);
    },

    /**
     * Get a session
     * @param id Session ID
     * @returns Session or undefined if not found
     */
    getSession(id: string): Session | undefined {
      return activeSessions.get(id);
    },

    /**
     * Disconnect a session
     * @param id Session ID
     */
    disconnectSession(id: string): void {
      const session = activeSessions.get(id);
      if (session) {
        // Check if we need to remove an old disconnected session
        if (disconnectedSessions.size >= MAX_DISCONNECTED_SESSIONS) {
          removeOldestDisconnectedSession();
        }

        // Add the session to disconnected sessions
        disconnectedSessions.set(id, {
          session,
          disconnectedAt: new Date(),
        });

        // Remove the session from active sessions
        activeSessions.delete(id);

        // Schedule removal of the disconnected session after TTL
        setTimeout(() => {
          disconnectedSessions.delete(id);
        }, DISCONNECTED_SESSION_TTL);
      }
    },
  };
}
