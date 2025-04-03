import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../';
import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from './error';
import type { ErrorResponse } from '@/api/instance';
import {
  createEntityAdapter,
  createSlice,
  type EntityState,
  type PayloadAction,
} from '@reduxjs/toolkit';
import type { Page } from '@/api/modules/pages';
import type { ExceptionJson } from '@trysourcetool/proto/exception/v1/exception';

// =============================================
// asyncActions
// =============================================

export const listPages = createAsyncThunk(
  'pages/listPages',
  async (params: { environmentId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.pages.listPages(params);

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

const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // listPages
      .addCase(listPages.pending, (state) => {
        state.isListPagesWaiting = true;
      })
      .addCase(listPages.fulfilled, (state, action) => {
        state.isListPagesWaiting = false;
        pagesAdapter.setAll(state.pages, action.payload.pages);
      })
      .addCase(listPages.rejected, (state) => {
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

// =============================================
// selectors
// =============================================

const getPageIds = createSelector(
  (state: RootState) => state.pages,
  (values) => values.pages.ids,
);

const getPageEntities = createSelector(
  (state: RootState) => state.pages,
  (values) => values.pages.entities,
);

const getPages = createSelector(
  (state: RootState) => state.pages,
  (values) => values.pages.ids.map((id) => values.pages.entities[id]),
);

const getPage = createSelector(
  (state: RootState, pageId: string) => state.pages.pages.entities[pageId],
  (values) => values || null,
);

const getPageFromPath = createSelector(
  (state: RootState, path: string) => ({
    pages: getPages(state),
    path,
  }),
  ({ pages, path }) => {
    const page = pages.find((page) => page.route === path);
    console.log({ page });
    return page || null;
  },
);

const getPermissionPages = createSelector(
  (state: RootState) => ({
    account: state.users.me,
    pages: state.pages,
    groups: state.groups,
  }),
  ({ account, pages, groups }) => {
    const userGroups = groups.userGroups.ids
      .map((id) => groups.userGroups.entities[id])
      .filter((userGroup) => userGroup.userId === account?.id);

    const groupPages = groups.groupPages.ids
      .map((id) => groups.groupPages.entities[id])
      .filter((groupPage) =>
        userGroups.some((userGroup) => userGroup.groupId === groupPage.groupId),
      );

    console.log(
      { groupPages, userGroups },
      pages.pages,
      pages.pages.ids
        .filter((id) => groupPages.some((page) => page.pageId === id))
        .map((id) => pages.pages.entities[id]),
    );

    return pages.pages.ids
      .filter(
        (id) =>
          !groups.groups.ids.length ||
          groupPages.some((page) => page.pageId === id),
      )
      .map((id) => pages.pages.entities[id]);
  },
);

// =============================================
// exports
// =============================================
export const pagesStore = {
  actions: slice.actions,
  asyncActions: {
    listPages,
  },
  reducer: slice.reducer,
  selector: {
    getPageIds,
    getPageEntities,
    getPages,
    getPage,
    getPageFromPath,
    getPermissionPages,
  },
};

export type PagesState = State;
