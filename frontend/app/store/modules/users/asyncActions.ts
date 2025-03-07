import { createAsyncThunk } from '@reduxjs/toolkit';
import { ENVIRONMENTS } from '@/environments';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';
import type { UserRole } from '@/api/modules/users';

export const refreshToken = createAsyncThunk(
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

export const obtainAuthToken = createAsyncThunk(
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

export const saveAuth = createAsyncThunk(
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

export const getUsersMe = createAsyncThunk(
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

export const listUsers = createAsyncThunk(
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

export const signin = createAsyncThunk(
  'users/signin',
  async (
    params: { data: { email: string; password: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersSignin(params);

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

export const signout = createAsyncThunk(
  'users/signout',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.usersSignout();

      location.href = '/signin';

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

export const signup = createAsyncThunk(
  'users/signup',
  async (
    params: {
      data: {
        firstName: string;
        lastName: string;
        password: string;
        passwordConfirmation: string;
        token: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersSignup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const signupInstructions = createAsyncThunk(
  'users/signupInstructions',
  async (
    params: { data: { email: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersSignupInstructions(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const oauthGoogleUserInfo = createAsyncThunk(
  'users/oauthGoogleUserInfo',
  async (
    params: { data: { code: string; state: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersOauthGoogleUserInfo(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const oauthGoogleSignin = createAsyncThunk(
  'users/oauthGoogleSignin',
  async (
    params: { data: { sessionToken: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersOauthGoogleSignin(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const oauthGoogleSignup = createAsyncThunk(
  'users/oauthGoogleSignup',
  async (
    params: {
      data: { firstName: string; lastName: string; sessionToken: string };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersOauthGoogleSignup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const oauthGoogleAuthCodeUrl = createAsyncThunk(
  'users/oauthGoogleAuthCodeUrl',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.users.usersOauthGoogleAuthCodeUrl();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const invite = createAsyncThunk(
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

export const invitationsResend = createAsyncThunk(
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

export const invitationsSignup = createAsyncThunk(
  'users/invitationsSignup',
  async (
    params: {
      data: {
        firstName: string;
        invitationToken: string;
        lastName: string;
        password: string;
        passwordConfirmation: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersInvitationsSignup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const invitationsSignin = createAsyncThunk(
  'users/invitationsSignin',
  async (
    params: { data: { invitationToken: string; password: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersInvitationsSignin(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const invitationsOauthGoogleUserInfo = createAsyncThunk(
  'users/invitationsOauthGoogleUserInfo',
  async (
    params: { data: { code: string; state: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersInvitationsOauthGoogleUserInfo(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const invitationsOauthGoogleSignup = createAsyncThunk(
  'users/invitationsOauthGoogleSignup',
  async (
    params: {
      data: {
        firstName: string;
        invitationToken: string;
        lastName: string;
        sessionToken: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersInvitationsOauthGoogleSignup(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const invitationsOauthGoogleSignin = createAsyncThunk(
  'users/invitationsOauthGoogleSignin',
  async (
    params: { data: { invitationToken: string; sessionToken: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.usersInvitationsOauthGoogleSignin(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const invitationsOauthGoogleAuthCodeUrl = createAsyncThunk(
  'users/invitationsOauthGoogleAuthCodeUrl',
  async (
    params: { data: { invitationToken: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res =
        await api.users.usersInvitationsOauthGoogleAuthCodeUrl(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const updateUserEmail = createAsyncThunk(
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

export const updateUserPassword = createAsyncThunk(
  'users/updateUserPassword',
  async (
    params: {
      data: {
        currentPassword: string;
        password: string;
        passwordConfirmation: string;
      };
    },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.users.updateUserPassword(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

export const updateUser = createAsyncThunk(
  'users/updateUser',
  async (
    params: { data: { firstName: string; lastName: string } },
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

export const usersSendUpdateEmailInstructions = createAsyncThunk(
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
