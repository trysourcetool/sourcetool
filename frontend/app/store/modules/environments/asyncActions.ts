import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';

export const listEnvironments = createAsyncThunk(
  'environments/listEnvironments',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.environments.listEnvironments();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const getEnvironment = createAsyncThunk(
  'environments/getEnvironment',
  async (params: { environmentId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.environments.getEnvironment(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const createEnvironment = createAsyncThunk(
  'environments/createEnvironment',
  async (
    params: { data: { color: string; name: string; slug: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.environments.createEnvironment(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const updateEnvironment = createAsyncThunk(
  'environments/updateEnvironment',
  async (
    params: { environmentId: string; data: { color: string; name: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.environments.updateEnvironment(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const deleteEnvironment = createAsyncThunk(
  'environments/deleteEnvironment',
  async (params: { environmentId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.environments.deleteEnvironment(params);

      return res;
    } catch (error: any) {
      console.error({ error });
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);
