import * as users from './modules/users';
import * as pages from './modules/pages';
import * as organizations from './modules/organizations';
import * as environments from './modules/environments';
import * as hostInstances from './modules/hostInstances';
import * as apiKeys from './modules/apiKeys';
import * as auth from './modules/auth';
import * as groups from './modules/groups';
import { api as apiInstance } from './instance';

export const api = {
  apiKeys,
  auth,
  users,
  pages,
  organizations,
  environments,
  hostInstances,
  groups,
  setExpiresAt: (expiresAt: string) => {
    apiInstance.setExpiresAt(expiresAt);
  },
};
