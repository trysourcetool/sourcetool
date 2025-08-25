import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Send, Bot, User, Loader2, AlertCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { useAgentWebSocket } from '@/hooks/use-agent-websocket';
import { ToolCallCard } from '@/components/agent/tool-call-card';
import { cn } from '@/lib/utils';
import type { Agent } from '@/store/modules/agents';

interface AgentChatProps {
  agentId: string;
  agent: Agent;
}

interface ChatMessage {
  id: string;
  role: 'user' | 'assistant' | 'system' | 'tool';
  content: string;
  toolCalls?: Array<{
    id: string;
    name: string;
    arguments: string;
  }>;
  timestamp: number;
  isStreaming?: boolean;
}

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

interface DisplayItem {
  id: string;
  type: 'message' | 'toolcall';
  timestamp: number;
  data: ChatMessage | ToolCallStatus;
}

export function AgentChat({ agentId, agent }: AgentChatProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isComposing, setIsComposing] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const {
    isConnected,
    sendMessage,
    streamingContent,
    isStreaming,
    toolCall,
    toolResult,
    activeToolCalls,
    completedToolCalls,
    error: wsError,
  } = useAgentWebSocket(agentId);

  // Scroll to bottom when new messages arrive
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  // Merge messages and tool calls into a single chronological display list
  const displayItems: DisplayItem[] = React.useMemo(() => {
    const items: DisplayItem[] = [];
    
    // Add all messages
    messages.forEach(msg => {
      items.push({
        id: msg.id,
        type: 'message',
        timestamp: msg.timestamp,
        data: msg
      });
    });
    
    // Add active tool calls with proper timestamp to maintain chronological order
    Array.from(activeToolCalls.values()).forEach((toolCall) => {
      items.push({
        id: toolCall.toolId,
        type: 'toolcall',
        timestamp: toolCall.startTime || Date.now(),
        data: toolCall
      });
    });
    
    // Add completed tool calls to maintain them in history
    Array.from(completedToolCalls.values()).forEach((toolCall) => {
      items.push({
        id: `completed-${toolCall.toolId}`, // Different ID to avoid conflicts
        type: 'toolcall',
        timestamp: toolCall.startTime || Date.now(),
        data: toolCall
      });
    });
    
    // Sort by timestamp to maintain chronological order
    return items.sort((a, b) => a.timestamp - b.timestamp);
  }, [messages, activeToolCalls, completedToolCalls]);

  useEffect(() => {
    scrollToBottom();
  }, [displayItems]);

  // Handle streaming content
  useEffect(() => {
    if (streamingContent && isStreaming) {
      setMessages((prev) => {
        const lastMessage = prev[prev.length - 1];
        if (lastMessage && lastMessage.role === 'assistant' && lastMessage.isStreaming) {
          // Update existing streaming message - keep original timestamp
          return [
            ...prev.slice(0, -1),
            { ...lastMessage, content: streamingContent },
          ];
        } else {
          // Create new streaming message with fixed timestamp
          const fixedTimestamp = Date.now();
          return [
            ...prev,
            {
              id: `stream-${fixedTimestamp}`,
              role: 'assistant',
              content: streamingContent,
              timestamp: fixedTimestamp, // Fixed timestamp that won't change
              isStreaming: true,
            },
          ];
        }
      });
    } else if (!isStreaming && messages.length > 0) {
      // Mark streaming as complete - keep original timestamp
      setMessages((prev) =>
        prev.map((msg) =>
          msg.isStreaming ? { ...msg, isStreaming: false } : msg
        )
      );
    }
  }, [streamingContent, isStreaming, messages.length]);

  // Handle tool calls
  useEffect(() => {
    if (toolCall) {
      setMessages((prev) => [
        ...prev,
        {
          id: `tool-${toolCall.id}`,
          role: 'assistant',
          content: `Calling tool: ${toolCall.name}`,
          toolCalls: [toolCall],
          timestamp: Date.now(),
        },
      ]);
    }
  }, [toolCall]);

  // Handle tool results
  useEffect(() => {
    if (toolResult) {
      setMessages((prev) => [
        ...prev,
        {
          id: `tool-result-${toolResult.toolCallId}`,
          role: 'tool',
          content: toolResult.success
            ? `Tool result: ${toolResult.result}`
            : `Tool error: ${toolResult.result}`,
          timestamp: Date.now(),
        },
      ]);
    }
  }, [toolResult]);

  // Handle errors
  useEffect(() => {
    if (wsError) {
      setError(wsError);
      setIsLoading(false);
    }
  }, [wsError]);

  const handleSend = useCallback(() => {
    if (!input.trim() || !isConnected || isLoading) return;

    const userMessage: ChatMessage = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: input.trim(),
      timestamp: Date.now(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);
    setError(null);

    sendMessage(input.trim(), messages);
    
    // Focus back on input
    textareaRef.current?.focus();
  }, [input, isConnected, isLoading, sendMessage, messages]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey && !isComposing) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleCompositionStart = () => {
    setIsComposing(true);
  };

  const handleCompositionEnd = () => {
    setIsComposing(false);
  };

  useEffect(() => {
    if (isStreaming && isLoading) {
      setIsLoading(false);
    }
  }, [isStreaming, isLoading]);

  return (
    <div className="flex h-full flex-col">
      {/* Connection Status */}
      {!isConnected && (
        <Alert className="m-4">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            Connecting to agent service...
          </AlertDescription>
        </Alert>
      )}

      {/* Error Display */}
      {error && (
        <Alert variant="destructive" className="m-4">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Messages Area */}
      <ScrollArea className="flex-1 p-4">
        <div className="space-y-4">
          {displayItems.length === 0 && (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <Bot className="mb-4 h-12 w-12 text-muted-foreground" />
              <h3 className="mb-2 text-lg font-medium">Start a conversation</h3>
              <p className="text-sm text-muted-foreground">
                Send a message to {agent.name} to begin
              </p>
            </div>
          )}

          {displayItems.map((item) => {
            if (item.type === 'message') {
              const message = item.data as ChatMessage;
              return (
                <div
                  key={item.id}
                  className={cn(
                    'flex gap-3',
                    message.role === 'user' && 'justify-end'
                  )}
                >
                  {message.role !== 'user' && (
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                      <Bot className="h-4 w-4" />
                    </div>
                  )}

                  <div
                    className={cn(
                      'max-w-[80%] rounded-lg px-4 py-2',
                      message.role === 'user'
                        ? 'bg-primary text-primary-foreground'
                        : message.role === 'tool'
                        ? 'bg-muted font-mono text-xs'
                        : 'bg-muted'
                    )}
                  >
                    <div className="whitespace-pre-wrap break-words">
                      {message.content}
                      {message.isStreaming && (
                        <span className="ml-1 inline-block animate-pulse">▊</span>
                      )}
                    </div>

                    {message.toolCalls && message.toolCalls.length > 0 && (
                      <div className="mt-2 space-y-1 border-t pt-2">
                        {message.toolCalls.map((tool) => (
                          <div key={tool.id} className="text-xs opacity-70">
                            <span className="font-medium">{tool.name}</span>
                            {tool.arguments && (
                              <span className="ml-2 font-mono">
                                {tool.arguments}
                              </span>
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  {message.role === 'user' && (
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary">
                      <User className="h-4 w-4 text-primary-foreground" />
                    </div>
                  )}
                </div>
              );
            } else {
              // Tool call display
              const toolCall = item.data as ToolCallStatus;
              return (
                <div key={item.id} className="flex gap-3">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                    <Bot className="h-4 w-4" />
                  </div>
                  <div className="max-w-[80%]">
                    <ToolCallCard toolCall={toolCall} />
                  </div>
                </div>
              );
            }
          })}

          {isLoading && !isStreaming && activeToolCalls.size === 0 && (
            <div className="flex gap-3">
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                <Bot className="h-4 w-4" />
              </div>
              <div className="flex items-center gap-2 rounded-lg bg-muted px-4 py-2">
                <Loader2 className="h-4 w-4 animate-spin" />
                <span className="text-sm">Thinking...</span>
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>
      </ScrollArea>

      {/* Input Area */}
      <div className="border-t p-4">
        <div className="flex gap-2">
          <Textarea
            ref={textareaRef}
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            onCompositionStart={handleCompositionStart}
            onCompositionEnd={handleCompositionEnd}
            placeholder={`Message ${agent.name}...`}
            className="min-h-[60px] resize-none"
            disabled={!isConnected || isLoading}
          />
          <Button
            onClick={handleSend}
            disabled={!input.trim() || !isConnected || isLoading}
            size="icon"
            className="h-[60px] w-[60px]"
          >
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Send className="h-4 w-4" />
            )}
          </Button>
        </div>
        <div className="mt-2 text-xs text-muted-foreground">
          Press Enter to send, Shift+Enter for new line
        </div>
      </div>
    </div>
  );
}