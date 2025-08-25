import { useEffect } from 'react';
import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '@/components/common/page-header';
import { AgentChat } from '@/components/agent/agent-chat';
import { useDispatch, useSelector } from '@/store';
import { agentsStore } from '@/store/modules/agents';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useTranslation } from 'react-i18next';
import { Bot } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export default function AgentDetail() {
  const { agentId } = Route.useParams();
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  
  const agent = useSelector((state) => 
    agentsStore.selector.getAgent(state, agentId)
  );

  useEffect(() => {
    if (agentId) {
      // Fetch agent details if not already loaded
      dispatch(agentsStore.asyncActions.getAgent({ agentId }));
    }
  }, [agentId, dispatch]);

  useEffect(() => {
    if (setBreadcrumbsState && agent) {
      setBreadcrumbsState([
        { label: t('breadcrumbs_agents'), to: '/agents' },
        { label: agent.name, to: `/agents/${agentId}` }
      ]);
    }
  }, [setBreadcrumbsState, t, agent, agentId]);

  if (!agent) {
    return (
      <div className="flex h-full flex-col">
        <PageHeader label={t('routes_agents_loading')} />
        <div className="flex items-center justify-center p-8">
          <div className="text-muted-foreground">Loading agent...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      <PageHeader label={agent.name} />
      
      <div className="flex flex-1 flex-col lg:flex-row">
        {/* Agent Info Sidebar */}
        <div className="border-r lg:w-80">
          <div className="p-4 md:p-6">
            <Card>
              <CardHeader>
                <div className="flex items-center gap-3">
                  <div className="rounded-full bg-primary/10 p-3">
                    <Bot className="h-6 w-6 text-primary" />
                  </div>
                  <div>
                    <CardTitle>{agent.name}</CardTitle>
                    <div className="text-sm text-muted-foreground">
                      Model: {agent.model}
                    </div>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <CardDescription className="mb-4">
                  {agent.description || 'No description available'}
                </CardDescription>
                
                {agent.instructions && (
                  <div className="space-y-2">
                    <h4 className="text-sm font-medium">Instructions</h4>
                    <p className="text-sm text-muted-foreground">
                      {agent.instructions}
                    </p>
                  </div>
                )}

                {agent.tools && agent.tools.length > 0 && (
                  <div className="mt-4 space-y-2">
                    <h4 className="text-sm font-medium">Available Tools</h4>
                    <div className="space-y-1">
                      {agent.tools.map((tool: any, index: number) => (
                        <div key={index} className="text-sm">
                          <span className="font-medium">{tool.name}:</span>{' '}
                          <span className="text-muted-foreground">{tool.description}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Chat Area */}
        <div className="flex-1">
          <AgentChat agentId={agentId} agent={agent} />
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_auth/agents/$agentId')({
  component: AgentDetail,
});