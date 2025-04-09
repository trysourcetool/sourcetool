import {
  Message,
  MessageSchema,
} from '@trysourcetool/proto/websocket/v1/message';
import * as logger from '../../logger';
import { MessageHandlerFunc, WebSocketClient } from '../websocket';
import { create } from '@bufbuild/protobuf';

/**
 * Mock WebSocket client implementation
 */
export class MockClient implements WebSocketClient {
  private handler: MessageHandlerFunc | null = null;
  private messages: Message[] = [];
  private doneResolve: () => void = () => {};
  private done: Promise<void>;

  /**
   * Constructor
   */
  constructor() {
    this.done = new Promise<void>((resolve) => {
      this.doneResolve = resolve;
    });
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
      const msg = this.newMessage(id, payload);
      this.messages.push(msg);
      if (this.handler) {
        this.handler(msg);
      }
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
  public enqueueWithResponse(
    id: string,
    payload: Message['type']['value'],
  ): Promise<Message> {
    try {
      const msg = this.newMessage(id, payload);
      this.messages.push(msg);
      return Promise.resolve(msg);
    } catch (err) {
      return Promise.reject(err);
    }
  }

  /**
   * Get all messages
   * @returns Messages
   */
  public getMessages(): Message[] {
    return this.messages;
  }

  /**
   * Close the WebSocket connection
   */
  public close(): void {
    this.doneResolve();
  }

  /**
   * Wait for the WebSocket connection to close
   */
  public wait(): Promise<void> {
    return this.done;
  }

  /**
   * Create a new message
   * @param id Message ID
   * @param payload Message payload
   * @returns Message
   */
  private newMessage(id: string, payload: Message['type']['value']): Message {
    const msg: {
      id: string;
      type?: Message['type'];
    } = {
      id,
    };

    // Set the message type based on the payload type
    if (payload?.$typeName === 'sourcetool.websocket.v1.InitializeHost') {
      msg.type = { case: 'initializeHost', value: payload };
    } else if (
      payload?.$typeName === 'sourcetool.websocket.v1.InitializeClient'
    ) {
      msg.type = { case: 'initializeClient', value: payload };
    } else if (payload?.$typeName === 'sourcetool.websocket.v1.RenderWidget') {
      msg.type = { case: 'renderWidget', value: payload };
    } else if (payload?.$typeName === 'sourcetool.websocket.v1.RerunPage') {
      msg.type = { case: 'rerunPage', value: payload };
    } else if (payload?.$typeName === 'sourcetool.websocket.v1.CloseSession') {
      msg.type = { case: 'closeSession', value: payload };
    } else if (
      payload?.$typeName === 'sourcetool.websocket.v1.ScriptFinished'
    ) {
      msg.type = { case: 'scriptFinished', value: payload };
    } else if (payload?.$typeName === 'sourcetool.exception.v1.Exception') {
      msg.type = { case: 'exception', value: payload };
    } else {
      throw new Error(`Unsupported message type: ${payload?.$typeName}`);
    }

    return create(MessageSchema, msg);
  }
}

/**
 * Create a new mock WebSocket client
 * @returns Mock WebSocket client
 */
export function createMockClient(): MockClient {
  return new MockClient();
}
