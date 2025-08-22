import { PageHeader } from '@/components/common/page-header';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useDispatch, useSelector } from '@/store';
import { environmentsStore } from '@/store/modules/environments';
import { agentsStore } from '@/store/modules/agents';
import { usersStore } from '@/store/modules/users';
import { apiKeysStore } from '@/store/modules/apiKeys';
import { Bot, AlertCircle } from 'lucide-react';
import { useEffect, useRef, useState, useMemo } from 'react';
import { Link, createFileRoute } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from '@/components/ui/select';
import { SelectValue } from '@radix-ui/react-select';
import { CodeBlock } from '@/components/common/code-block';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';


export default function Agents() {
  const isInitialLoading = useRef(false);
  const [isInitialLoaded, setIsInitialLoaded] = useState(false);
  const [selectedEnvironmentId, setSelectedEnvironmentId] = useState<
    string | null
  >(null);
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const environments = useSelector(environmentsStore.selector.getEnvironments);
  const agents = useSelector(agentsStore.selector.getAgents);
  const user = useSelector(usersStore.selector.getUserMe);
  const devKey = useSelector(apiKeysStore.selector.getDevKey);
  const apiKeys = useSelector(apiKeysStore.selector.getApiKeys);

  const selectedApiKey = useMemo(() => {
    if (!selectedEnvironmentId) {
      return null;
    }
    if (
      environments.find((e) => e.id === selectedEnvironmentId)?.slug ===
      'development'
    ) {
      return devKey;
    }
    return (
      apiKeys.find(
        (apiKey) => apiKey.environment.id === selectedEnvironmentId,
      ) ?? null
    );
  }, [apiKeys, devKey, environments, selectedEnvironmentId]);

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
    await dispatch(agentsStore.asyncActions.listAgents({ environmentId }));
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
          dispatch(apiKeysStore.asyncActions.listApiKeys()),
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
          await dispatch(agentsStore.asyncActions.listAgents({ environmentId }));
        }
        isInitialLoading.current = false;
        setIsInitialLoaded(true);
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
            {agents.length} Agents in
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
        
        {isInitialLoaded && agents.length === 0 && (
          <div className="flex w-full flex-col gap-4 rounded-md md:border md:p-6">
            <h2 className="text-xl font-bold">
              {t('routes_agents_placeholder_title')}
            </h2>
            <p className="text-sidebar-foreground font-normal">
              {t('routes_agents_placeholder_description')}
            </p>
            {!selectedApiKey?.key && (
              <Alert
                variant="destructive"
                className="border-destructive bg-destructive/10"
              >
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>No API Key Found for This Environment</AlertTitle>
                <AlertDescription>
                  <p>
                    Create an API key for this environment to access the setup
                    code and get started with your integration.
                  </p>
                  <Button variant="destructive" asChild>
                    <Link to={'/apiKeys/new'}>Edit API Key</Link>
                  </Button>
                </AlertDescription>
              </Alert>
            )}
            <CodeBlock
              code={`package main

import (
	"context"
	"log"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/agent"
	"github.com/trysourcetool/sourcetool-go/agent/models"
)

// Weather tool parameters
type WeatherParams struct {
	Location string \`json:"location" desc:"City name or coordinates" required:"true"\`
}

func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "${selectedApiKey?.key ?? 'your_api_key'}",
		Endpoint: "${user?.organization?.webSocketEndpoint}",
	})

	// Create weather tool
	weatherTool := agent.NewTool("get_weather", "Get current weather information",
		func(ctx context.Context, params WeatherParams) (interface{}, error) {
			// Your weather API logic here
			return map[string]interface{}{
				"location":    params.Location,
				"temperature": "72°F",
				"condition":   "Sunny",
				"humidity":    "45%",
			}, nil
		},
	)

	// Create and register an agent
	s.Agent("assistant", &sourcetool.Agent{
		Name:        "assistant", 
		Description: "A helpful AI assistant for general tasks",
		Instructions: "You are a helpful assistant. Be concise and accurate.",
		Model:       models.OpenAI("gpt-4o-mini"),
		Tools:       []agent.Tool{weatherTool},
	})

	if err := s.Listen(); err != nil {
		log.Fatal(err)
	}
}`}
              language="go"
            />

            <p className="text-sidebar-foreground font-normal">
              {t('routes_agents_placeholder_restart_server')}
            </p>
            <p className="text-sidebar-foreground font-normal">
              {t('routes_agents_placeholder_agent_added')}
            </p>
            <p className="text-sidebar-foreground font-normal">
              {t('routes_agents_placeholder_documentation')}
            </p>
          </div>
        )}
        
        {agents.length > 0 && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {agents.map((agent) => (
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
                        <div className="mt-1 text-sm text-muted-foreground">
                          Model: {agent.model}
                        </div>
                      </div>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <CardDescription className="mb-4">
                    {agent.description || 'No description available'}
                  </CardDescription>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
        )}
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_auth/agents/')({
  component: Agents,
});