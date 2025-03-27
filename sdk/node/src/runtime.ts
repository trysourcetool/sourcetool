import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import { Page, PageManager, newPageManager } from './internal/page';
import { newSession } from './internal/session';

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
  constructor(wsClient: any, sessionManager: any, pageManager: PageManager) {
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
  async handleInitializeClient(msg: any): Promise<void> {
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
  async handleRerunPage(msg: any): Promise<void> {
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
      switch (widget.type) {
        case 'TextInput':
          newWidgetStates[id] = { id, value: widget.textInput.value };
          break;
        case 'Button':
          newWidgetStates[id] = { id, value: widget.button.value };
          break;
        // Add other widget types here
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
  handleCloseSession(msg: any): void {
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
  const sessionManager = {
    setSession: (session: any) => {},
    getSession: (id: string) => null,
    disconnectSession: (id: string) => {},
  };

  // Create page manager
  const pageManager = newPageManager(pages);

  // Create WebSocket client
  const wsClient = {
    enqueue: (id: string, msg: any) => {},
    enqueueWithResponse: async (id: string, msg: any) => ({}),
    registerHandler: (handler: (msg: any) => Promise<void>) => {},
    wait: async () => {},
    close: async () => {},
  };

  // Create runtime
  const runtime = new Runtime(wsClient, sessionManager, pageManager);

  // Register message handlers
  wsClient.registerHandler(async (msg: any) => {
    try {
      switch (msg.type) {
        case 'InitializeClient':
          await runtime.handleInitializeClient(msg.initializeClient);
          break;
        case 'RerunPage':
          await runtime.handleRerunPage(msg.rerunPage);
          break;
        case 'CloseSession':
          runtime.handleCloseSession(msg.closeSession);
          break;
        default:
          throw new Error(`Unknown message type: ${msg.type}`);
      }
    } catch (err) {
      if (msg.type === 'InitializeClient') {
        runtime.sendException(
          msg.id,
          msg.initializeClient.sessionId,
          err as Error,
        );
      } else if (msg.type === 'RerunPage') {
        runtime.sendException(msg.id, msg.rerunPage.sessionId, err as Error);
      } else if (msg.type === 'CloseSession') {
        runtime.sendException(msg.id, msg.closeSession.sessionId, err as Error);
      }
    }
  });

  // Send initialize host message
  runtime.sendInitializeHost(apiKey, pages);

  return runtime;
}
