import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import { Page, PageManager, newPageManager } from './internal/page';
import { createSessionManager, newSession } from './internal/session';
import { Message, toJson } from '@bufbuild/protobuf';
import { createWebSocketClient, WebSocketClient } from './internal/websocket';
import { CloseSession, InitializeClient, MessageSchema, RerunPage } from '@trysourcetool/proto/websocket/v1/message';
import { convertTextInputProtoToState } from './textinput';
import { convertButtonProtoToState } from './button';
import { convertNumberInputProtoToState } from './numberinput';
import { convertDateInputProtoToState } from './dateinput';
import { convertDateTimeInputProtoToState } from './datetimeinput';
import { convertTimeInputProtoToState } from './timeinput';
import { convertFormProtoToState } from './form';
import { convertMarkdownProtoToState } from './markdown';
import { convertColumnItemProtoToState, convertColumnsProtoToState } from './columns';
import { convertCheckboxProtoToState } from './checkbox';
import { convertCheckboxGroupProtoToState } from './checkboxgroup';
import { convertRadioProtoToState } from './radio';
import { convertSelectboxProtoToState } from './selectbox';
import { convertTextAreaProtoToState } from './textarea';
import { convertTableProtoToState } from './table';
import { convertMultiSelectProtoToState } from './multiselect';

/**
 * Runtime class
 */
export class Runtime {
  /**
   * WebSocket client
   */
  wsClient: any;

  /**
   * Session manager
   */
  sessionManager: any;

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
  constructor(wsClient: WebSocketClient, sessionManager: any, pageManager: PageManager) {
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
    const pagesPayload = Object.entries(pages).map(([_, page]) => ({
      id: page.id,
      name: page.name,
      route: page.route,
      path: page.path,
      groups: page.accessGroups,
    }));

    const msg = {
      apiKey,
      sdkName: 'sourcetool-node',
      sdkVersion: '0.1.0',
      pages: pagesPayload,
    };

    this.wsClient
      .enqueueWithResponse(uuidv4(), msg)
      .then((resp: any) => {
        if (resp.exception) {
          console.error(
            'Initialize host message failed:',
            resp.exception.message,
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
    if (session.pageID !== pageID) {
      session.state.resetStates();
    }

    // Update widget states
    const newWidgetStates: Record<string, any> = {};
    for (const widget of msg.states) {
      const id = widget.id;

      // Convert widget state based on type
      // This is a simplified version, the actual implementation would handle all widget types
      switch (widget.type.case) {
        case 'textInput':
          newWidgetStates[id] = convertTextInputProtoToState(id, widget.type.value);
          break;
        case 'numberInput':
          newWidgetStates[id] = convertNumberInputProtoToState(id, widget.type.value);
          break;
        case 'dateInput':
          newWidgetStates[id] = convertDateInputProtoToState(id, widget.type.value);
          break;
        case 'dateTimeInput':
          newWidgetStates[id] = convertDateTimeInputProtoToState(id, widget.type.value);
          break;
        case 'timeInput':
          newWidgetStates[id] = convertTimeInputProtoToState(id, widget.type.value);
          break;
        case 'form':
          newWidgetStates[id] = convertFormProtoToState(id, widget.type.value);
          break;
        case 'button':
          newWidgetStates[id] = convertButtonProtoToState(id, widget.type.value);
          break;
        case 'markdown':
          newWidgetStates[id] = convertMarkdownProtoToState(id, widget.type.value);
          break;
        case 'columns':
          newWidgetStates[id] = convertColumnsProtoToState(id, widget.type.value);
          break;
        case 'columnItem':
          newWidgetStates[id] = convertColumnItemProtoToState(id, widget.type.value);
          break;
        case 'checkbox':
          newWidgetStates[id] = convertCheckboxProtoToState(id, widget.type.value);
          break;
        case 'checkboxGroup':
          newWidgetStates[id] = convertCheckboxGroupProtoToState(id, widget.type.value);
          break;
        case 'radio':
          newWidgetStates[id] = convertRadioProtoToState(id, widget.type.value);
          break;
        case 'selectbox':
          newWidgetStates[id] = convertSelectboxProtoToState(id, widget.type.value);
          break;
        case 'multiSelect':
          newWidgetStates[id] = convertMultiSelectProtoToState(id, widget.type.value);
          break;
        case 'table':
          newWidgetStates[id] = convertTableProtoToState(id, widget.type.value);
          break;
        case 'textArea':
          newWidgetStates[id] = convertTextAreaProtoToState(id, widget.type.value);
          break;

        default:
          throw new Error(`Unknown widget type: ${widget.type}`);
      }
    }

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
