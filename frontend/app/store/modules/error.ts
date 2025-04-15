import { createAsyncThunk } from '@reduxjs/toolkit';
import { createSlice } from '@reduxjs/toolkit';

// =============================================
// asyncActions
// =============================================
const handleError = createAsyncThunk('error/handleError', async () => {});

// =============================================
//  slice
// =============================================
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
      .addCase(handleError.pending, () => {})
      .addCase(handleError.fulfilled, () => {})
      .addCase(handleError.rejected, () => {});
  },
  initialState,
  name: 'error',
  reducers: {},
});

// =============================================
// exports
// =============================================

export const errorStore = {
  actions: slice.actions,
  asyncActions: {
    handleError,
  },
  reducer: slice.reducer,
};

export type ErrorState = State;
