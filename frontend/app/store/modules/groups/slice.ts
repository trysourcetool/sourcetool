import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';
import type { Group, GroupPage, UserGroup } from '@/api/modules/groups';
import { pagesStore } from '../pages';

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
      .addCase(asyncActions.listGroups.pending, (state) => {
        state.isListGroupsWaiting = true;
      })
      .addCase(asyncActions.listGroups.fulfilled, (state, action) => {
        state.isListGroupsWaiting = false;
        groupsAdapter.setAll(state.groups, action.payload.groups);
        userGroupsAdapter.setAll(state.userGroups, action.payload.userGroups);
      })
      .addCase(asyncActions.listGroups.rejected, (state) => {
        state.isListGroupsWaiting = false;
      })

      // getGroup
      .addCase(asyncActions.getGroup.pending, (state) => {
        state.isGetGroupWaiting = true;
      })
      .addCase(asyncActions.getGroup.fulfilled, (state, action) => {
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
      .addCase(asyncActions.getGroup.rejected, (state) => {
        state.isGetGroupWaiting = false;
      })

      // createEnvironment
      .addCase(asyncActions.createGroup.pending, (state) => {
        state.isCreateGroupWaiting = true;
      })
      .addCase(asyncActions.createGroup.fulfilled, (state, action) => {
        state.isCreateGroupWaiting = false;
        groupsAdapter.addOne(state.groups, action.payload.group);
      })
      .addCase(asyncActions.createGroup.rejected, (state) => {
        state.isCreateGroupWaiting = false;
      })

      // updateEnvironment
      .addCase(asyncActions.updateGroup.pending, (state) => {
        state.isUpdateGroupWaiting = true;
      })
      .addCase(asyncActions.updateGroup.fulfilled, (state, action) => {
        state.isUpdateGroupWaiting = false;
        groupsAdapter.updateOne(state.groups, {
          id: action.payload.group.id,
          changes: action.payload.group,
        });
      })
      .addCase(asyncActions.updateGroup.rejected, (state) => {
        state.isUpdateGroupWaiting = false;
      })

      // deleteEnvironment
      .addCase(asyncActions.deleteGroup.pending, (state) => {
        state.isDeleteGroupWaiting = true;
      })
      .addCase(asyncActions.deleteGroup.fulfilled, (state, action) => {
        state.isDeleteGroupWaiting = false;
        groupsAdapter.removeOne(state.groups, action.payload.group.id);
      })
      .addCase(asyncActions.deleteGroup.rejected, (state) => {
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
