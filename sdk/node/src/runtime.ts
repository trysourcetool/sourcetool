import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import { Page, PageManager, newPageManager } from './internal/page';
import {
  createSessionManager,
  newSession,
  SessionManager,
  WidgetState,
} from './internal/session';
import {
  createWebSocketClient,
  WebSocketClient,
} from './internal/websocket/websocket';
import {
  CloseSession,
  InitializeClient,
  InitializeHostSchema,
  Message,
  RerunPage,
} from '@trysourcetool/proto/websocket/v1/message';
import { convertTextInputProtoToState } from './textinput';
import { convertButtonProtoToState } from './button';
import { convertNumberInputProtoToState } from './numberinput';
import { convertDateInputProtoToState } from './dateinput';
import { convertDateTimeInputProtoToState } from './datetimeinput';
import { convertTimeInputProtoToState } from './timeinput';
import { convertFormProtoToState } from './form';
import { convertMarkdownProtoToState } from './markdown';
import {
  convertColumnItemProtoToState,
  convertColumnsProtoToState,
} from './columns';
import { convertCheckboxProtoToState } from './checkbox';
import { convertCheckboxGroupProtoToState } from './checkboxgroup';
import { convertRadioProtoToState } from './radio';
import { convertSelectboxProtoToState } from './selectbox';
import { convertTextAreaProtoToState } from './textarea';
import { convertTableProtoToState } from './table';
import { convertMultiSelectProtoToState } from './multiselect';
import { create } from '@bufbuild/protobuf';

/**
 * Runtime class
 */
export class Runtime {
  /**
   * WebSocket client
   */
  wsClient: WebSocketClient;

  /**
   * Session manager
   */
  sessionManager: SessionManager;

  /**
   * Page manager
   */
  pageManager: PageManager;

  /**
   * Constructor
   * @param wsClient WebSocket client
   * @param sessionManager Session manager
   * @param pageManager Page manager
   */
  constructor(
    wsClient: WebSocketClient,
    sessionManager: SessionManager,
    pageManager: PageManager,
  ) {
    this.wsClient = wsClient;
    this.sessionManager = sessionManager;
    this.pageManager = pageManager;
  }

  /**
   * Send initialize host message
   * @param apiKey API key
   * @param pages Pages
   */
  sendInitializeHost(apiKey: string, pages: Record<string, Page>): void {
    const pagesPayload = Object.entries(pages).map(([, page]) => ({
      id: page.id,
      name: page.name,
      route: page.route,
      path: page.path,
      groups: page.accessGroups,
    }));

    const msg = create(InitializeHostSchema, {
      apiKey,
      sdkName: 'sourcetool-js',
      sdkVersion: '0.0.1',
      pages: pagesPayload,
    });

    this.wsClient
      .enqueueWithResponse(uuidv4(), msg)
      .then((resp: Message) => {
        if (resp.type.case === 'exception') {
          console.error(
            'Initialize host message failed:',
            resp.type.value.message,
          );
          throw new Error('Initialize host message failed');
        }
        console.info('Initialize host message sent:', resp);
      })
      .catch((err: Error) => {
        console.error('Failed to send initialize host message:', err);
        throw err;
      });
  }

  /**
   * Handle initialize client message
   * @param msg Message
   * @returns Promise
   */
  async handleInitializeClient(msg: InitializeClient): Promise<void> {
    if (!msg.sessionId) {
      throw new Error('Session ID is required');
    }

    const sessionID = msg.sessionId;
    const pageID = msg.pageId;

    // Create session
    const session = newSession(sessionID, pageID);

    this.sessionManager.setSession(session);

    // Get page
    const page = this.pageManager.getPage(pageID);
    if (!page) {
      throw new Error(`Page not found: ${pageID}`);
    }

    // Create UI builder
    const ui = new UIBuilder(this, session, page);

    try {
      // Run page
      await page.run(ui);

      // Send script finished message
      this.wsClient.enqueue(uuidv4(), {
        sessionId: sessionID,
        status: 'SUCCESS',
      });
    } catch (err) {
      // Send script finished message with failure status
      this.wsClient.enqueue(uuidv4(), {
        sessionId: sessionID,
        status: 'FAILURE',
      });

      throw err;
    }
  }

  /**
   * Handle rerun page message
   * @param msg Message
   * @returns Promise
   */
  async handleRerunPage(msg: RerunPage): Promise<void> {
    const sessionID = msg.sessionId;
    const pageID = msg.pageId;

    // Get session
    const session = this.sessionManager.getSession(sessionID);
    if (!session) {
      throw new Error(`Session not found: ${sessionID}`);
    }

    // Get page
    const page = this.pageManager.getPage(pageID);
    if (!page) {
      throw new Error(`Page not found: ${pageID}`);
    }

    // Reset states if page changed
    if (session.pageId !== pageID) {
      session.state.resetStates();
    }

    // Update widget states
    const newWidgetStates = new Map<string, WidgetState>();
    for (const widget of msg.states) {
      const id = widget.id;

      // Convert widget state based on type
      // This is a simplified version, the actual implementation would handle all widget types
      switch (widget.type.case) {
        case 'textInput':
          const textInputState = convertTextInputProtoToState(
            id,
            widget.type.value,
          );
          if (textInputState) {
            newWidgetStates.set(id, textInputState);
          }
          break;
        case 'numberInput':
          const numberInputState = convertNumberInputProtoToState(
            id,
            widget.type.value,
          );
          if (numberInputState) {
            newWidgetStates.set(id, numberInputState);
          }
          break;
        case 'dateInput':
          const dateInputState = convertDateInputProtoToState(
            id,
            widget.type.value,
          );
          if (dateInputState) {
            newWidgetStates.set(id, dateInputState);
          }
          break;
        case 'dateTimeInput':
          const dateTimeInputState = convertDateTimeInputProtoToState(
            id,
            widget.type.value,
          );
          if (dateTimeInputState) {
            newWidgetStates.set(id, dateTimeInputState);
          }
          break;
        case 'timeInput':
          const timeInputState = convertTimeInputProtoToState(
            id,
            widget.type.value,
          );
          if (timeInputState) {
            newWidgetStates.set(id, timeInputState);
          }
          break;
        case 'form':
          const formState = convertFormProtoToState(id, widget.type.value);
          if (formState) {
            newWidgetStates.set(id, formState);
          }
          break;
        case 'button':
          const buttonState = convertButtonProtoToState(id, widget.type.value);
          if (buttonState) {
            newWidgetStates.set(id, buttonState);
          }
          break;
        case 'markdown':
          const markdownState = convertMarkdownProtoToState(
            id,
            widget.type.value,
          );
          if (markdownState) {
            newWidgetStates.set(id, markdownState);
          }
          break;
        case 'columns':
          const columnsState = convertColumnsProtoToState(
            id,
            widget.type.value,
          );
          if (columnsState) {
            newWidgetStates.set(id, columnsState);
          }
          break;
          break;
        case 'columnItem':
          const columnItemState = convertColumnItemProtoToState(
            id,
            widget.type.value,
          );
          if (columnItemState) {
            newWidgetStates.set(id, columnItemState);
          }
          break;
        case 'checkbox':
          const checkboxState = convertCheckboxProtoToState(
            id,
            widget.type.value,
          );
          if (checkboxState) {
            newWidgetStates.set(id, checkboxState);
          }
          break;
        case 'checkboxGroup':
          const checkboxGroupState = convertCheckboxGroupProtoToState(
            id,
            widget.type.value,
          );
          if (checkboxGroupState) {
            newWidgetStates.set(id, checkboxGroupState);
          }
          break;
        case 'radio':
          const radioState = convertRadioProtoToState(id, widget.type.value);
          if (radioState) {
            newWidgetStates.set(id, radioState);
          }
          break;
        case 'selectbox':
          const selectboxState = convertSelectboxProtoToState(
            id,
            widget.type.value,
          );
          if (selectboxState) {
            newWidgetStates.set(id, selectboxState);
          }
          break;
        case 'multiSelect':
          const multiSelectState = convertMultiSelectProtoToState(
            id,
            widget.type.value,
          );
          if (multiSelectState) {
            newWidgetStates.set(id, multiSelectState);
          }
          break;
        case 'table':
          const tableState = convertTableProtoToState(id, widget.type.value);
          if (tableState) {
            newWidgetStates.set(id, tableState);
          }
          break;
        case 'textArea':
          const textAreaState = convertTextAreaProtoToState(
            id,
            widget.type.value,
          );
          if (textAreaState) {
            newWidgetStates.set(id, textAreaState);
          }
          break;

        default:
          throw new Error(`Unknown widget type: ${widget.type}`);
      }
    }

    console.log('newWidgetStates', newWidgetStates);

    // Set new widget states
    session.state.setStates(newWidgetStates);

    // Create UI builder
    const ui = new UIBuilder(this, session, page);

    try {
      // Run page
      await page.run(ui);

      // Send script finished message
      this.wsClient.enqueue(uuidv4(), {
        sessionId: sessionID,
        status: 'SUCCESS',
      });

      // Reset buttons
      session.state.resetButtons();
    } catch (err) {
      // Send script finished message with failure status
      this.wsClient.enqueue(uuidv4(), {
        sessionId: sessionID,
        status: 'FAILURE',
      });

      throw err;
    }
  }

  /**
   * Handle close session message
   * @param msg Message
   */
  handleCloseSession(msg: CloseSession): void {
    const sessionID = msg.sessionId;
    this.sessionManager.disconnectSession(sessionID);
  }

  /**
   * Send exception
   * @param id Message ID
   * @param sessionID Session ID
   * @param err Error
   */
  sendException(id: string, sessionID: string, err: Error): void {
    const exception = {
      title: 'Error',
      message: err.message,
      stackTrace: err.stack,
      sessionId: sessionID,
    };

    this.wsClient.enqueue(id, exception);
  }
}

/**
 * Start the runtime
 * @param apiKey API key
 * @param endpoint Endpoint URL
 * @param pages Pages
 * @returns Promise
 */
export async function startRuntime(
  apiKey: string,
  endpoint: string,
  pages: Record<string, Page>,
): Promise<Runtime> {
  // Create session manager
  const sessionManager = createSessionManager();

  // Create page manager
  const pageManager = newPageManager(pages);

  // Create WebSocket client
  const wsClient = createWebSocketClient({
    url: endpoint,
    apiKey,
    instanceID: uuidv4(),
    pingInterval: 1000,
    reconnectDelay: 1000,
    onReconnecting: () => {
      console.info('Reconnecting...');
    },
    onReconnected: () => {
      console.info('Reconnected!');
      runtime.sendInitializeHost(apiKey, pages);
    },
  });

  // Create runtime
  const runtime = new Runtime(wsClient, sessionManager, pageManager);

  // Register message handlers
  wsClient.registerHandler(async (msg) => {
    try {
      switch (msg.type.case) {
        case 'initializeClient':
          await runtime.handleInitializeClient(msg.type.value);
          break;
        case 'rerunPage':
          await runtime.handleRerunPage(msg.type.value);
          break;
        case 'closeSession':
          runtime.handleCloseSession(msg.type.value);
          break;
        default:
          throw new Error(`Unknown message type: ${msg.type.case}`);
      }
    } catch (err) {
      console.error('Error processing message:', err);
      if (msg.type.case === 'initializeClient') {
        runtime.sendException(
          msg.id,
          msg.type.value.sessionId ?? '',
          err as Error,
        );
      } else if (msg.type.case === 'rerunPage') {
        runtime.sendException(msg.id, msg.type.value.sessionId, err as Error);
      } else if (msg.type.case === 'closeSession') {
        runtime.sendException(msg.id, msg.type.value.sessionId, err as Error);
      }
    }
  });

  // Send initialize host message
  runtime.sendInitializeHost(apiKey, pages);

  return runtime;
}
