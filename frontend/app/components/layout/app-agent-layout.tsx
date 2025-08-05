import { useEffect, useRef, useState, type PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '../ui/sidebar';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { ModeToggle } from '../common/mode-toggle';
import { Link, useNavigate, useParams } from '@tanstack/react-router';
import { useAuth } from '@/hooks/use-auth';
import { ArrowLeft, Bot, Loader2, MessageSquare, Plus } from 'lucide-react';
import { ENVIRONMENTS } from '@/environments';

// Mock data for agents and chats
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
  {
    id: '3',
    title: 'API design questions',
    lastMessage: 'What about REST vs GraphQL?',
    timestamp: new Date(Date.now() - 7200000),
    messageCount: 15,
  },
];

export function AppAgentLayout(props: PropsWithChildren) {
  const { agentId, chatId } = useParams({
    from: '/_agent/agents/$agentId/chat/$chatId',
  });
  const { isAuthChecked, handleNoAuthRoute } = useAuth();
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const [chats] = useState(mockChats);

  // Find the agent by ID
  const agent = mockAgents.find(a => a.id === agentId);

  useEffect(() => {
    if (isAuthChecked === 'checked' && !ENVIRONMENTS.IS_CLOUD_EDITION) {
      // Add auth check logic if needed
    }
  }, [isAuthChecked, handleNoAuthRoute]);

  const handleNewChat = () => {
    // Generate new chat ID and navigate
    const newChatId = Date.now().toString();
    navigate({
      to: '/agents/$agentId/chat/$chatId',
      params: { agentId, chatId: newChatId },
    });
  };

  if (!agent) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground">Agent not found</p>
          <Link to="/agents" className="mt-4 inline-block">
            <Button variant="outline">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Agents
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  return isAuthChecked === 'checked' ? (
    <>
      <Sidebar collapsible="icon">
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                size="lg"
                className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground w-full cursor-default"
              >
                <div className="flex flex-1 items-center gap-2 data-[state=open]:px-2 data-[state=open]:py-1">
                  <Link to={'/agents'} className="size-8">
                    <img
                      src="/images/logo-sidebar.png"
                      alt="Sourcetool"
                      className="size-full"
                    />
                  </Link>
                  <div className="flex flex-1 flex-col gap-0.5">
                    <p className="text-sidebar-foreground text-sm font-semibold">
                      {agent.name}
                    </p>
                    <div className="flex items-center gap-2">
                      <Badge
                        variant={agent.status === 'online' ? 'default' : 'secondary'}
                        className="text-xs"
                      >
                        {agent.status}
                      </Badge>
                    </div>
                  </div>
                  <ModeToggle />
                </div>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <div className="px-2 pb-2">
              <Button 
                onClick={handleNewChat}
                className="w-full justify-start" 
                variant="outline"
              >
                <Plus className="mr-2 h-4 w-4" />
                New Chat
              </Button>
            </div>
            {chats.map((chat) => (
              <SidebarMenu key={chat.id}>
                <SidebarMenuButton 
                  asChild 
                  isActive={chatId === chat.id}
                >
                  <Link 
                    to={'/agents/$agentId/chat/$chatId'} 
                    params={{ agentId, chatId: chat.id }}
                  >
                    <span className="truncate text-sm">
                      {chat.title}
                    </span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenu>
            ))}
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton asChild>
                <Link to="/agents">
                  <ArrowLeft className="h-4 w-4" />
                  <span>Back to Agents</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sidebar>
      <SidebarInset>
        <div>{props.children}</div>
      </SidebarInset>
    </>
  ) : (
    <div className="flex h-screen w-full items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}