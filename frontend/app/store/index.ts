import {
  combineReducers,
  configureStore as createConfigureStore,
} from '@reduxjs/toolkit';
import {
  useSelector as useAppSelector,
  useDispatch as useAppDispatch,
  type TypedUseSelectorHook,
} from 'react-redux';
import { ENVIRONMENTS } from '@/environments';
import { type UsersState, usersStore } from './modules/users';
import {
  organizationsStore,
  type OrganizationsState,
} from './modules/organizations';
import {
  environmentsStore,
  type EnvironmentsState,
} from './modules/environments';
import { pagesStore, type PagesState } from './modules/pages';
import { apiKeysStore, type ApiKeysState } from './modules/apiKeys';
import { widgetsStore, type WidgetsState } from './modules/widgets';
import { groupsStore, type GroupsState } from './modules/groups';
import {
  hostInstancesStore,
  type HostInstancesState,
} from './modules/hostInstances';
import { type AuthState, authStore } from './modules/auth';

export type RootState = {
  users: UsersState;
  organizations: OrganizationsState;
  environments: EnvironmentsState;
  pages: PagesState;
  widgets: WidgetsState;
  apiKeys: ApiKeysState;
  groups: GroupsState;
  hostInstances: HostInstancesState;
  auth: AuthState;
};

export const configureStore = () => {
  const rootReducer = combineReducers({
    users: usersStore.reducer,
    organizations: organizationsStore.reducer,
    environments: environmentsStore.reducer,
    pages: pagesStore.reducer,
    apiKeys: apiKeysStore.reducer,
    widgets: widgetsStore.reducer,
    groups: groupsStore.reducer,
    hostInstances: hostInstancesStore.reducer,
    auth: authStore.reducer,
  });

  const store = createConfigureStore({
    devTools:
      ENVIRONMENTS.MODE !== 'production'
        ? {
            traceLimit: 100,
          }
        : false,
    reducer: rootReducer,
  });

  return {
    store,
  };
};

export const useSelector: TypedUseSelectorHook<RootState> = useAppSelector;
export type AppDispatch = ReturnType<
  typeof configureStore
>['store']['dispatch'];

export const useDispatch = () => useAppDispatch<AppDispatch>();
