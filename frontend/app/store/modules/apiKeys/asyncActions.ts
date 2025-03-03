import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';

export const listApiKeys = createAsyncThunk(
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

export const getApiKey = createAsyncThunk(
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

export const createApiKey = createAsyncThunk(
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

export const updateApiKey = createAsyncThunk(
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

export const deleteApiKey = createAsyncThunk(
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
