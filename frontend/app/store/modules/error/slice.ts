import { createSlice } from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';

// =============================================
// State

export type State = {};

const initialState: State = {};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // handleError
      .addCase(asyncActions.handleError.pending, () => {})
      .addCase(asyncActions.handleError.fulfilled, () => {})
      .addCase(asyncActions.handleError.rejected, () => {});
  },
  initialState,
  name: 'error',
  reducers: {},
});
