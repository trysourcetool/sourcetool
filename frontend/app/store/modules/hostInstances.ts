import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from './error';
import type { ErrorResponse } from '@/api/instance';
import { createSlice } from '@reduxjs/toolkit';
import type { HostInstance } from '@/api/modules/hostInstances';

// =============================================
// asyncActions
// =============================================
const getHostInstancePing = createAsyncThunk(
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

// =============================================
// slice
// =============================================

// =============================================
// schema

// =============================================
// State

export type State = {
  isHostInstancePingError: boolean;
  isGetHostInstancePingWaiting: boolean;
  hostInstance: HostInstance | null;
};

const initialState: State = {
  isHostInstancePingError: false,
  isGetHostInstancePingWaiting: false,
  hostInstance: null,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // getHostInstancePing
      .addCase(getHostInstancePing.pending, (state) => {
        state.isGetHostInstancePingWaiting = true;
      })
      .addCase(getHostInstancePing.fulfilled, (state, action) => {
        state.isGetHostInstancePingWaiting = false;
        state.isHostInstancePingError = false;
        state.hostInstance = action.payload.hostInstance;
      })
      .addCase(getHostInstancePing.rejected, (state) => {
        state.isGetHostInstancePingWaiting = false;
        state.isHostInstancePingError = true;
        state.hostInstance = null;
      });
  },
  initialState,
  name: 'hostInstances',
  reducers: {},
});

// =============================================
// exports
// =============================================

export const hostInstancesStore = {
  actions: slice.actions,
  asyncActions: {
    getHostInstancePing,
  },
  reducer: slice.reducer,
};

export type HostInstancesState = State;
