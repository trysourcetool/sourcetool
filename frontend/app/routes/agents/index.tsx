import { PageHeader } from '@/components/common/page-header';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useDispatch, useSelector } from '@/store';
import { environmentsStore } from '@/store/modules/environments';
import { Bot, MessageSquare } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { Link, createFileRoute } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from '@/components/ui/select';
import { SelectValue } from '@radix-ui/react-select';

// Mock data for agents
const mockAgents = [
  { 
    id: '1', 
    name: 'Assistant AI', 
    description: 'General purpose AI assistant for everyday tasks',
    status: 'online', 
    avatar: null,
    lastActive: new Date(),
    messagesCount: 156
  },
  { 
    id: '2', 
    name: 'Code Helper', 
    description: 'Specialized in code review and debugging',
    status: 'online', 
    avatar: null,
    lastActive: new Date(),
    messagesCount: 89
  },
  { 
    id: '3', 
    name: 'Data Analyst', 
    description: 'Expert in data analysis and visualization',
    status: 'offline', 
    avatar: null,
    lastActive: new Date(Date.now() - 3600000), // 1 hour ago
    messagesCount: 234
  },
];

export default function Agents() {
  const isInitialLoading = useRef(false);
  const [selectedEnvironmentId, setSelectedEnvironmentId] = useState<
    string | null
  >(null);
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const environments = useSelector(environmentsStore.selector.getEnvironments);

  // TODO: Consider using redux-persist if localStorage is used frequently
  const setLocalStorageSelectedEnvironmentId = (environmentId: string) => {
    localStorage.setItem('selectedEnvironmentId', environmentId);
  };

  const getLocalStorageSelectedEnvironmentId = (): string | null => {
    const environmentId = localStorage.getItem('selectedEnvironmentId');
    return environmentId || null;
  };

  const handleSelectEnvironment = async (environmentId: string) => {
    setSelectedEnvironmentId(environmentId);
    setLocalStorageSelectedEnvironmentId(environmentId);
    // TODO: Add agent filtering by environment when backend is ready
    // await dispatch(agentsStore.asyncActions.listAgents({ environmentId }));
  };

  useEffect(() => {
    if (setBreadcrumbsState) {
      setBreadcrumbsState([{ label: t('breadcrumbs_agents'), to: '/agents' }]);
    }
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        const resultActions = await Promise.all([
          dispatch(environmentsStore.asyncActions.listEnvironments()),
        ]);
        if (
          environmentsStore.asyncActions.listEnvironments.fulfilled.match(
            resultActions[0],
          )
        ) {
          const localStorageEnvironmentId =
            getLocalStorageSelectedEnvironmentId();
          let environmentId =
            localStorageEnvironmentId ||
            resultActions[0].payload.environments[0].id;
          const hasEnvironmentId = resultActions[0].payload.environments.some(
            (e) => e.id === environmentId,
          );
          if (!hasEnvironmentId) {
            environmentId = resultActions[0].payload.environments[0].id;
          }
          setSelectedEnvironmentId(environmentId);
          if (
            !localStorageEnvironmentId ||
            localStorageEnvironmentId !== environmentId
          ) {
            setLocalStorageSelectedEnvironmentId(environmentId);
          }
          // TODO: Load agents for selected environment
          // await dispatch(agentsStore.asyncActions.listAgents({ environmentId }));
        }
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch]);

  return (
    <div className="flex h-full flex-col">
      <PageHeader 
        label={t('routes_agents_page_header')} 
      />
      
      <div className="p-4 md:p-6">
        <div className="mb-6 flex flex-col items-start justify-between gap-3 md:flex-row md:items-center">
          <div className="flex gap-2 text-lg font-bold">
            {mockAgents.length} Agents in
            <div className="flex items-center gap-2">
              <div
                className="size-3 rounded-full"
                style={{
                  backgroundColor: environments.find(
                    (e) => e.id === selectedEnvironmentId,
                  )?.color,
                }}
              />
              {environments.find((e) => e.id === selectedEnvironmentId)?.name}
            </div>
          </div>

          <div className="w-full md:max-w-72">
            {selectedEnvironmentId && (
              <Select
                value={selectedEnvironmentId ?? ''}
                onValueChange={handleSelectEnvironment}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a environment" />
                </SelectTrigger>
                <SelectContent>
                  {environments.map((environment) => (
                    <SelectItem key={environment.id} value={environment.id}>
                      {environment.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {mockAgents.map((agent) => (
            <Link
              key={agent.id}
              to="/agents/$agentId"
              params={{ agentId: agent.id }}
              className="block transition-transform hover:scale-[1.02]"
            >
              <Card className="h-full hover:shadow-lg">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                      <div className="rounded-full bg-primary/10 p-3">
                        <Bot className="h-6 w-6 text-primary" />
                      </div>
                      <div>
                        <CardTitle className="text-lg">{agent.name}</CardTitle>
                        <Badge
                          variant={agent.status === 'online' ? 'default' : 'secondary'}
                          className="mt-1"
                        >
                          {agent.status}
                        </Badge>
                      </div>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <CardDescription className="mb-4">
                    {agent.description}
                  </CardDescription>
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <MessageSquare className="h-4 w-4" />
                      <span>{agent.messagesCount} messages</span>
                    </div>
                    <span>
                      {agent.status === 'online' 
                        ? 'Active now' 
                        : `Last active ${new Date(agent.lastActive).toLocaleTimeString()}`
                      }
                    </span>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_auth/agents/')({
  component: Agents,
});