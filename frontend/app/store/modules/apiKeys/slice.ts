import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';
import type { ApiKey } from '@/api/modules/apiKeys';

// =============================================
// schema

const apiKeysAdapter = createEntityAdapter<ApiKey, string>({
  selectId: (apiKey) => apiKey.id,
});

// =============================================
// State

export type State = {
  devKey: null | ApiKey;
  liveKeys: EntityState<ApiKey, string>;
  isListApiKeysWaiting: boolean;
  isCreateApiKeyWaiting: boolean;
  isDeleteApiKeyWaiting: boolean;
  isGetApiKeyWaiting: boolean;
  isUpdateApiKeyWaiting: boolean;
};

const initialState: State = {
  devKey: null,
  liveKeys: apiKeysAdapter.getInitialState(),
  isListApiKeysWaiting: false,
  isCreateApiKeyWaiting: false,
  isDeleteApiKeyWaiting: false,
  isGetApiKeyWaiting: false,
  isUpdateApiKeyWaiting: false,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // listApiKeys
      .addCase(asyncActions.listApiKeys.pending, (state) => {
        state.isListApiKeysWaiting = true;
      })
      .addCase(asyncActions.listApiKeys.fulfilled, (state, action) => {
        state.isListApiKeysWaiting = false;
        state.devKey = action.payload.devKey;
        apiKeysAdapter.setAll(state.liveKeys, action.payload.liveKeys);
      })
      .addCase(asyncActions.listApiKeys.rejected, (state) => {
        state.isListApiKeysWaiting = false;
      })

      // getApiKey
      .addCase(asyncActions.getApiKey.pending, (state) => {
        state.isGetApiKeyWaiting = true;
      })
      .addCase(asyncActions.getApiKey.fulfilled, (state, action) => {
        state.isGetApiKeyWaiting = false;

        if (state.liveKeys.entities[action.payload.apiKey.id]) {
          apiKeysAdapter.updateOne(state.liveKeys, {
            id: action.payload.apiKey.id,
            changes: action.payload.apiKey,
          });
        } else {
          apiKeysAdapter.addOne(state.liveKeys, action.payload.apiKey);
        }
      })
      .addCase(asyncActions.getApiKey.rejected, (state) => {
        state.isGetApiKeyWaiting = false;
      })

      // createApiKey
      .addCase(asyncActions.createApiKey.pending, (state) => {
        state.isCreateApiKeyWaiting = true;
      })
      .addCase(asyncActions.createApiKey.fulfilled, (state, action) => {
        state.isCreateApiKeyWaiting = false;
        apiKeysAdapter.addOne(state.liveKeys, action.payload.apiKey);
      })
      .addCase(asyncActions.createApiKey.rejected, (state) => {
        state.isCreateApiKeyWaiting = false;
      })

      // updateApiKey
      .addCase(asyncActions.updateApiKey.pending, (state) => {
        state.isUpdateApiKeyWaiting = true;
      })
      .addCase(asyncActions.updateApiKey.fulfilled, (state, action) => {
        state.isUpdateApiKeyWaiting = false;
        apiKeysAdapter.updateOne(state.liveKeys, {
          id: action.payload.apiKey.id,
          changes: action.payload.apiKey,
        });
      })
      .addCase(asyncActions.updateApiKey.rejected, (state) => {
        state.isUpdateApiKeyWaiting = false;
      })

      // deleteApiKey
      .addCase(asyncActions.deleteApiKey.pending, (state) => {
        state.isDeleteApiKeyWaiting = true;
      })
      .addCase(asyncActions.deleteApiKey.fulfilled, (state, action) => {
        state.isDeleteApiKeyWaiting = false;
        apiKeysAdapter.removeOne(state.liveKeys, action.payload.apiKey.id);
      });
  },
  initialState,
  name: 'apiKeys',
  reducers: {},
});
