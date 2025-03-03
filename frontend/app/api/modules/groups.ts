import * as api from '@/api/instance';
import type { User } from './users';

export type Group = {
  createdAt: string;
  id: string;
  name: string;
  slug: string;
  updatedAt: string;
};

export type UserGroup = {
  createdAt: string;
  groupId: string;
  id: string;
  updatedAt: string;
  userId: string;
};

export type GroupPage = {
  createdAt: string;
  groupId: string;
  id: string;
  pageId: string;
  updatedAt: string;
};

export const listGroups = async () => {
  const res = await api.get<{
    groups: Group[];
    userGroups: UserGroup[];
    users: User[];
  }>({ path: '/groups', auth: true });

  return res;
};

export const getGroup = async (params: { groupId: string }) => {
  const res = await api.get<{
    group: Group;
  }>({ path: `/groups/${params.groupId}`, auth: true });

  return res;
};

export const createGroup = async (params: {
  data: {
    name: string;
    slug: string;
    userIds: string[];
  };
}) => {
  const res = await api.post<{
    group: Group;
  }>({ path: '/groups', data: params.data, auth: true });

  return res;
};

export const updateGroup = async (params: {
  groupId: string;
  data: {
    name: string;
    userIds: string[];
  };
}) => {
  const res = await api.put<{
    group: Group;
  }>({
    path: `/groups/${params.groupId}`,
    data: params.data,
    auth: true,
  });

  return res;
};

export const deleteGroup = async (params: { groupId: string }) => {
  const res = await api.del<{
    group: Group;
  }>({
    path: `/groups/${params.groupId}`,
    auth: true,
  });

  return res;
};
