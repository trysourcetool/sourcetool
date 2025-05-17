import WebSocket from 'ws';
import { create, fromBinary, toBinary } from '@bufbuild/protobuf';
import { Message, MessageSchema } from '../pb/websocket/v1/message_pb';

// WebSocket constants
const MIN_PING_INTERVAL = 100; // ms
const MAX_PING_INTERVAL = 30000; // ms
const MIN_RECONNECT_DELAY = 100; // ms
const MAX_RECONNECT_DELAY = 30000; // ms
const DEFAULT_PING_INTERVAL = 1000; // ms
const DEFAULT_RECONNECT_DELAY = 1000; // ms
const DEFAULT_QUEUE_SIZE = 250;
const MIN_QUEUE_SIZE = 50;
const MAX_QUEUE_SIZE = 1000;
const MAX_MESSAGE_RETRIES = 3;
const MESSAGE_RETRY_DELAY = 100; // ms
const BATCH_INTERVAL = 10; // ms
const SHUTDOWN_TIMEOUT = 5000; // ms
const INITIAL_RECONNECT_DELAY = 500; // ms
const MAX_RECONNECT_ATTEMPTS = 26;

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
export type MessageHandlerFunc = (message: Message) => Promise<void> | void;

/**
 * WebSocket client interface
 */
export interface WebSocketClient {
  enqueue(id: string, message: Message['type']['value']): void;
  enqueueWithResponse(
    id: string,
    message: Message['type']['value'],
  ): Promise<Message>;
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
function newMessage(
  id: string,
  payload: Message['type']['value'],
): {
  id: string;
  type?: Message['type'];
} {
  const msg: {
    id: string;
    type?: Message['type'];
  } = {
    id,
  };

  // Set the message type based on the payload type
  if (payload?.$typeName === 'websocket.v1.InitializeHost') {
    msg.type = { case: 'initializeHost', value: payload };
  } else if (payload?.$typeName === 'websocket.v1.InitializeClient') {
    msg.type = { case: 'initializeClient', value: payload };
  } else if (payload?.$typeName === 'websocket.v1.RenderWidget') {
    msg.type = { case: 'renderWidget', value: payload };
  } else if (payload?.$typeName === 'websocket.v1.RerunPage') {
    msg.type = { case: 'rerunPage', value: payload };
  } else if (payload?.$typeName === 'websocket.v1.CloseSession') {
    msg.type = { case: 'closeSession', value: payload };
  } else if (payload?.$typeName === 'websocket.v1.ScriptFinished') {
    msg.type = { case: 'scriptFinished', value: payload };
  } else if (payload?.$typeName === 'exception.v1.Exception') {
    msg.type = { case: 'exception', value: payload };
  } else {
    throw new Error(`Unsupported message type: ${payload?.$typeName}`);
  }

  return create(MessageSchema, msg);
}

function setConfigDefaults(config: WebSocketClientConfig & { queueSize?: number }): void {
  if (!config.pingInterval || config.pingInterval < MIN_PING_INTERVAL) {
    config.pingInterval = DEFAULT_PING_INTERVAL;
  }
  if (!config.reconnectDelay || config.reconnectDelay < MIN_RECONNECT_DELAY) {
    config.reconnectDelay = DEFAULT_RECONNECT_DELAY;
  }
  if (!config.queueSize || config.queueSize < MIN_QUEUE_SIZE) {
    config.queueSize = DEFAULT_QUEUE_SIZE;
  }
}

function validateConfig(config: WebSocketClientConfig & { queueSize?: number }): void {
  if (config.pingInterval < MIN_PING_INTERVAL) {
    throw new Error(`pingInterval must be at least ${MIN_PING_INTERVAL}ms`);
  }
  if (config.pingInterval > MAX_PING_INTERVAL) {
    throw new Error(`pingInterval must not exceed ${MAX_PING_INTERVAL}ms`);
  }
  if (config.reconnectDelay < MIN_RECONNECT_DELAY) {
    throw new Error(`reconnectDelay must be at least ${MIN_RECONNECT_DELAY}ms`);
  }
  if (config.reconnectDelay > MAX_RECONNECT_DELAY) {
    throw new Error(`reconnectDelay must not exceed ${MAX_RECONNECT_DELAY}ms`);
  }
  if (config.queueSize && config.queueSize < MIN_QUEUE_SIZE) {
    throw new Error(`queueSize must be at least ${MIN_QUEUE_SIZE}`);
  }
  if (config.queueSize && config.queueSize > MAX_QUEUE_SIZE) {
    throw new Error(`queueSize must not exceed ${MAX_QUEUE_SIZE}`);
  }
}

/**
 * WebSocket client implementation
 */
export class Client implements WebSocketClient {
  private config: WebSocketClientConfig & { queueSize?: number };
  private conn: WebSocket | null = null;
  private messageQueue: any[] = [];
  private done: Promise<void>;
  private doneResolve: () => void = () => {};
  private responses: Map<
    string,
    { resolve: (value: any) => void; reject: (reason: any) => void }
  > = new Map();
  private handler: MessageHandlerFunc | null = null;
  private isShutdown = false;
  private pingIntervalId: NodeJS.Timeout | null = null;
  private pongTimeoutId: NodeJS.Timeout | null = null;
  private sendIntervalId: NodeJS.Timeout | null = null;
  private isSending = false;
  private shutdownOnce = false;
  private sendingDone: Promise<void> | null = null;
  private sendingDoneResolve: (() => void) | null = null;

  /**
   * Constructor
   * @param config WebSocket client configuration
   */
  constructor(config: WebSocketClientConfig & { queueSize?: number }) {
    setConfigDefaults(config);
    validateConfig(config);
    this.config = config;

    // Initialize done promise
    this.done = new Promise<void>((resolve) => {
      this.doneResolve = resolve;
    });
  }

  public async createConnection(): Promise<void> {
    try {
      await this.connect();
      this.startPingPongLoop();
      this.startSendEnqueuedMessagesLoop();
    } catch (err) {
      console.error('[ERROR] Error creating connection', err);
      throw err;
    }
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
        console.info('[INFO] WebSocket connection established');
        this.startPingPongLoop();
        this.startSendEnqueuedMessagesLoop();
      });

      this.conn.on('message', (data: Buffer) => {
        try {
          const msg = fromBinary(MessageSchema, data);
          this.handleMessage(msg);
        } catch (err) {
          console.error('[ERROR] Error parsing message', err);
        }
      });

      this.conn.on('close', () => {
        console.info('[INFO] WebSocket connection closed');
        this.conn = null;
        this.reconnect();
      });

      this.conn.on('error', (err) => {
        console.error('[ERROR] WebSocket error', err);
        this.conn?.close();
        this.conn = null;
        this.reconnect();
      });
    } catch (err) {
      console.error('[ERROR] Error connecting to WebSocket server', err);
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

    let attempt = 0;
    let lastSuccessTime = Date.now();

    while (true) {
      // Calculate delay with exponential backoff, capped
      let delay = Math.min(INITIAL_RECONNECT_DELAY * Math.pow(2, attempt), MAX_RECONNECT_DELAY);
      // Add jitter (up to 1/4 of delay)
      const maxJitter = Math.floor(delay / 4);
      if (maxJitter > 0) {
        delay += Math.floor(Math.random() * maxJitter);
      }

      console.info(`[INFO] Attempting to reconnect (attempt ${attempt + 1}, delay ${delay}ms)`);

      try {
        await new Promise((resolve) => setTimeout(resolve, delay));
        if (this.isShutdown) {
          console.info('[INFO] Reconnection canceled during shutdown');
          return;
        }
        await this.connect();
        console.info(`[INFO] Reconnection successful (attempts: ${attempt + 1})`);
        if (this.config.onReconnected) {
          this.config.onReconnected();
        }
        lastSuccessTime = Date.now();
        return;
      } catch (err) {
        attempt++;
        console.error('[ERROR] Reconnection failed', err);
        // If max attempts reached within an hour, stop
        if (attempt >= MAX_RECONNECT_ATTEMPTS) {
          if (Date.now() - lastSuccessTime < 60 * 60 * 1000) {
            console.error('[ERROR] Max reconnection attempts reached within an hour');
            return;
          } else {
            // More than an hour has passed, reset counter
            console.warn('[WARN] Continuing reconnection attempts after an hour');
            attempt = 0;
          }
        }
      }
    }
  }

  /**
   * Start the ping-pong loop
   */
  private startPingPongLoop(): void {
    if (this.pingIntervalId) {
      clearInterval(this.pingIntervalId);
    }
    if (this.pongTimeoutId) {
      clearTimeout(this.pongTimeoutId);
    }
    if (!this.conn) return;

    // Set up pong handler
    this.conn.removeAllListeners('pong');
    this.conn.on('pong', () => {
      // On pong, reset the pong timeout
      if (this.pongTimeoutId) {
        clearTimeout(this.pongTimeoutId);
      }
      this.pongTimeoutId = setTimeout(() => {
        console.error('[ERROR] Pong timeout: no pong received');
        if (this.conn) {
          this.conn.close();
          this.conn = null;
        }
        this.reconnect();
      }, this.config.pingInterval * 2);
    });

    // Immediately set the first pong timeout
    this.pongTimeoutId = setTimeout(() => {
      console.error('[ERROR] Pong timeout: no pong received');
      if (this.conn) {
        this.conn.close();
        this.conn = null;
      }
      this.reconnect();
    }, this.config.pingInterval * 2);

    // Start ping interval
    this.pingIntervalId = setInterval(() => {
      if (!this.conn) {
        if (this.pingIntervalId) clearInterval(this.pingIntervalId);
        if (this.pongTimeoutId) clearTimeout(this.pongTimeoutId);
        return;
      }
      try {
        this.conn.ping();
      } catch (err) {
        console.error('[ERROR] Ping failed', err);
        this.conn.close();
        this.conn = null;
        if (this.pingIntervalId) clearInterval(this.pingIntervalId);
        if (this.pongTimeoutId) clearTimeout(this.pongTimeoutId);
        this.reconnect();
      }
    }, this.config.pingInterval);
  }

  /**
   * Send a message
   * @param msg Message
   */
  private send(msg: Message): void {
    if (!this.conn) {
      throw new Error('WebSocket connection not established');
    }
    const bytes = toBinary(MessageSchema, msg);
    this.conn.send(bytes);
  }

  /**
   * Start the send enqueued messages loop
   */
  private startSendEnqueuedMessagesLoop(): void {
    if (this.sendIntervalId) {
      clearInterval(this.sendIntervalId);
    }
    this.isSending = false;
    this.sendingDone = new Promise<void>((resolve) => {
      this.sendingDoneResolve = resolve;
    });

    const processQueue = async () => {
      if (this.isSending) return;
      this.isSending = true;
      let messageBuffer: Message[] = [];
      while (!this.isShutdown) {
        if (!this.conn) {
          await new Promise((resolve) => setTimeout(resolve, 1000));
          continue;
        }
        // Fill buffer
        while (this.messageQueue.length > 0) {
          messageBuffer.push(this.messageQueue.shift()!);
        }
        if (messageBuffer.length === 0) {
          await new Promise((resolve) => setTimeout(resolve, BATCH_INTERVAL));
          continue;
        }
        // Send messages with retry
        const remainingMessages: Message[] = [];
        for (const msg of messageBuffer) {
          try {
            await this.sendWithRetry(msg);
          } catch {
            remainingMessages.push(msg);
            break; // Stop batch on first hard failure
          }
          await new Promise((resolve) => setTimeout(resolve, 1));
        }
        messageBuffer = remainingMessages;
        if (remainingMessages.length === 0) {
          await new Promise((resolve) => setTimeout(resolve, BATCH_INTERVAL));
        }
      }
      this.isSending = false;
      if (this.sendingDoneResolve) {
        this.sendingDoneResolve();
        this.sendingDoneResolve = null;
      }
    };
    processQueue();
  }

  /**
   * Handle a message
   * @param msg Message
   */
  private handleMessage(msg: Message): void {
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
        console.error('[ERROR] Error handling message', err);
      }
    } else {
      console.error('[ERROR] No message handler registered');
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
  public enqueue(id: string, payload: Message['type']['value']): void {
    try {
      const msg = newMessage(id, payload);
      this.messageQueue.push(msg);
    } catch (err) {
      console.error('[ERROR] Error creating message', err);
    }
  }

  /**
   * Enqueue a message and wait for a response
   * @param id Message ID
   * @param payload Message payload
   * @returns Response
   */
  public enqueueWithResponse(
    id: string,
    payload: Message['type']['value'],
  ): Promise<Message> {
    return new Promise((resolve, reject) => {
      let timeoutId: NodeJS.Timeout | null = null;
      let handled = false;
      try {
        const msg = newMessage(id, payload);
        this.responses.set(id, {
          resolve: (resp: Message) => {
            if (handled) return;
            handled = true;
            if (timeoutId) clearTimeout(timeoutId);
            this.responses.delete(id);
            resolve(resp);
          },
          reject: (err: any) => {
            if (handled) return;
            handled = true;
            if (timeoutId) clearTimeout(timeoutId);
            this.responses.delete(id);
            reject(err);
          },
        });
        this.messageQueue.push(msg);
        // Set a timeout to reject the promise if no response is received
        timeoutId = setTimeout(() => {
          if (handled) return;
          handled = true;
          this.responses.delete(id);
          reject(new Error('Timeout waiting for response'));
        }, 30000); // 30 seconds
      } catch (err) {
        if (!handled) {
          handled = true;
          this.responses.delete(id);
          reject(err);
        }
      }
    });
  }

  private async sendWithRetry(msg: Message): Promise<void> {
    let lastErr: any = null;
    for (let attempt = 0; attempt < MAX_MESSAGE_RETRIES; attempt++) {
      try {
        this.send(msg);
        return;
      } catch (err) {
        lastErr = err;
        const delay = MESSAGE_RETRY_DELAY * Math.pow(2, attempt);
        console.warn(`[WARN] Message send failed (attempt ${attempt + 1}), retrying in ${delay}ms`, err);
        await new Promise((resolve) => setTimeout(resolve, delay));
      }
    }
    console.error('[ERROR] Failed to send message after retries', lastErr);
    throw lastErr;
  }

  /**
   * Close the WebSocket connection
   */
  public async close(): Promise<void> {
    if (this.shutdownOnce) return;
    this.shutdownOnce = true;
    this.isShutdown = true;

    // 1. Stop all intervals/timers
    if (this.sendIntervalId) {
      clearInterval(this.sendIntervalId);
      this.sendIntervalId = null;
    }
    if (this.pingIntervalId) {
      clearInterval(this.pingIntervalId);
      this.pingIntervalId = null;
    }
    if (this.pongTimeoutId) {
      clearTimeout(this.pongTimeoutId);
      this.pongTimeoutId = null;
    }

    // 2. Wait for message sending loop to finish (with timeout)
    if (this.sendingDone) {
      await Promise.race([
        this.sendingDone,
        new Promise((resolve) => setTimeout(resolve, SHUTDOWN_TIMEOUT)),
      ]);
      this.sendingDone = null;
    }

    // 3. Try to send remaining messages with retry and timeout
    if (this.messageQueue.length > 0 && this.conn) {
      const sendAll = async () => {
        for (const msg of this.messageQueue) {
          try {
            await this.sendWithRetry(msg);
          } catch {
            // Already logged in sendWithRetry
          }
        }
        this.messageQueue = [];
      };
      await Promise.race([
        sendAll(),
        new Promise((resolve) => setTimeout(resolve, SHUTDOWN_TIMEOUT)),
      ]);
    }

    // 4. Close the WebSocket connection
    if (this.conn) {
      this.conn.close();
      this.conn = null;
    }

    // 5. Clean up response promises
    for (const [id, { reject }] of this.responses.entries()) {
      reject(new Error('WebSocket client closed'));
      this.responses.delete(id);
    }

    // 6. Resolve the done promise
    this.doneResolve();
  }

  /**
   * Wait for the WebSocket connection to close
   */
  public wait(): Promise<void> {
    return this.done;
  }
}
