import WebSocket from 'ws';
import * as logger from './logger';

/**
 * WebSocket client configuration
 */
export interface WebSocketClientConfig {
  url: string;
  apiKey: string;
  instanceID: string;
  pingInterval: number;
  reconnectDelay: number;
  onReconnecting?: () => void;
  onReconnected?: () => void;
}

/**
 * Message handler function
 */
export type MessageHandlerFunc = (message: any) => Promise<void> | void;

/**
 * WebSocket client interface
 */
export interface WebSocketClient {
  enqueue(id: string, message: any): void;
  enqueueWithResponse(id: string, message: any): Promise<any>;
  registerHandler(handler: MessageHandlerFunc): void;
  close(): void;
  wait(): Promise<void>;
}

/**
 * Create a new message
 * @param id Message ID
 * @param payload Message payload
 * @returns Message
 */
function newMessage(id: string, payload: any): any {
  const msg: any = {
    id,
  };

  // Set the message type based on the payload type
  if (payload.initializeHost) {
    msg.type = { initializeHost: payload };
  } else if (payload.initializeClient) {
    msg.type = { initializeClient: payload };
  } else if (payload.renderWidget) {
    msg.type = { renderWidget: payload };
  } else if (payload.rerunPage) {
    msg.type = { rerunPage: payload };
  } else if (payload.closeSession) {
    msg.type = { closeSession: payload };
  } else if (payload.scriptFinished) {
    msg.type = { scriptFinished: payload };
  } else if (payload.exception) {
    msg.type = { exception: payload };
  } else {
    throw new Error(`Unsupported message type: ${typeof payload}`);
  }

  return msg;
}

/**
 * WebSocket client implementation
 */
class Client implements WebSocketClient {
  private config: WebSocketClientConfig;
  private conn: WebSocket | null = null;
  private messageQueue: any[] = [];
  private done: Promise<void>;
  private doneResolve: () => void = () => {};
  private responses: Map<
    string,
    { resolve: (value: any) => void; reject: (reason: any) => void }
  > = new Map();
  private handler: MessageHandlerFunc | null = null;

  /**
   * Constructor
   * @param config WebSocket client configuration
   */
  constructor(config: WebSocketClientConfig) {
    this.config = config;

    // Set default values
    if (!this.config.pingInterval) {
      this.config.pingInterval = 1000; // 1 second
    }
    if (!this.config.reconnectDelay) {
      this.config.reconnectDelay = 1000; // 1 second
    }

    // Initialize done promise
    this.done = new Promise<void>((resolve) => {
      this.doneResolve = resolve;
    });
  }

  /**
   * Connect to the WebSocket server
   */
  private async connect(): Promise<void> {
    try {
      // Create headers
      const headers = {
        Authorization: `Bearer ${this.config.apiKey}`,
        'X-Instance-Id': this.config.instanceID,
      };

      // Connect to the WebSocket server
      this.conn = new WebSocket(this.config.url, { headers });

      // Set up event handlers
      this.conn.on('open', () => {
        logger.info('WebSocket connection established');
        this.startPingPongLoop();
        this.startSendEnqueuedMessagesLoop();
      });

      this.conn.on('message', (data: Buffer) => {
        try {
          const msg = JSON.parse(data.toString());
          this.handleMessage(msg);
        } catch (err) {
          logger.error('Error parsing message', err);
        }
      });

      this.conn.on('close', () => {
        logger.info('WebSocket connection closed');
        this.conn = null;
        this.reconnect();
      });

      this.conn.on('error', (err) => {
        logger.error('WebSocket error', err);
        this.conn?.close();
        this.conn = null;
        this.reconnect();
      });
    } catch (err) {
      logger.error('Error connecting to WebSocket server', err);
      this.reconnect();
    }
  }

  /**
   * Reconnect to the WebSocket server
   */
  private async reconnect(): Promise<void> {
    if (this.config.onReconnecting) {
      this.config.onReconnecting();
    }

    while (true) {
      try {
        await this.connect();
        if (this.config.onReconnected) {
          this.config.onReconnected();
        }
        return;
      } catch (err) {
        logger.error('Reconnection failed, retrying', err);
        await new Promise((resolve) =>
          setTimeout(resolve, this.config.reconnectDelay),
        );
      }
    }
  }

  /**
   * Start the ping-pong loop
   */
  private startPingPongLoop(): void {
    const pingInterval = setInterval(() => {
      if (!this.conn) {
        clearInterval(pingInterval);
        return;
      }

      try {
        this.conn.ping();
      } catch (err) {
        logger.error('Ping failed', err);
        this.conn.close();
        this.conn = null;
        clearInterval(pingInterval);
        this.reconnect();
      }
    }, this.config.pingInterval);
  }

  /**
   * Start the send enqueued messages loop
   */
  private startSendEnqueuedMessagesLoop(): void {
    const sendInterval = setInterval(() => {
      if (!this.conn) {
        clearInterval(sendInterval);
        return;
      }

      if (this.messageQueue.length > 0) {
        const msg = this.messageQueue.shift();
        try {
          this.send(msg);
        } catch (err) {
          logger.error('Error sending message', err);
          this.messageQueue.unshift(msg);
        }
      }
    }, 10); // 10ms
  }

  /**
   * Send a message
   * @param msg Message
   */
  private send(msg: any): void {
    if (!this.conn) {
      throw new Error('WebSocket connection not established');
    }

    this.conn.send(JSON.stringify(msg));
  }

  /**
   * Handle a message
   * @param msg Message
   */
  private handleMessage(msg: any): void {
    // Handle responses
    const response = this.responses.get(msg.id);
    if (response) {
      response.resolve(msg);
      this.responses.delete(msg.id);
      return;
    }

    // Handle calls
    if (this.handler) {
      try {
        this.handler(msg);
      } catch (err) {
        logger.error('Error handling message', err);
      }
    } else {
      logger.error('No message handler registered');
    }
  }

  /**
   * Register a message handler
   * @param handler Message handler
   */
  public registerHandler(handler: MessageHandlerFunc): void {
    this.handler = handler;
  }

  /**
   * Enqueue a message
   * @param id Message ID
   * @param payload Message payload
   */
  public enqueue(id: string, payload: any): void {
    try {
      const msg = newMessage(id, payload);
      this.messageQueue.push(msg);
    } catch (err) {
      logger.error('Error creating message', err);
    }
  }

  /**
   * Enqueue a message and wait for a response
   * @param id Message ID
   * @param payload Message payload
   * @returns Response
   */
  public enqueueWithResponse(id: string, payload: any): Promise<any> {
    return new Promise((resolve, reject) => {
      try {
        const msg = newMessage(id, payload);
        this.responses.set(id, { resolve, reject });
        this.messageQueue.push(msg);

        // Set a timeout to reject the promise if no response is received
        setTimeout(() => {
          if (this.responses.has(id)) {
            this.responses.delete(id);
            reject(new Error('Timeout waiting for response'));
          }
        }, 30000); // 30 seconds
      } catch (err) {
        reject(err);
      }
    });
  }

  /**
   * Close the WebSocket connection
   */
  public close(): void {
    if (this.conn) {
      this.conn.close();
      this.conn = null;
    }

    this.doneResolve();
  }

  /**
   * Wait for the WebSocket connection to close
   */
  public wait(): Promise<void> {
    return this.done;
  }
}

/**
 * Create a WebSocket client
 * @param config WebSocket client configuration
 * @returns WebSocket client
 */
export function createWebSocketClient(
  config: WebSocketClientConfig,
): WebSocketClient {
  logger.info(`Creating WebSocket client for ${config.url}`);

  const client = new Client(config);

  return client;
}

/**
 * Create a mock WebSocket client
 * @returns Mock WebSocket client
 */
export function createMockWebSocketClient(): WebSocketClient {
  const messages: any[] = [];
  let handler: MessageHandlerFunc | null = null;
  let doneResolve: () => void = () => {};
  const done = new Promise<void>((resolve) => {
    doneResolve = resolve;
  });

  return {
    enqueue(id: string, payload: any): void {
      logger.debug(`Enqueuing message ${id}`, payload);

      try {
        const msg = newMessage(id, payload);
        messages.push(msg);

        if (handler) {
          handler(msg);
        }
      } catch (err) {
        logger.error('Error creating message', err);
      }
    },

    async enqueueWithResponse(id: string, payload: any): Promise<any> {
      logger.debug(`Enqueuing message with response ${id}`, payload);

      try {
        const msg = newMessage(id, payload);
        messages.push(msg);
        return msg;
      } catch (err) {
        throw err;
      }
    },

    registerHandler(h: MessageHandlerFunc): void {
      logger.debug('Registering handler');
      handler = h;
    },

    close(): void {
      logger.debug('Closing WebSocket client');
      doneResolve();
    },

    async wait(): Promise<void> {
      logger.debug('Waiting for WebSocket client');
      return done;
    },
  };
}
