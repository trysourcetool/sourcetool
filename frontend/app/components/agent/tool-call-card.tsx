import React, { useState } from 'react';
import { Wrench, Clock, CheckCircle, XCircle, ChevronDown, ChevronUp } from 'lucide-react';
import { cn } from '@/lib/utils';

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

interface ToolCallCardProps {
  toolCall: ToolCallStatus;
}

function ToolStatusBadge({ status }: { status: ToolCallStatus['status'] }) {
  const statusConfig = {
    starting: {
      icon: <Clock className="h-3 w-3" />,
      text: 'Starting',
      className: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
    },
    running: {
      icon: <Clock className="h-3 w-3 animate-spin" />,
      text: 'Running',
      className: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
    },
    completed: {
      icon: <CheckCircle className="h-3 w-3" />,
      text: 'Completed',
      className: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
    },
    error: {
      icon: <XCircle className="h-3 w-3" />,
      text: 'Error',
      className: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
    },
  };

  const config = statusConfig[status];

  return (
    <span className={cn(
      'inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs font-medium',
      config.className
    )}>
      {config.icon}
      {config.text}
    </span>
  );
}

function formatDuration(duration: number): string {
  if (duration < 1000) {
    return `${duration}ms`;
  }
  return `${(duration / 1000).toFixed(1)}s`;
}

export function ToolCallCard({ toolCall }: ToolCallCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  const hasDetails = toolCall.parameters || toolCall.result || toolCall.error;

  return (
    <div className={cn(
      'border rounded-lg p-3 transition-colors',
      toolCall.status === 'error' 
        ? 'bg-red-50 border-red-200 dark:bg-red-950 dark:border-red-800'
        : 'bg-blue-50 border-blue-200 dark:bg-blue-950 dark:border-blue-800'
    )}>
      <div className="flex items-center gap-2">
        <Wrench className={cn(
          'h-4 w-4',
          toolCall.status === 'error' ? 'text-red-600' : 'text-blue-600'
        )} />
        <span className="font-medium text-sm">{toolCall.toolName}</span>
        <ToolStatusBadge status={toolCall.status} />
        <div className="flex-1" />
        {toolCall.duration && (
          <span className="text-xs text-muted-foreground">
            {formatDuration(toolCall.duration)}
          </span>
        )}
        {hasDetails && (
          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className="text-muted-foreground hover:text-foreground transition-colors"
          >
            {isExpanded ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
          </button>
        )}
      </div>

      {isExpanded && hasDetails && (
        <div className="mt-3 space-y-3">
          {toolCall.parameters && (
            <div>
              <div className="text-xs font-medium text-muted-foreground mb-1">
                Parameters
              </div>
              <pre className="text-xs p-2 bg-muted rounded border overflow-x-auto">
                {(() => {
                  try {
                    return JSON.stringify(JSON.parse(toolCall.parameters), null, 2);
                  } catch {
                    return toolCall.parameters;
                  }
                })()}
              </pre>
            </div>
          )}

          {toolCall.result && (
            <div>
              <div className="text-xs font-medium text-muted-foreground mb-1">
                Result
              </div>
              <div className="text-xs p-2 bg-muted rounded border">
                {(() => {
                  try {
                    const parsed = JSON.parse(toolCall.result);
                    return (
                      <pre className="whitespace-pre-wrap">
                        {JSON.stringify(parsed, null, 2)}
                      </pre>
                    );
                  } catch {
                    return <span className="whitespace-pre-wrap">{toolCall.result}</span>;
                  }
                })()}
              </div>
            </div>
          )}

          {toolCall.error && (
            <div>
              <div className="text-xs font-medium text-red-600 dark:text-red-400 mb-1">
                Error
              </div>
              <div className="text-xs p-2 bg-red-100 dark:bg-red-900 rounded border text-red-800 dark:text-red-200">
                {toolCall.error}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}