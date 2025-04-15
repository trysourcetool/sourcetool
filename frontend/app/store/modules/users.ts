import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../';
import { createAsyncThunk } from '@reduxjs/toolkit';
import { ENVIRONMENTS } from '@/environments';
import { api } from '@/api';
import type { ErrorResponse } from '@/api/instance';
import type { UserRole } from '@/api/modules/users';
import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import type { User, UserInvitation } from '@/api/modules/users';
import { groupsStore } from './groups';
import { pagesStore } from './pages';

// =============================================
// asyncActions
// =============================================

const getMe = createAsyncThunk(
  'users/getMe',
  async (_, { rejectWithValue }) => {
    try {
      const res = await api.users.getMe();
      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateMe = createAsyncThunk(
  'users/updateMe',
  async (
    params: {
      data: {
        firstName?: string;
        lastName?: string;
      };
    },
    { rejectWithValue },
  ) => {
    try {
      const res = await api.users.updateMe(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const sendUpdateMeEmailInstructions = createAsyncThunk(
  'users/sendUpdateMeEmailInstructions',
  async (
    params: { data: { email: string; emailConfirmation: string } },
    { rejectWithValue },
  ) => {
    try {
      const res = await api.users.sendUpdateMeEmailInstructions(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateMeEmail = createAsyncThunk(
  'users/updateMeEmail',
  async (
    params: {
      data: {
        token: string;
      };
    },
    { rejectWithValue },
  ) => {
    try {
      const res = await api.users.updateMeEmail(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const listUsers = createAsyncThunk(
  'users/listUsers',
  async (_, { rejectWithValue }) => {
    try {
      const res = await api.users.listUsers();
      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateUser = createAsyncThunk(
  'users/updateUser',
  async (
    params: {
      userId: string;
      data: {
        role?: UserRole;
        groupIds?: string[];
      };
    },
    { rejectWithValue },
  ) => {
    try {
      const res = await api.users.updateUser(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const deleteUser = createAsyncThunk(
  'users/deleteUser',
  async (params: { userId: string }, { rejectWithValue }) => {
    try {
      const res = await api.users.deleteUser(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const createUserInvitations = createAsyncThunk(
  'users/createUserInvitations',
  async (
    params: { data: { emails: string[]; role: UserRole } },
    { rejectWithValue },
  ) => {
    try {
      const res = await api.users.createUserInvitations(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const resendUserInvitation = createAsyncThunk(
  'users/resendUserInvitation',
  async (params: { invitationId: string }, { rejectWithValue }) => {
    try {
      const res = await api.users.resendUserInvitation(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

// =============================================
// schema
// =============================================

const usersAdapter = createEntityAdapter<User, string>({
  selectId: (user) => user.id,
});

const userInvitationsAdapter = createEntityAdapter<UserInvitation, string>({
  selectId: (userInvitation) => userInvitation.id,
});

// =============================================
// State
// =============================================

export type State = {
  me: User | null;
  users: EntityState<User, string>;
  userInvitations: EntityState<UserInvitation, string>;
  isGetMeWaiting: boolean;
  isUpdateMeWaiting: boolean;
  isSendUpdateMeEmailInstructionsWaiting: boolean;
  isUpdateMeEmailWaiting: boolean;
  isListUsersWaiting: boolean;
  isUpdateUserWaiting: boolean;
  isDeleteUserWaiting: boolean;
  isCreateUserInvitationsWaiting: boolean;
  isResendUserInvitationWaiting: boolean;
};

const initialState: State = {
  me: null,
  users: usersAdapter.getInitialState(),
  userInvitations: userInvitationsAdapter.getInitialState(),
  isGetMeWaiting: false,
  isUpdateMeWaiting: false,
  isSendUpdateMeEmailInstructionsWaiting: false,
  isUpdateMeEmailWaiting: false,
  isListUsersWaiting: false,
  isUpdateUserWaiting: false,
  isDeleteUserWaiting: false,
  isCreateUserInvitationsWaiting: false,
  isResendUserInvitationWaiting: false,
};

// =============================================
// slice
// =============================================

export const slice = createSlice({
  name: 'users',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      // getMe
      .addCase(getMe.pending, (state) => {
        state.isGetMeWaiting = true;
      })
      .addCase(getMe.fulfilled, (state, action) => {
        state.me = action.payload.user;
        state.isGetMeWaiting = false;
      })
      .addCase(getMe.rejected, (state) => {
        state.isGetMeWaiting = false;
      })

      // updateMe
      .addCase(updateMe.pending, (state) => {
        state.isUpdateMeWaiting = true;
      })
      .addCase(updateMe.fulfilled, (state, action) => {
        state.me = action.payload.user;
        state.isUpdateMeWaiting = false;
      })
      .addCase(updateMe.rejected, (state) => {
        state.isUpdateMeWaiting = false;
      })

      // sendUpdateMeEmailInstructions
      .addCase(sendUpdateMeEmailInstructions.pending, (state) => {
        state.isSendUpdateMeEmailInstructionsWaiting = true;
      })
      .addCase(sendUpdateMeEmailInstructions.fulfilled, (state) => {
        state.isSendUpdateMeEmailInstructionsWaiting = false;
      })
      .addCase(sendUpdateMeEmailInstructions.rejected, (state) => {
        state.isSendUpdateMeEmailInstructionsWaiting = false;
      })

      // updateMeEmail
      .addCase(updateMeEmail.pending, (state) => {
        state.isUpdateMeEmailWaiting = true;
      })
      .addCase(updateMeEmail.fulfilled, (state, action) => {
        state.me = action.payload.user;
        state.isUpdateMeEmailWaiting = false;
      })
      .addCase(updateMeEmail.rejected, (state) => {
        state.isUpdateMeEmailWaiting = false;
      })

      // createUserInvitations
      .addCase(createUserInvitations.pending, (state) => {
        state.isCreateUserInvitationsWaiting = true;
      })
      .addCase(createUserInvitations.fulfilled, (state) => {
        state.isCreateUserInvitationsWaiting = false;
      })
      .addCase(createUserInvitations.rejected, (state) => {
        state.isCreateUserInvitationsWaiting = false;
      })

      // resendUserInvitation
      .addCase(resendUserInvitation.pending, (state) => {
        state.isResendUserInvitationWaiting = true;
      })
      .addCase(resendUserInvitation.fulfilled, (state) => {
        state.isResendUserInvitationWaiting = false;
      })
      .addCase(resendUserInvitation.rejected, (state) => {
        state.isResendUserInvitationWaiting = false;
      })

      // listUsers
      .addCase(listUsers.pending, (state) => {
        state.isListUsersWaiting = true;
      })
      .addCase(listUsers.fulfilled, (state, action) => {
        usersAdapter.setAll(state.users, action.payload.users);
        userInvitationsAdapter.setAll(
          state.userInvitations,
          action.payload.userInvitations,
        );
        state.isListUsersWaiting = false;
      })
      .addCase(listUsers.rejected, (state) => {
        state.isListUsersWaiting = false;
      })

      // updateUser
      .addCase(updateUser.pending, (state) => {
        state.isUpdateUserWaiting = true;
      })
      .addCase(updateUser.fulfilled, (state, action) => {
        usersAdapter.updateOne(state.users, {
          id: action.payload.user.id,
          changes: action.payload.user,
        });
        state.isUpdateUserWaiting = false;
      })
      .addCase(updateUser.rejected, (state) => {
        state.isUpdateUserWaiting = false;
      })

      // deleteUser
      .addCase(deleteUser.pending, (state) => {
        state.isDeleteUserWaiting = true;
      })
      .addCase(deleteUser.fulfilled, (state) => {
        state.isDeleteUserWaiting = false;
      })
      .addCase(deleteUser.rejected, (state) => {
        state.isDeleteUserWaiting = false;
      })

      // getUserGroups
      .addCase(groupsStore.asyncActions.listGroups.pending, () => {})
      .addCase(
        groupsStore.asyncActions.listGroups.fulfilled,
        (state, action) => {
          usersAdapter.setAll(state.users, action.payload.users);
        },
      )
      .addCase(groupsStore.asyncActions.listGroups.rejected, () => {})

      // listPages
      .addCase(pagesStore.asyncActions.listPages.pending, () => {})
      .addCase(pagesStore.asyncActions.listPages.fulfilled, (state, action) => {
        usersAdapter.setAll(state.users, action.payload.users);
      })
      .addCase(pagesStore.asyncActions.listPages.rejected, () => {});
  },
});

// =============================================
// selectors
// =============================================

const getUserMe = createSelector(
  (state: RootState) => state.users,
  (values) => values.me,
);

const getUserIds = createSelector(
  (state: RootState) => state.users,
  (values) => values.users.ids,
);

const getUserEntities = createSelector(
  (state: RootState) => state.users,
  (values) => values.users.entities,
);

const getUsers = createSelector(
  (state: RootState) => state.users,
  (values) => values.users.ids.map((id) => values.users.entities[id]),
);

const getUserInvitations = createSelector(
  (state: RootState) => state.users,
  (values) =>
    values.userInvitations.ids.map((id) => values.userInvitations.entities[id]),
);

const getUser = createSelector(
  (state: RootState, userId: string) => ({
    users: state.users.users,
    userId,
  }),
  ({ users, userId }) => users.entities[userId],
);

const getSubDomainMatched = createSelector(
  (state: RootState, subDomain: string | null) => {
    const isAuthChecked =
      state.auth.isAuthChecked &&
      (state.auth.isAuthFailed ||
        (state.auth.isAuthSucceeded && state.users.me));
    const matched = state.users.me?.organization?.subdomain === subDomain;
    return {
      isMatched: matched,
      status: !isAuthChecked ? 'checking' : 'checked',
    } as {
      isMatched: boolean;
      status: 'checking' | 'checked';
    };
  },
  (values) => values,
);

// =============================================
// exports
// =============================================

export const usersStore = {
  actions: slice.actions,
  asyncActions: {
    getMe,
    listUsers,
    createUserInvitations,
    resendUserInvitation,
    updateUser,
    updateMe,
    updateMeEmail,
    sendUpdateMeEmailInstructions,
    deleteUser,
  },
  reducer: slice.reducer,
  selector: {
    getUserMe,
    getUserIds,
    getUserEntities,
    getUsers,
    getUserInvitations,
    getUser,
    getSubDomainMatched,
  },
};

export type UsersState = State;
