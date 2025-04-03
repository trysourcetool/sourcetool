import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';
import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';
import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import type { ApiKey } from '@/api/modules/apiKeys';

// =============================================
// asyncActions
// =============================================
const listApiKeys = createAsyncThunk(
  'apiKeys/listApiKeys',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.apiKeys.listApiKeys();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const getApiKey = createAsyncThunk(
  'apiKeys/getApiKey',
  async (params: { apiKeyId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.apiKeys.getApiKey(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const createApiKey = createAsyncThunk(
  'apiKeys/createApiKey',
  async (
    params: { environmentId: string; name: string },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.apiKeys.createApiKey({
        data: {
          environmentId: params.environmentId,
          name: params.name,
        },
      });

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateApiKey = createAsyncThunk(
  'apiKeys/updateApiKey',
  async (
    params: {
      apiKeyId: string;
      data: {
        name: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.apiKeys.updateApiKey(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const deleteApiKey = createAsyncThunk(
  'apiKeys/deleteApiKey',
  async (params: { apiKeyId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.apiKeys.deleteApiKey(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

// =============================================
// slice
// =============================================
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
      .addCase(listApiKeys.pending, (state) => {
        state.isListApiKeysWaiting = true;
      })
      .addCase(listApiKeys.fulfilled, (state, action) => {
        state.isListApiKeysWaiting = false;
        state.devKey = action.payload.devKey;
        apiKeysAdapter.setAll(state.liveKeys, action.payload.liveKeys);
      })
      .addCase(listApiKeys.rejected, (state) => {
        state.isListApiKeysWaiting = false;
      })

      // getApiKey
      .addCase(getApiKey.pending, (state) => {
        state.isGetApiKeyWaiting = true;
      })
      .addCase(getApiKey.fulfilled, (state, action) => {
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
      .addCase(getApiKey.rejected, (state) => {
        state.isGetApiKeyWaiting = false;
      })

      // createApiKey
      .addCase(createApiKey.pending, (state) => {
        state.isCreateApiKeyWaiting = true;
      })
      .addCase(createApiKey.fulfilled, (state, action) => {
        state.isCreateApiKeyWaiting = false;
        apiKeysAdapter.addOne(state.liveKeys, action.payload.apiKey);
      })
      .addCase(createApiKey.rejected, (state) => {
        state.isCreateApiKeyWaiting = false;
      })

      // updateApiKey
      .addCase(updateApiKey.pending, (state) => {
        state.isUpdateApiKeyWaiting = true;
      })
      .addCase(updateApiKey.fulfilled, (state, action) => {
        state.isUpdateApiKeyWaiting = false;
        apiKeysAdapter.updateOne(state.liveKeys, {
          id: action.payload.apiKey.id,
          changes: action.payload.apiKey,
        });
      })
      .addCase(updateApiKey.rejected, (state) => {
        state.isUpdateApiKeyWaiting = false;
      })

      // deleteApiKey
      .addCase(deleteApiKey.pending, (state) => {
        state.isDeleteApiKeyWaiting = true;
      })
      .addCase(deleteApiKey.fulfilled, (state, action) => {
        state.isDeleteApiKeyWaiting = false;
        apiKeysAdapter.removeOne(state.liveKeys, action.payload.apiKey.id);
      });
  },
  initialState,
  name: 'apiKeys',
  reducers: {},
});

// =============================================
// selectors
// =============================================
const getApiKeysIds = createSelector(
  (state: RootState) => state.apiKeys,
  (values) => values.liveKeys.ids,
);

const getApiKeysEntities = createSelector(
  (state: RootState) => state.apiKeys,
  (values) => values.liveKeys.entities,
);

const getApiKeys = createSelector(
  (state: RootState) => state.apiKeys,
  (values) => values.liveKeys.ids.map((id) => values.liveKeys.entities[id]),
);

const getApiKeyValue = createSelector(
  (state: RootState, apiKeyId: string) =>
    state.apiKeys.liveKeys.entities[apiKeyId],
  (values) => values,
);

const getDevKey = createSelector(
  (state: RootState) => state.apiKeys.devKey,
  (values) => values,
);

// =============================================
// exports
// =============================================

export const apiKeysStore = {
  actions: slice.actions,
  asyncActions: {
    listApiKeys,
    getApiKey,
    createApiKey,
    updateApiKey,
    deleteApiKey,
  },
  reducer: slice.reducer,
  selector: {
    getApiKeysIds,
    getApiKeysEntities,
    getApiKeys,
    getApiKey: getApiKeyValue,
    getDevKey,
  },
};

export type ApiKeysState = State;
