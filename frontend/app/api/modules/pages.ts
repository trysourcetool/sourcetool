import * as api from '@/api/instance';
import type { Group, GroupPage, UserGroup } from './groups';
import type { User } from './users';

export type Page = {
  createdAt: string;
  id: string;
  name: string;
  updatedAt: string;
  route: string;
};

export const listPages = async () => {
  const res = await api.get<{
    pages: Page[];
    groupPages: GroupPage[];
    groups: Group[];
    userGroups: UserGroup[];
    users: User[];
  }>({ path: '/pages', auth: true });

  return res;
};
