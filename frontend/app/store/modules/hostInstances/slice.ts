import { createSlice } from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';
import type { HostInstance } from '@/api/modules/hostInstances';

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
      .addCase(asyncActions.getHostInstancePing.pending, (state) => {
        state.isGetHostInstancePingWaiting = true;
      })
      .addCase(asyncActions.getHostInstancePing.fulfilled, (state, action) => {
        state.isGetHostInstancePingWaiting = false;
        state.isHostInstancePingError = false;
        state.hostInstance = action.payload.hostInstance;
      })
      .addCase(asyncActions.getHostInstancePing.rejected, (state) => {
        state.isGetHostInstancePingWaiting = false;
        state.isHostInstancePingError = true;
        state.hostInstance = null;
      });
  },
  initialState,
  name: 'hostInstances',
  reducers: {},
});
