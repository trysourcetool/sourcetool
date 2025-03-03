import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';

export const getHostInstancePing = createAsyncThunk(
  'hostInstances/getHostInstancePing',
  async (params: { pageId?: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.hostInstances.getHostInstancePing(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);
