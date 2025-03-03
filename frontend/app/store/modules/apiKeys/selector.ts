import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';

export const getApiKeysIds = createSelector(
  (state: RootState) => state.apiKeys,
  (values) => values.liveKeys.ids,
);

export const getApiKeysEntities = createSelector(
  (state: RootState) => state.apiKeys,
  (values) => values.liveKeys.entities,
);

export const getApiKeys = createSelector(
  (state: RootState) => state.apiKeys,
  (values) => values.liveKeys.ids.map((id) => values.liveKeys.entities[id]),
);

export const getApiKey = createSelector(
  (state: RootState, apiKeyId: string) =>
    state.apiKeys.liveKeys.entities[apiKeyId],
  (values) => values,
);

export const getDevKey = createSelector(
  (state: RootState) => state.apiKeys.devKey,
  (values) => values,
);
