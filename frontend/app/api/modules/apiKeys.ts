import * as api from '@/api/instance';
import type { Environment } from './environments';

export type ApiKey = {
  createdAt: string;
  id: string;
  key: string;
  name: string;
  updatedAt: string;
  environment: Environment;
};

export const listApiKeys = async () => {
  const res = await api.get<{
    devKey: ApiKey;
    liveKeys: ApiKey[];
  }>({ path: '/apiKeys', auth: true });

  return res;
};

export const getApiKey = async (params: { apiKeyId: string }) => {
  const res = await api.get<{
    apiKey: ApiKey;
  }>({ path: `/apiKeys/${params.apiKeyId}`, auth: true });

  return res;
};

export const createApiKey = async (params: {
  data: {
    environmentId: string;
    name: string;
  };
}) => {
  const res = await api.post<{
    apiKey: ApiKey;
  }>({ path: '/apiKeys', data: params.data, auth: true });

  return res;
};

export const updateApiKey = async (params: {
  apiKeyId: string;
  data: {
    name: string;
  };
}) => {
  const res = await api.put<{
    apiKey: ApiKey;
  }>({
    path: `/apiKeys/${params.apiKeyId}`,
    data: params.data,
    auth: true,
  });

  return res;
};

export const deleteApiKey = async (params: { apiKeyId: string }) => {
  const res = await api.del<{
    apiKey: ApiKey;
  }>({
    path: `/apiKeys/${params.apiKeyId}`,
    auth: true,
  });

  return res;
};
