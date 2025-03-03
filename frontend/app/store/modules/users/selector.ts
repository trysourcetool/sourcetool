import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';

export const getMe = createSelector(
  (state: RootState) => state.users,
  (values) => values.me,
);

export const getUserIds = createSelector(
  (state: RootState) => state.users,
  (values) => values.users.ids,
);

export const getUserEntities = createSelector(
  (state: RootState) => state.users,
  (values) => values.users.entities,
);

export const getUsers = createSelector(
  (state: RootState) => state.users,
  (values) => values.users.ids.map((id) => values.users.entities[id]),
);

export const getUser = createSelector(
  (state: RootState, userId: string) => ({
    users: state.users.users,
    userId,
  }),
  ({ users, userId }) => users.entities[userId],
);

export const getSubDomainMatched = createSelector(
  (state: RootState, subDomain: string | null) => {
    const isAuthChecked =
      state.users.isAuthChecked &&
      (state.users.isAuthFailed ||
        (state.users.isAuthSucceeded && state.users.me));
    const matched = state.users.me?.organization?.subdomain === subDomain;
    return {
      isMatched: matched,
      status: !isAuthChecked ? 'checking' : 'checked',
    } as {
      isMatched: boolean;
      status: 'checking' | 'checked';
    };
  },
  (values) => values,
);
