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
} from '@reduxjs/toolkit';
import type { Group, GroupPage, UserGroup } from '@/api/modules/groups';
import { pagesStore } from './pages';

// =============================================
// asyncActions
// =============================================
const listGroups = createAsyncThunk(
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

const getGroup = createAsyncThunk(
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

const createGroup = createAsyncThunk(
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

const updateGroup = createAsyncThunk(
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

const deleteGroup = createAsyncThunk(
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

// =============================================
// slice
// =============================================
// =============================================
// schema

const groupsAdapter = createEntityAdapter<Group, string>({
  selectId: (group) => group.id,
});

const userGroupsAdapter = createEntityAdapter<UserGroup, string>({
  selectId: (userGroup) => userGroup.id,
});

const groupPagesAdapter = createEntityAdapter<GroupPage, string>({
  selectId: (groupPage) => groupPage.id,
});

// =============================================
// State

export type State = {
  groups: EntityState<Group, string>;
  userGroups: EntityState<UserGroup, string>;
  groupPages: EntityState<GroupPage, string>;
  isListGroupsWaiting: boolean;
  isGetGroupWaiting: boolean;
  isCreateGroupWaiting: boolean;
  isUpdateGroupWaiting: boolean;
  isDeleteGroupWaiting: boolean;
};

const initialState: State = {
  groups: groupsAdapter.getInitialState(),
  userGroups: userGroupsAdapter.getInitialState(),
  groupPages: groupPagesAdapter.getInitialState(),
  isListGroupsWaiting: false,
  isGetGroupWaiting: false,
  isCreateGroupWaiting: false,
  isUpdateGroupWaiting: false,
  isDeleteGroupWaiting: false,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // listGroups
      .addCase(listGroups.pending, (state) => {
        state.isListGroupsWaiting = true;
      })
      .addCase(listGroups.fulfilled, (state, action) => {
        state.isListGroupsWaiting = false;
        groupsAdapter.setAll(state.groups, action.payload.groups);
        userGroupsAdapter.setAll(state.userGroups, action.payload.userGroups);
      })
      .addCase(listGroups.rejected, (state) => {
        state.isListGroupsWaiting = false;
      })

      // getGroup
      .addCase(getGroup.pending, (state) => {
        state.isGetGroupWaiting = true;
      })
      .addCase(getGroup.fulfilled, (state, action) => {
        state.isGetGroupWaiting = false;

        if (state.groups.entities[action.payload.group.id]) {
          groupsAdapter.updateOne(state.groups, {
            id: action.payload.group.id,
            changes: action.payload.group,
          });
        } else {
          groupsAdapter.addOne(state.groups, action.payload.group);
        }
      })
      .addCase(getGroup.rejected, (state) => {
        state.isGetGroupWaiting = false;
      })

      // createEnvironment
      .addCase(createGroup.pending, (state) => {
        state.isCreateGroupWaiting = true;
      })
      .addCase(createGroup.fulfilled, (state, action) => {
        state.isCreateGroupWaiting = false;
        groupsAdapter.addOne(state.groups, action.payload.group);
      })
      .addCase(createGroup.rejected, (state) => {
        state.isCreateGroupWaiting = false;
      })

      // updateEnvironment
      .addCase(updateGroup.pending, (state) => {
        state.isUpdateGroupWaiting = true;
      })
      .addCase(updateGroup.fulfilled, (state, action) => {
        state.isUpdateGroupWaiting = false;
        groupsAdapter.updateOne(state.groups, {
          id: action.payload.group.id,
          changes: action.payload.group,
        });
      })
      .addCase(updateGroup.rejected, (state) => {
        state.isUpdateGroupWaiting = false;
      })

      // deleteEnvironment
      .addCase(deleteGroup.pending, (state) => {
        state.isDeleteGroupWaiting = true;
      })
      .addCase(deleteGroup.fulfilled, (state, action) => {
        state.isDeleteGroupWaiting = false;
        groupsAdapter.removeOne(state.groups, action.payload.group.id);
      })
      .addCase(deleteGroup.rejected, (state) => {
        state.isDeleteGroupWaiting = false;
      })

      // listPages
      .addCase(pagesStore.asyncActions.listPages.pending, () => {})
      .addCase(pagesStore.asyncActions.listPages.fulfilled, (state, action) => {
        groupPagesAdapter.setAll(state.groupPages, action.payload.groupPages);
        userGroupsAdapter.setAll(state.userGroups, action.payload.userGroups);
        groupsAdapter.setAll(state.groups, action.payload.groups);
      })
      .addCase(pagesStore.asyncActions.listPages.rejected, () => {});
  },
  initialState,
  name: 'groups',
  reducers: {},
});

// =============================================
// selectors
// =============================================
const getGroupIds = createSelector(
  (state: RootState) => state.groups,
  (values) => values.groups.ids,
);

const getGroupEntities = createSelector(
  (state: RootState) => state.groups,
  (values) => values.groups.entities,
);

const getGroups = createSelector(
  (state: RootState) => state.groups,
  (values) => values.groups.ids.map((id) => values.groups.entities[id]),
);

const getGroupValues = createSelector(
  (state: RootState, groupId: string) => state.groups.groups.entities[groupId],
  (values) => values,
);

const getUserGroupIds = createSelector(
  (state: RootState) => state.groups,
  (values) => values.userGroups.ids,
);

const getUserGroupEntities = createSelector(
  (state: RootState) => state.groups,
  (values) => values.userGroups.entities,
);

const getUserGroups = createSelector(
  (state: RootState) => state.groups,
  (values) => values.userGroups.ids.map((id) => values.userGroups.entities[id]),
);

// =============================================
// exports
// =============================================

export const groupsStore = {
  actions: slice.actions,
  asyncActions: {
    listGroups,
    getGroup,
    createGroup,
    updateGroup,
    deleteGroup,
  },
  reducer: slice.reducer,
  selector: {
    getGroupIds,
    getGroupEntities,
    getGroups,
    getGroup: getGroupValues,
    getUserGroupIds,
    getUserGroupEntities,
    getUserGroups,
  },
};

export type GroupsState = State;
