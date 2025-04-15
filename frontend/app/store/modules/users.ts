import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../';
import { createAsyncThunk } from '@reduxjs/toolkit';
import { ENVIRONMENTS } from '@/environments';
import { api } from '@/api';
import { errorStore } from './error';
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
import { organizationsStore } from './organizations';

// =============================================
// asyncActions
// =============================================
const refreshToken = createAsyncThunk(
  'users/refreshToken',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.usersRefreshToken();

      api.setExpiresAt(res.expiresAt);
      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const obtainAuthToken = createAsyncThunk(
  'users/obtainAuthToken',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.usersObtainAuthToken();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const saveAuth = createAsyncThunk(
  'users/saveAuth',
  async (
    params: { authUrl: string; data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersSaveAuth(params);

      api.setExpiresAt(res.expiresAt);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const getUsersMe = createAsyncThunk(
  'users/getUsersMe',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.getUsersMe();

      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const listUsers = createAsyncThunk(
  'users/listUsers',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.listUsers();

      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const signout = createAsyncThunk(
  'users/signout',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.usersSignout();

      location.href = '/login';

      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestGoogleAuthLink = createAsyncThunk(
  'users/requestGoogleAuthLink',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.usersRequestGoogleAuthLink();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const authenticateWithGoogle = createAsyncThunk(
  'users/authenticateWithGoogle',
  async (
    params: {
      data: { code: string; state: string };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersAuthenticateWithGoogle(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const registerWithGoogle = createAsyncThunk(
  'users/registerWithGoogle',
  async (
    params: { data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersRegisterWithGoogle(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const invite = createAsyncThunk(
  'users/invite',
  async (
    params: { data: { emails: string[]; role: UserRole } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersInvite(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const invitationsResend = createAsyncThunk(
  'users/invitationsResend',
  async (
    params: { data: { invitationId: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersInvitationsResend(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestInvitationGoogleAuthLink = createAsyncThunk(
  'users/requestInvitationGoogleAuthLink',
  async (
    params: {
      data: {
        invitationToken: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersRequestInvitationGoogleAuthLink(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const authenticateWithInvitationGoogleAuthLink = createAsyncThunk(
  'users/authenticateWithInvitationGoogleAuthLink',
  async (
    params: { data: { code: string; state: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res =
        await api.users.usersAuthenticateWithInvitationGoogleAuthLink(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const registerWithInvitationGoogleAuthLink = createAsyncThunk(
  'users/registerWithInvitationGoogleAuthLink',
  async (
    params: { data: { token: string; firstName: string; lastName: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res =
        await api.users.usersRegisterWithInvitationGoogleAuthLink(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateUserEmail = createAsyncThunk(
  'users/updateUserEmail',
  async (
    params: { data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.updateUserEmail(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateUser = createAsyncThunk(
  'users/updateUser',
  async (
    params: {
      data: {
        firstName: string;
        lastName: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.updateUser(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const usersSendUpdateEmailInstructions = createAsyncThunk(
  'users/usersSendUpdateEmailInstructions',
  async (
    params: { data: { email: string; emailConfirmation: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersSendUpdateEmailInstructions(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestMagicLink = createAsyncThunk(
  'users/requestMagicLink',
  async (
    params: { data: { email: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersRequestMagicLink(params);
      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const authenticateWithMagicLink = createAsyncThunk(
  'users/authenticateWithMagicLink',
  async (
    params: { data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersAuthenticateWithMagicLink(params);
      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const registerWithMagicLink = createAsyncThunk(
  'users/registerWithMagicLink',
  async (
    params: { data: { token: string; firstName: string; lastName: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersRegisterWithMagicLink(params);
      api.setExpiresAt(res.expiresAt);
      return res;
    } catch (error: any) {
      if (ENVIRONMENTS.MODE === 'development') {
        console.log({ error });
      }
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestInvitationMagicLink = createAsyncThunk(
  'users/requestInvitationMagicLink',
  async (
    params: { data: { invitationToken: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersRequestInvitationMagicLink(params);
      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const authenticateWithInvitationMagicLink = createAsyncThunk(
  'users/authenticateWithInvitationMagicLink',
  async (
    params: { data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res =
        await api.users.usersAuthenticateWithInvitationMagicLink(params);
      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const registerWithInvitationMagicLink = createAsyncThunk(
  'users/registerWithInvitationMagicLink',
  async (
    params: {
      data: {
        token: string;
        firstName: string;
        lastName: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersRegisterWithInvitationMagicLink(params);
      api.setExpiresAt(res.expiresAt);
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

const usersAdapter = createEntityAdapter<User, string>({
  selectId: (user) => user.id,
});

const userInvitationsAdapter = createEntityAdapter<UserInvitation, string>({
  selectId: (userInvitation) => userInvitation.id,
});

// =============================================
// State

export type State = {
  me: User | null;
  users: EntityState<User, string>;
  userInvitations: EntityState<UserInvitation, string>;
  isGetMeWaiting: boolean;
  isRefreshTokenWaiting: boolean;
  isListUsersWaiting: boolean;
  isRequestMagicLinkWaiting: boolean;
  isRegisterWithMagicLinkWaiting: boolean;
  isAuthenticateWithMagicLinkWaiting: boolean;
  isAuthChecked: boolean;
  isAuthSucceeded: boolean;
  isAuthFailed: boolean;
  isInviteWaiting: boolean;
  isSaveAuthWaiting: boolean;
  isSignoutWaiting: boolean;
  isObtainAuthTokenWaiting: boolean;
  isRequestGoogleAuthLinkWaiting: boolean;
  isAuthenticateWithGoogleWaiting: boolean;
  isRegisterWithGoogleWaiting: boolean;
  isUpdateUserWaiting: boolean;
  isUpdateUserEmailWaiting: boolean;
  isUsersSendUpdateEmailInstructionsWaiting: boolean;
  isInvitationsResendWaiting: boolean;
  isRequestInvitationMagicLinkWaiting: boolean;
  isAuthenticateWithInvitationMagicLinkWaiting: boolean;
  isRegisterWithInvitationMagicLinkWaiting: boolean;
  isRequestInvitationGoogleAuthLinkWaiting: boolean;
  isAuthenticateWithInvitationGoogleAuthLinkWaiting: boolean;
  isRegisterWithInvitationGoogleAuthLinkWaiting: boolean;
};

const initialState: State = {
  me: null,
  users: usersAdapter.getInitialState(),
  userInvitations: userInvitationsAdapter.getInitialState(),
  isGetMeWaiting: false,
  isRefreshTokenWaiting: false,
  isListUsersWaiting: false,
  isRequestMagicLinkWaiting: false,
  isRegisterWithMagicLinkWaiting: false,
  isAuthenticateWithMagicLinkWaiting: false,
  isAuthChecked: false,
  isAuthSucceeded: false,
  isAuthFailed: false,
  isInviteWaiting: false,
  isObtainAuthTokenWaiting: false,
  isSaveAuthWaiting: false,
  isSignoutWaiting: false,
  isRequestGoogleAuthLinkWaiting: false,
  isAuthenticateWithGoogleWaiting: false,
  isRegisterWithGoogleWaiting: false,
  isUpdateUserWaiting: false,
  isUpdateUserEmailWaiting: false,
  isUsersSendUpdateEmailInstructionsWaiting: false,
  isInvitationsResendWaiting: false,
  isRequestInvitationMagicLinkWaiting: false,
  isAuthenticateWithInvitationMagicLinkWaiting: false,
  isRegisterWithInvitationMagicLinkWaiting: false,
  isRequestInvitationGoogleAuthLinkWaiting: false,
  isAuthenticateWithInvitationGoogleAuthLinkWaiting: false,
  isRegisterWithInvitationGoogleAuthLinkWaiting: false,
};
// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // refreshToken
      .addCase(refreshToken.pending, (state) => {
        state.isRefreshTokenWaiting = true;
      })
      .addCase(refreshToken.fulfilled, (state) => {
        state.isRefreshTokenWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(refreshToken.rejected, (state) => {
        state.isRefreshTokenWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // obtainAuthToken
      .addCase(obtainAuthToken.pending, (state) => {
        state.isSaveAuthWaiting = true;
      })
      .addCase(obtainAuthToken.fulfilled, (state) => {
        state.isSaveAuthWaiting = false;
      })
      .addCase(obtainAuthToken.rejected, (state) => {
        state.isSaveAuthWaiting = false;
      })

      // saveAuth
      .addCase(saveAuth.pending, (state) => {
        state.isSaveAuthWaiting = true;
      })
      .addCase(saveAuth.fulfilled, (state) => {
        state.isSaveAuthWaiting = false;
      })
      .addCase(saveAuth.rejected, (state) => {
        state.isSaveAuthWaiting = false;
      })

      // getUsersMe
      .addCase(getUsersMe.pending, (state) => {
        state.isGetMeWaiting = true;
      })
      .addCase(getUsersMe.fulfilled, (state, action) => {
        state.me = action.payload.user;
        state.isGetMeWaiting = false;
      })
      .addCase(getUsersMe.rejected, (state) => {
        state.isGetMeWaiting = false;
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

      // requestMagicLink
      .addCase(requestMagicLink.pending, (state) => {
        state.isRequestMagicLinkWaiting = true;
      })
      .addCase(requestMagicLink.fulfilled, (state) => {
        state.isRequestMagicLinkWaiting = false;
      })
      .addCase(requestMagicLink.rejected, (state) => {
        state.isRequestMagicLinkWaiting = false;
      })

      // authenticateWithMagicLink
      .addCase(authenticateWithMagicLink.pending, (state) => {
        state.isAuthenticateWithMagicLinkWaiting = true;
      })
      .addCase(authenticateWithMagicLink.fulfilled, (state) => {
        state.isAuthenticateWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(authenticateWithMagicLink.rejected, (state) => {
        state.isAuthenticateWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // registerWithMagicLink
      .addCase(registerWithMagicLink.pending, (state) => {
        state.isRegisterWithMagicLinkWaiting = true;
      })
      .addCase(registerWithMagicLink.fulfilled, (state) => {
        state.isRegisterWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(registerWithMagicLink.rejected, (state) => {
        state.isRegisterWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // signout
      .addCase(signout.pending, (state) => {
        state.isSignoutWaiting = true;
      })
      .addCase(signout.fulfilled, (state) => {
        state.isSignoutWaiting = false;
      })
      .addCase(signout.rejected, (state) => {
        state.isSignoutWaiting = false;
      })

      // invite
      .addCase(invite.pending, (state) => {
        state.isInviteWaiting = true;
      })
      .addCase(invite.fulfilled, (state) => {
        state.isInviteWaiting = false;
      })
      .addCase(invite.rejected, (state) => {
        state.isInviteWaiting = false;
      })

      // requestGoogleAuthLink
      .addCase(requestGoogleAuthLink.pending, (state) => {
        state.isRequestGoogleAuthLinkWaiting = true;
      })
      .addCase(requestGoogleAuthLink.fulfilled, (state) => {
        state.isRequestGoogleAuthLinkWaiting = false;
      })
      .addCase(requestGoogleAuthLink.rejected, (state) => {
        state.isRequestGoogleAuthLinkWaiting = false;
      })

      // authenticateWithGoogle
      .addCase(authenticateWithGoogle.pending, (state) => {
        state.isAuthenticateWithGoogleWaiting = true;
      })
      .addCase(authenticateWithGoogle.fulfilled, (state) => {
        state.isAuthenticateWithGoogleWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(authenticateWithGoogle.rejected, (state) => {
        state.isAuthenticateWithGoogleWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // registerWithGoogle
      .addCase(registerWithGoogle.pending, (state) => {
        state.isRegisterWithGoogleWaiting = true;
      })
      .addCase(registerWithGoogle.fulfilled, (state) => {
        state.isRegisterWithGoogleWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(registerWithGoogle.rejected, (state) => {
        state.isRegisterWithGoogleWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
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

      // updateUser
      .addCase(updateUser.pending, (state) => {
        state.isUpdateUserWaiting = true;
      })
      .addCase(updateUser.fulfilled, (state, action) => {
        state.isUpdateUserWaiting = false;
        if (state.me) {
          state.me = action.payload.user;
        }
      })
      .addCase(updateUser.rejected, (state) => {
        state.isUpdateUserWaiting = false;
      })

      // updateUserEmail
      .addCase(updateUserEmail.pending, (state) => {
        state.isUpdateUserEmailWaiting = true;
      })
      .addCase(updateUserEmail.fulfilled, (state, action) => {
        state.isUpdateUserEmailWaiting = false;
        state.me = action.payload.user;
      })
      .addCase(updateUserEmail.rejected, (state) => {
        state.isUpdateUserEmailWaiting = false;
      })

      // usersSendUpdateEmailInstructions
      .addCase(usersSendUpdateEmailInstructions.pending, (state) => {
        state.isUsersSendUpdateEmailInstructionsWaiting = true;
      })
      .addCase(usersSendUpdateEmailInstructions.fulfilled, (state) => {
        state.isUsersSendUpdateEmailInstructionsWaiting = false;
      })
      .addCase(usersSendUpdateEmailInstructions.rejected, (state) => {
        state.isUsersSendUpdateEmailInstructionsWaiting = false;
      })

      // listPages
      .addCase(pagesStore.asyncActions.listPages.pending, () => {})
      .addCase(pagesStore.asyncActions.listPages.fulfilled, (state, action) => {
        usersAdapter.setAll(state.users, action.payload.users);
      })
      .addCase(pagesStore.asyncActions.listPages.rejected, () => {})

      // updateOrganizationUser
      .addCase(
        organizationsStore.asyncActions.updateOrganizationUser.pending,
        () => {},
      )
      .addCase(
        organizationsStore.asyncActions.updateOrganizationUser.fulfilled,
        (state, action) => {
          usersAdapter.updateOne(state.users, {
            id: action.payload.id,
            changes: action.payload,
          });
        },
      )
      .addCase(
        organizationsStore.asyncActions.updateOrganizationUser.rejected,
        () => {},
      )

      // invitationsResend
      .addCase(invitationsResend.pending, (state) => {
        state.isInvitationsResendWaiting = true;
      })
      .addCase(invitationsResend.fulfilled, (state) => {
        state.isInvitationsResendWaiting = false;
      })
      .addCase(invitationsResend.rejected, (state) => {
        state.isInvitationsResendWaiting = false;
      })

      // requestInvitationMagicLink
      .addCase(requestInvitationMagicLink.pending, (state) => {
        state.isRequestInvitationMagicLinkWaiting = true;
      })
      .addCase(requestInvitationMagicLink.fulfilled, (state) => {
        state.isRequestInvitationMagicLinkWaiting = false;
      })
      .addCase(requestInvitationMagicLink.rejected, (state) => {
        state.isRequestInvitationMagicLinkWaiting = false;
      })

      // authenticateWithInvitationMagicLink
      .addCase(authenticateWithInvitationMagicLink.pending, (state) => {
        state.isAuthenticateWithInvitationMagicLinkWaiting = true;
      })
      .addCase(authenticateWithInvitationMagicLink.fulfilled, (state) => {
        state.isAuthenticateWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(authenticateWithInvitationMagicLink.rejected, (state) => {
        state.isAuthenticateWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // registerWithInvitationMagicLink
      .addCase(registerWithInvitationMagicLink.pending, (state) => {
        state.isRegisterWithInvitationMagicLinkWaiting = true;
      })
      .addCase(registerWithInvitationMagicLink.fulfilled, (state) => {
        state.isRegisterWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(registerWithInvitationMagicLink.rejected, (state) => {
        state.isRegisterWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // requestInvitationGoogleAuthLink
      .addCase(requestInvitationGoogleAuthLink.pending, (state) => {
        state.isRequestInvitationGoogleAuthLinkWaiting = true;
      })
      .addCase(requestInvitationGoogleAuthLink.fulfilled, (state) => {
        state.isRequestInvitationGoogleAuthLinkWaiting = false;
      })
      .addCase(requestInvitationGoogleAuthLink.rejected, (state) => {
        state.isRequestInvitationGoogleAuthLinkWaiting = false;
      })

      // authenticateWithInvitationGoogleAuthLink
      .addCase(authenticateWithInvitationGoogleAuthLink.pending, (state) => {
        state.isAuthenticateWithInvitationGoogleAuthLinkWaiting = true;
      })
      .addCase(authenticateWithInvitationGoogleAuthLink.fulfilled, (state) => {
        state.isAuthenticateWithInvitationGoogleAuthLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(authenticateWithInvitationGoogleAuthLink.rejected, (state) => {
        state.isAuthenticateWithInvitationGoogleAuthLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // registerWithInvitationGoogleAuthLink
      .addCase(registerWithInvitationGoogleAuthLink.pending, (state) => {
        state.isRegisterWithInvitationGoogleAuthLinkWaiting = true;
      })
      .addCase(registerWithInvitationGoogleAuthLink.fulfilled, (state) => {
        state.isRegisterWithInvitationGoogleAuthLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(registerWithInvitationGoogleAuthLink.rejected, (state) => {
        state.isRegisterWithInvitationGoogleAuthLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      });
  },
  initialState,
  name: 'users',
  reducers: {},
});

// =============================================
// selectors
// =============================================

const getMe = createSelector(
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
      state.users.isAuthChecked &&
      (state.users.isAuthFailed ||
        (state.users.isAuthSucceeded && state.users.me));
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
    refreshToken,
    obtainAuthToken,
    saveAuth,
    getUsersMe,
    listUsers,
    signout,
    requestGoogleAuthLink,
    authenticateWithGoogle,
    registerWithGoogle,
    invite,
    invitationsResend,
    updateUserEmail,
    updateUser,
    usersSendUpdateEmailInstructions,
    requestMagicLink,
    authenticateWithMagicLink,
    registerWithMagicLink,
    requestInvitationMagicLink,
    authenticateWithInvitationMagicLink,
    registerWithInvitationMagicLink,
    requestInvitationGoogleAuthLink,
    authenticateWithInvitationGoogleAuthLink,
    registerWithInvitationGoogleAuthLink,
  },
  reducer: slice.reducer,
  selector: {
    getMe,
    getUserIds,
    getUserEntities,
    getUsers,
    getUserInvitations,
    getUser,
    getSubDomainMatched,
  },
};

export type UsersState = State;
