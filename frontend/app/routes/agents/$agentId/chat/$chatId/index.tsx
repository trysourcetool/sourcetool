import { PageHeader } from '@/components/common/page-header';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { Bot, MessageSquare, Send, User, Trash2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { createFileRoute, useParams } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/utils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { MoreHorizontal } from 'lucide-react';

// Mock data for agents
const mockAgents = [
  { 
    id: '1', 
    name: 'Assistant AI', 
    description: 'General purpose AI assistant for everyday tasks',
    status: 'online', 
    avatar: null,
  },
  { 
    id: '2', 
    name: 'Code Helper', 
    description: 'Specialized in code review and debugging',
    status: 'online', 
    avatar: null,
  },
  { 
    id: '3', 
    name: 'Data Analyst', 
    description: 'Expert in data analysis and visualization',
    status: 'offline', 
    avatar: null,
  },
];

const mockChats = [
  {
    id: '1',
    title: 'Help with React components',
    lastMessage: 'How can I optimize this component?',
    timestamp: new Date(),
    messageCount: 12,
  },
  {
    id: '2',
    title: 'Database optimization',
    lastMessage: 'Thanks for the SQL tips!',
    timestamp: new Date(Date.now() - 3600000),
    messageCount: 8,
  },
];

// Mock messages
const mockMessages = [
  {
    id: '1',
    sender: 'agent',
    content: 'Hello! How can I help you today?',
    timestamp: new Date(),
  },
  {
    id: '2',
    sender: 'user',
    content: 'I need help with my React components',
    timestamp: new Date(),
  },
  {
    id: '3',
    sender: 'agent',
    content: "I'd be happy to help with your React components. What specific issue are you facing?",
    timestamp: new Date(),
  },
];

export default function ChatDetail() {
  const { agentId, chatId } = useParams({
    from: '/_agent/agents/$agentId/chat/$chatId',
  });
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const [messages, setMessages] = useState(mockMessages);
  const [inputMessage, setInputMessage] = useState('');

  // Find the agent and chat by ID
  const agent = mockAgents.find(a => a.id === agentId);
  const chat = mockChats.find(c => c.id === chatId);

  const isNewChat = !chat;
  const chatTitle = chat?.title || 'New Chat';

  useEffect(() => {
    if (setBreadcrumbsState && agent) {
      setBreadcrumbsState([
        { label: t('breadcrumbs_agents'), to: '/agents' },
        { label: agent.name, to: `/agents/${agentId}` },
        { label: chatTitle },
      ]);
    }
  }, [setBreadcrumbsState, t, agent, agentId, chatTitle]);

  const handleSendMessage = () => {
    if (!inputMessage.trim()) return;

    const newMessage = {
      id: Date.now().toString(),
      sender: 'user' as const,
      content: inputMessage,
      timestamp: new Date(),
    };

    setMessages([...messages, newMessage]);
    setInputMessage('');

    // Simulate agent response
    setTimeout(() => {
      const agentResponse = {
        id: (Date.now() + 1).toString(),
        sender: 'agent' as const,
        content: 'I understand. Let me help you with that...',
        timestamp: new Date(),
      };
      setMessages(prev => [...prev, agentResponse]);
    }, 1000);
  };

  const handleDeleteChat = () => {
    // TODO: Implement chat deletion
    console.log('Delete chat:', chatId);
  };

  if (!agent) {
    return (
      <div className="flex h-full flex-col">
        <PageHeader label="Agent Not Found" />
        <div className="flex flex-1 items-center justify-center">
          <p className="text-muted-foreground">The requested agent could not be found.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      <PageHeader 
        label={chatTitle}
      />
      
      <div className="flex flex-1 p-4">
        <Card className="flex-1">
          <CardHeader className="border-b">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Avatar>
                  <AvatarImage src={agent.avatar || undefined} />
                  <AvatarFallback>
                    <Bot className="h-5 w-5" />
                  </AvatarFallback>
                </Avatar>
                <div>
                  <CardTitle className="flex items-center gap-2">
                    <MessageSquare className="h-5 w-5" />
                    {chatTitle}
                  </CardTitle>
                  <p className="text-sm text-muted-foreground mt-1">
                    {agent.description}
                  </p>
                </div>
              </div>
              {!isNewChat && (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={handleDeleteChat}>
                      <Trash2 className="mr-2 h-4 w-4" />
                      Delete Chat
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
            </div>
          </CardHeader>
          <CardContent className="flex h-[calc(100vh-280px)] flex-col p-0">
            {/* Messages */}
            <ScrollArea className="flex-1 p-4">
              <div className="space-y-4">
                {messages.map((message) => (
                  <div
                    key={message.id}
                    className={cn(
                      "flex gap-3",
                      message.sender === 'user' && "flex-row-reverse"
                    )}
                  >
                    <Avatar className="h-8 w-8">
                      <AvatarFallback>
                        {message.sender === 'agent' ? (
                          <Bot className="h-4 w-4" />
                        ) : (
                          <User className="h-4 w-4" />
                        )}
                      </AvatarFallback>
                    </Avatar>
                    <div
                      className={cn(
                        "rounded-lg px-4 py-2 max-w-[70%]",
                        message.sender === 'agent'
                          ? "bg-muted"
                          : "bg-primary text-primary-foreground"
                      )}
                    >
                      <p className="text-sm">{message.content}</p>
                      <p className="mt-1 text-xs opacity-70">
                        {message.timestamp.toLocaleTimeString()}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </ScrollArea>

            {/* Input Area */}
            <div className="border-t p-4">
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  handleSendMessage();
                }}
                className="flex gap-2"
              >
                <Input
                  value={inputMessage}
                  onChange={(e) => setInputMessage(e.target.value)}
                  placeholder={`Type your message to ${agent.name}...`}
                  className="flex-1"
                  disabled={agent.status === 'offline'}
                />
                <Button 
                  type="submit" 
                  size="icon"
                  disabled={agent.status === 'offline' || !inputMessage.trim()}
                >
                  <Send className="h-4 w-4" />
                </Button>
              </form>
              {agent.status === 'offline' && (
                <p className="mt-2 text-sm text-muted-foreground">
                  This agent is currently offline. Messages will be delivered when they come back online.
                </p>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_agent/agents/$agentId/chat/$chatId')({
  component: ChatDetail,
});