import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';

export const listGroups = createAsyncThunk(
  'groups/listGroups',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.groups.listGroups();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const getGroup = createAsyncThunk(
  'groups/getGroup',
  async (params: { groupId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.groups.getGroup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const createGroup = createAsyncThunk(
  'groups/createGroup',
  async (
    params: { data: { name: string; slug: string; userIds: string[] } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.groups.createGroup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const updateGroup = createAsyncThunk(
  'groups/updateGroup',
  async (
    params: { groupId: string; data: { name: string; userIds: string[] } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.groups.updateGroup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const deleteGroup = createAsyncThunk(
  'groups/deleteGroup',
  async (params: { groupId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.groups.deleteGroup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);
