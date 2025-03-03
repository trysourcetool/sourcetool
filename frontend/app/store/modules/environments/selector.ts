import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';

export const getEnvironmentIds = createSelector(
  (state: RootState) => state.environments,
  (values) => values.environments.ids,
);

export const getEnvironmentEntities = createSelector(
  (state: RootState) => state.environments,
  (values) => values.environments.entities,
);

export const getEnvironments = createSelector(
  (state: RootState) => state.environments,
  (values) =>
    values.environments.ids.map((id) => values.environments.entities[id]),
);

export const getEnvironment = createSelector(
  (state: RootState, environmentId: string) =>
    state.environments.environments.entities[environmentId],
  (values) => values,
);
