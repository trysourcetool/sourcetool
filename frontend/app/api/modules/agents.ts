import * as api from '@/api/instance';

export type Agent = {
  id: string;
  name: string;
  description: string;
  instructions: string;
  model: string;
};

export const listAgents = async ({ environmentId }: { environmentId: string }) => {
  const res = await api.get<{
    agents: Agent[];
  }>({ path: `/agents?environmentId=${environmentId}`, auth: true });

  return res;
};

export const getAgent = async ({ agentId }: { agentId: string }) => {
  const res = await api.get<Agent>({ path: `/agents/${agentId}`, auth: true });

  return { agent: res };
};