import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';

export const listPages = createAsyncThunk(
  'pages/listPages',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.pages.listPages();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);
