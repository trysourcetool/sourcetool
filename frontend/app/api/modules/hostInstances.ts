import * as api from '@/api/instance';

export type HostInstance = {
  createdAt: string;
  id: string;
  sdkName: string;
  sdkVersion: string;
  status: string;
  updatedAt: string;
};

export const getHostInstancePing = async (params?: { pageId?: string }) => {
  const res = await api.get<{
    hostInstance: HostInstance;
  }>({
    path: '/hostInstances/ping',
    auth: true,
    params,
  });

  return res;
};
