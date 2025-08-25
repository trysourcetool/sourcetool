import { useState, useEffect, useRef, useCallback } from 'react';
import { v4 as uuidv4 } from 'uuid';
import {
  MessageSchema,
  ChatMessageSchema,
  type InitializeAgentChatJson,
  type ChatMessage,
  type ToolCall,
  type ToolResult,
  ChatMessage_Role,
} from '@/pb/ts/websocket/v1/message_pb';
import { create, toBinary } from '@bufbuild/protobuf';

interface ToolCallStatus {
  toolId: string;
  toolName: string;
  status: 'starting' | 'running' | 'completed' | 'error';
  parameters?: string;
  result?: string;
  error?: string;
  startTime?: number;
  duration?: number;
}

interface UseAgentWebSocketReturn {
  sessionId: string | null;
  isConnected: boolean;
  sendMessage: (message: string, conversationHistory?: any[]) => void;
  streamingContent: string;
  isStreaming: boolean;
  toolCall: ToolCall | null;
  toolResult: ToolResult | null;
  activeToolCalls: Map<string, ToolCallStatus>;
  completedToolCalls: Map<string, ToolCallStatus>;
  error: string | null;
}

export function useAgentWebSocket(agentId: string): UseAgentWebSocketReturn {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [streamingContent, setStreamingContent] = useState('');
  const [isStreaming, setIsStreaming] = useState(false);
  const [toolCall, setToolCall] = useState<ToolCall | null>(null);
  const [toolResult, setToolResult] = useState<ToolResult | null>(null);
  const [activeToolCalls, setActiveToolCalls] = useState<
    Map<string, ToolCallStatus>
  >(new Map());
  const [completedToolCalls, setCompletedToolCalls] = useState<
    Map<string, ToolCallStatus>
  >(new Map());
  const [error, setError] = useState<string | null>(null);
  const isInitialized = useRef(false);
  const initializationTimeout = useRef<NodeJS.Timeout | null>(null);

  // Listen for agent messages from the main WebSocket controller
  useEffect(() => {
    const handleAgentMessage = (event: CustomEvent) => {
      const message = event.detail;
      console.log('Agent WebSocket message received:', message);

      if (message.initializeAgentChatCompleted) {
        setSessionId(message.initializeAgentChatCompleted.sessionId);
        setError(null);
        console.log(
          'Agent chat session initialized:',
          message.initializeAgentChatCompleted.sessionId,
        );
      }

      if (message.agentResponse) {
        const response = message.agentResponse;

        switch (response.type) {
          case 'RESPONSE_TYPE_TEXT_CHUNK':
            if (response.textChunk) {
              setStreamingContent((prev) => prev + response.textChunk);
              setIsStreaming(true);
            }
            break;

          case 'RESPONSE_TYPE_TOOL_CALL':
            if (response.toolCall) {
              setToolCall(response.toolCall);
              setIsStreaming(false);
            }
            break;

          case 'RESPONSE_TYPE_TOOL_RESULT':
            if (response.toolResult) {
              setToolResult(response.toolResult);
            }
            break;

          case 'RESPONSE_TYPE_ERROR':
            if (response.errorMessage) {
              setError(response.errorMessage);
              setIsStreaming(false);
            }
            break;

          case 'RESPONSE_TYPE_TOOL_CALL_START':
            if (response.toolCallInfo) {
              const toolInfo = response.toolCallInfo;
              // Use current time but ensure it's after the last streaming message
              const now = Date.now();
              setActiveToolCalls((prev) =>
                new Map(prev).set(toolInfo.toolId, {
                  toolId: toolInfo.toolId,
                  toolName: toolInfo.toolName,
                  status: 'starting',
                  parameters: toolInfo.parameters,
                  startTime: now,
                }),
              );
            }
            break;

          case 'RESPONSE_TYPE_TOOL_CALL_COMPLETE':
            if (response.toolCallInfo) {
              const toolInfo = response.toolCallInfo;

              // Update the tool call status to completed
              setActiveToolCalls((prev) => {
                const updated = new Map(prev);
                const existing = updated.get(toolInfo.toolId);
                if (existing) {
                  const completedTool = {
                    ...existing,
                    status: 'completed' as const,
                    result: toolInfo.result,
                    duration: toolInfo.durationMs
                      ? Number(toolInfo.durationMs)
                      : undefined,
                  };

                  // Move to completed tools after a short delay to show completion status
                  setTimeout(() => {
                    setCompletedToolCalls((prevCompleted) =>
                      new Map(prevCompleted).set(
                        toolInfo.toolId,
                        completedTool,
                      ),
                    );
                    setActiveToolCalls((prevActive) => {
                      const newActive = new Map(prevActive);
                      newActive.delete(toolInfo.toolId);
                      return newActive;
                    });
                  }, 1000); // Show completed status for 1 second

                  updated.set(toolInfo.toolId, completedTool);
                }
                return updated;
              });
            }
            break;

          case 'RESPONSE_TYPE_TOOL_CALL_ERROR':
            if (response.toolCallInfo) {
              const toolInfo = response.toolCallInfo;

              // Update the tool call status to error
              setActiveToolCalls((prev) => {
                const updated = new Map(prev);
                const existing = updated.get(toolInfo.toolId);
                if (existing) {
                  const errorTool = {
                    ...existing,
                    status: 'error' as const,
                    error: toolInfo.errorMessage,
                    duration: toolInfo.durationMs
                      ? Number(toolInfo.durationMs)
                      : undefined,
                  };

                  // Move to completed tools after a longer delay to show error status
                  setTimeout(() => {
                    setCompletedToolCalls((prevCompleted) =>
                      new Map(prevCompleted).set(toolInfo.toolId, errorTool),
                    );
                    setActiveToolCalls((prevActive) => {
                      const newActive = new Map(prevActive);
                      newActive.delete(toolInfo.toolId);
                      return newActive;
                    });
                  }, 2000); // Show error status for 2 seconds

                  updated.set(toolInfo.toolId, errorTool);
                }
                return updated;
              });
            }
            break;

          default:
            console.warn('Unknown agent response type:', response.type);
        }
      }

      if (message.agentChatComplete) {
        setIsStreaming(false);
        console.log('Agent chat completed:', message.agentChatComplete);
      }
    };

    window.addEventListener(
      'agentWebSocketMessage',
      handleAgentMessage as EventListener,
    );

    return () => {
      window.removeEventListener(
        'agentWebSocketMessage',
        handleAgentMessage as EventListener,
      );
    };
  }, []);

  const isConnected = true; // Always connected through main WebSocket

  // Initialize agent chat session when connected
  useEffect(() => {
    if (isConnected && !isInitialized.current && agentId) {
      // Clear any existing timeout
      if (initializationTimeout.current) {
        clearTimeout(initializationTimeout.current);
      }

      // Delay initialization to handle StrictMode double calls
      initializationTimeout.current = setTimeout(() => {
        if (!isInitialized.current) {
          isInitialized.current = true;

          const initMessage = create(MessageSchema, {
            id: uuidv4(),
            type: {
              case: 'initializeAgentChat',
              value: {
                agentId: agentId,
              } satisfies InitializeAgentChatJson,
            },
          });

          const binaryData = toBinary(MessageSchema, initMessage);
          window.dispatchEvent(
            new CustomEvent('sendAgentMessage', { detail: binaryData }),
          );

          console.log('Initializing agent chat for:', agentId);
        }
      }, 50); // Small delay to handle React StrictMode double rendering
    }

    // Cleanup timeout on unmount or agentId change
    return () => {
      if (initializationTimeout.current) {
        clearTimeout(initializationTimeout.current);
        initializationTimeout.current = null;
      }
    };
  }, [isConnected, agentId]);

  const sendMessage = useCallback(
    (message: string, conversationHistory: any[] = []) => {
      if (!isConnected || !sessionId) {
        setError('Not connected to agent service');
        return;
      }

      // Clear previous streaming state (but keep tool call history)
      setStreamingContent('');
      setIsStreaming(false);
      setToolCall(null);
      setToolResult(null);
      // Don't clear activeToolCalls or completedToolCalls to preserve history
      setError(null);

      // Convert conversation history to the correct format
      const chatHistory: ChatMessage[] = conversationHistory.map((msg) =>
        create(ChatMessageSchema, {
          role:
            msg.role === 'user'
              ? ChatMessage_Role.USER
              : msg.role === 'assistant'
                ? ChatMessage_Role.ASSISTANT
                : msg.role === 'system'
                  ? ChatMessage_Role.SYSTEM
                  : ChatMessage_Role.TOOL,
          content: msg.content,
          toolCallId: msg.toolCallId,
          toolCalls: msg.toolCalls || [],
          timestamp: BigInt(msg.timestamp || Date.now()),
        }),
      );

      const sendMsg = create(MessageSchema, {
        id: uuidv4(),
        type: {
          case: 'sendAgentMessage',
          value: {
            sessionId,
            agentId,
            message,
            conversationHistory: chatHistory,
          },
        },
      });

      // Send message through main WebSocket controller
      const binaryData = toBinary(MessageSchema, sendMsg);
      window.dispatchEvent(
        new CustomEvent('sendAgentMessage', { detail: binaryData }),
      );

      console.log('Sending message to agent:', { sessionId, agentId, message });
    },
    [isConnected, sessionId, agentId],
  );

  // Cleanup function to close session
  const closeSession = useCallback(() => {
    if (sessionId) {
      console.log('Closing agent session:', sessionId);

      const closeMessage = create(MessageSchema, {
        id: uuidv4(),
        type: {
          case: 'closeSession',
          value: {
            sessionId: sessionId,
          },
        },
      });

      const binaryData = toBinary(MessageSchema, closeMessage);
      window.dispatchEvent(
        new CustomEvent('sendAgentMessage', { detail: binaryData }),
      );
    }
  }, [sessionId]);

  // Reset state when agent changes
  useEffect(() => {
    // Close previous session before resetting
    if (sessionId) {
      closeSession();
    }

    isInitialized.current = false;
    setSessionId(null);
    setStreamingContent('');
    setIsStreaming(false);
    setToolCall(null);
    setToolResult(null);
    setActiveToolCalls(new Map());
    setCompletedToolCalls(new Map()); // Clear completed tools when changing agent
    setError(null);
  }, [agentId]); // Removed closeSession from dependencies to avoid re-running

  // Cleanup on unmount
  useEffect(() => {
    const currentSessionId = sessionId;
    return () => {
      if (currentSessionId) {
        console.log(
          'Component unmounting, closing agent session:',
          currentSessionId,
        );
        const closeMessage = create(MessageSchema, {
          id: uuidv4(),
          type: {
            case: 'closeSession',
            value: {
              sessionId: currentSessionId,
            },
          },
        });
        const binaryData = toBinary(MessageSchema, closeMessage);
        window.dispatchEvent(
          new CustomEvent('sendAgentMessage', { detail: binaryData }),
        );
      }
    };
  }, [sessionId]);

  // Handle page unload and visibility changes
  useEffect(() => {
    const currentSessionId = sessionId;

    const handleBeforeUnload = () => {
      if (currentSessionId) {
        console.log('Page unloading, closing agent session:', currentSessionId);
        const closeMessage = create(MessageSchema, {
          id: uuidv4(),
          type: {
            case: 'closeSession',
            value: {
              sessionId: currentSessionId,
            },
          },
        });
        const binaryData = toBinary(MessageSchema, closeMessage);
        window.dispatchEvent(
          new CustomEvent('sendAgentMessage', { detail: binaryData }),
        );
      }
    };

    const handleVisibilityChange = () => {
      if (document.visibilityState === 'hidden' && currentSessionId) {
        console.log(
          'Page becoming hidden, closing agent session:',
          currentSessionId,
        );
        const closeMessage = create(MessageSchema, {
          id: uuidv4(),
          type: {
            case: 'closeSession',
            value: {
              sessionId: currentSessionId,
            },
          },
        });
        const binaryData = toBinary(MessageSchema, closeMessage);
        window.dispatchEvent(
          new CustomEvent('sendAgentMessage', { detail: binaryData }),
        );
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    document.addEventListener('visibilitychange', handleVisibilityChange);

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [sessionId]);

  return {
    sessionId,
    isConnected,
    sendMessage,
    streamingContent,
    isStreaming,
    toolCall,
    toolResult,
    activeToolCalls,
    completedToolCalls,
    error,
  };
}
