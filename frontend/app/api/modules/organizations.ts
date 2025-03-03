import * as api from '@/api/instance';
import type { User, UserRole } from './users';

export type Organization = {
  createdAt: string;
  id: string;
  subdomain: string;
  updatedAt: string;
};

export const createOrganization = async (params: {
  data: {
    subdomain: string;
  };
}) => {
  const res = await api.post<{
    organization: Organization;
  }>({ path: '/organizations', data: params.data, auth: true });

  return res;
};

export const checkSubdomainAvailability = async (params: {
  subdomain: string;
}) => {
  const res = await api.get<{
    available: boolean;
  }>({
    path: '/organizations/checkSubdomainAvailability',
    params: params,
    auth: true,
  });

  return res;
};

export const updateOrganizationUser = async (params: {
  userId: string;
  data: {
    groupIds: string[];
    role: UserRole;
  };
}) => {
  const res = await api.put<User>({
    path: `/organizations/users/${params.userId}`,
    data: params.data,
    auth: true,
  });

  return res;
};
