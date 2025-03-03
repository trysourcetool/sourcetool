import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';

export const getGroupIds = createSelector(
  (state: RootState) => state.groups,
  (values) => values.groups.ids,
);

export const getGroupEntities = createSelector(
  (state: RootState) => state.groups,
  (values) => values.groups.entities,
);

export const getGroups = createSelector(
  (state: RootState) => state.groups,
  (values) => values.groups.ids.map((id) => values.groups.entities[id]),
);

export const getGroup = createSelector(
  (state: RootState, groupId: string) => state.groups.groups.entities[groupId],
  (values) => values,
);

export const getUserGroupIds = createSelector(
  (state: RootState) => state.groups,
  (values) => values.userGroups.ids,
);

export const getUserGroupEntities = createSelector(
  (state: RootState) => state.groups,
  (values) => values.userGroups.entities,
);

export const getUserGroups = createSelector(
  (state: RootState) => state.groups,
  (values) => values.userGroups.ids.map((id) => values.userGroups.entities[id]),
);
