import {
  createEntityAdapter,
  createSlice,
  type EntityState,
  type PayloadAction,
} from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';
import type { Page } from '@/api/modules/pages';
import type { ExceptionJson } from '@trysourcetool/proto/exception/v1/exception';

// =============================================
// schema

const pagesAdapter = createEntityAdapter<Page, string>({
  selectId: (page) => page.id,
});

// =============================================
// State

export type State = {
  pages: EntityState<Page, string>;
  exception: ExceptionJson | null;
  isListPagesWaiting: boolean;
};

const initialState: State = {
  exception: null,
  pages: pagesAdapter.getInitialState(),
  isListPagesWaiting: false,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // listPages
      .addCase(asyncActions.listPages.pending, (state) => {
        state.isListPagesWaiting = true;
      })
      .addCase(asyncActions.listPages.fulfilled, (state, action) => {
        state.isListPagesWaiting = false;
        pagesAdapter.setAll(state.pages, action.payload.pages);
      })
      .addCase(asyncActions.listPages.rejected, (state) => {
        state.isListPagesWaiting = false;
      });
  },
  initialState,
  name: 'pages',
  reducers: {
    setException: (state, action: PayloadAction<ExceptionJson>) => {
      state.exception = action.payload;
    },
    clearException: (state) => {
      state.exception = null;
    },
  },
});
