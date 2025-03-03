import * as api from '@/api/instance';

export type Environment = {
  color: string;
  createdAt: string;
  id: string;
  name: string;
  slug: string;
  updatedAt: string;
};

export const listEnvironments = async () => {
  const res = await api.get<{
    environments: Environment[];
  }>({ path: '/environments', auth: true });

  return res;
};

export const getEnvironment = async (params: { environmentId: string }) => {
  const res = await api.get<{
    environment: Environment;
  }>({ path: `/environments/${params.environmentId}`, auth: true });

  return res;
};

export const createEnvironment = async (params: {
  data: {
    color: string;
    name: string;
    slug: string;
  };
}) => {
  const res = await api.post<{
    environment: Environment;
  }>({ path: '/environments', data: params.data, auth: true });

  return res;
};

export const updateEnvironment = async (params: {
  environmentId: string;
  data: {
    color: string;
    name: string;
  };
}) => {
  const res = await api.put<{
    environment: Environment;
  }>({
    path: `/environments/${params.environmentId}`,
    data: params.data,
    auth: true,
  });

  return res;
};

export const deleteEnvironment = async (params: { environmentId: string }) => {
  const res = await api.del<{
    environment: Environment;
  }>({
    path: `/environments/${params.environmentId}`,
    auth: true,
  });

  return res;
};
