import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import type { ErrorResponse } from '@/api/instance';
import { createSlice } from '@reduxjs/toolkit';

// =============================================
// asyncActions
// =============================================

const refreshToken = createAsyncThunk(
  'auth/refreshToken',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.auth.refreshToken();
      api.setExpiresAt(res.expiresAt);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const obtainAuthToken = createAsyncThunk(
  'auth/obtainAuthToken',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.auth.obtainAuthToken();
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const saveAuth = createAsyncThunk(
  'auth/saveAuth',
  async (
    params: { authUrl: string; data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.saveAuth(params);
      api.setExpiresAt(res.expiresAt);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const logout = createAsyncThunk(
  'auth/logout',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.auth.logout();
      location.href = '/login';
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestGoogleAuthLink = createAsyncThunk(
  'auth/requestGoogleAuthLink',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.auth.requestGoogleAuthLink();
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const authenticateWithGoogle = createAsyncThunk(
  'auth/authenticateWithGoogle',
  async (
    params: { data: { code: string; state: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.authenticateWithGoogle(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const registerWithGoogle = createAsyncThunk(
  'auth/registerWithGoogle',
  async (
    params: { data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.registerWithGoogle(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestMagicLink = createAsyncThunk(
  'auth/requestMagicLink',
  async (
    params: { data: { email: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.requestMagicLink(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const authenticateWithMagicLink = createAsyncThunk(
  'auth/authenticateWithMagicLink',
  async (
    params: { data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.authenticateWithMagicLink(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const registerWithMagicLink = createAsyncThunk(
  'auth/registerWithMagicLink',
  async (
    params: { data: { token: string; firstName: string; lastName: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.registerWithMagicLink(params);
      api.setExpiresAt(res.expiresAt);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestInvitationGoogleAuthLink = createAsyncThunk(
  'auth/requestInvitationGoogleAuthLink',
  async (
    params: { data: { invitationToken: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.requestInvitationGoogleAuthLink(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const requestInvitationMagicLink = createAsyncThunk(
  'auth/requestInvitationMagicLink',
  async (
    params: { data: { invitationToken: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.requestInvitationMagicLink(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const authenticateWithInvitationMagicLink = createAsyncThunk(
  'auth/authenticateWithInvitationMagicLink',
  async (
    params: { data: { token: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.authenticateWithInvitationMagicLink(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const registerWithInvitationMagicLink = createAsyncThunk(
  'auth/registerWithInvitationMagicLink',
  async (
    params: { data: { token: string; firstName: string; lastName: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.auth.registerWithInvitationMagicLink(params);
      api.setExpiresAt(res.expiresAt);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

// =============================================
// State
// =============================================

export type State = {
  isRefreshTokenWaiting: boolean;
  isObtainAuthTokenWaiting: boolean;
  isSaveAuthWaiting: boolean;
  isLogoutWaiting: boolean;
  isRequestGoogleAuthLinkWaiting: boolean;
  isAuthenticateWithGoogleWaiting: boolean;
  isRegisterWithGoogleWaiting: boolean;
  isRequestMagicLinkWaiting: boolean;
  isAuthenticateWithMagicLinkWaiting: boolean;
  isRegisterWithMagicLinkWaiting: boolean;
  isRequestInvitationGoogleAuthLinkWaiting: boolean;
  isRequestInvitationMagicLinkWaiting: boolean;
  isAuthenticateWithInvitationMagicLinkWaiting: boolean;
  isRegisterWithInvitationMagicLinkWaiting: boolean;
  isAuthChecked: boolean;
  isAuthSucceeded: boolean;
  isAuthFailed: boolean;
};

const initialState: State = {
  isRefreshTokenWaiting: false,
  isObtainAuthTokenWaiting: false,
  isSaveAuthWaiting: false,
  isLogoutWaiting: false,
  isRequestGoogleAuthLinkWaiting: false,
  isAuthenticateWithGoogleWaiting: false,
  isRegisterWithGoogleWaiting: false,
  isRequestMagicLinkWaiting: false,
  isAuthenticateWithMagicLinkWaiting: false,
  isRegisterWithMagicLinkWaiting: false,
  isRequestInvitationGoogleAuthLinkWaiting: false,
  isRequestInvitationMagicLinkWaiting: false,
  isAuthenticateWithInvitationMagicLinkWaiting: false,
  isRegisterWithInvitationMagicLinkWaiting: false,
  isAuthChecked: false,
  isAuthSucceeded: false,
  isAuthFailed: false,
};

// =============================================
// slice
// =============================================

export const slice = createSlice({
  name: 'auth',
  initialState,
  reducers: {},
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
        state.isObtainAuthTokenWaiting = true;
      })
      .addCase(obtainAuthToken.fulfilled, (state) => {
        state.isObtainAuthTokenWaiting = false;
      })
      .addCase(obtainAuthToken.rejected, (state) => {
        state.isObtainAuthTokenWaiting = false;
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

      // logout
      .addCase(logout.pending, (state) => {
        state.isLogoutWaiting = true;
      })
      .addCase(logout.fulfilled, (state) => {
        state.isLogoutWaiting = false;
      })
      .addCase(logout.rejected, (state) => {
        state.isLogoutWaiting = false;
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
      })
      .addCase(authenticateWithInvitationMagicLink.rejected, (state) => {
        state.isAuthenticateWithInvitationMagicLinkWaiting = false;
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
      });
  },
});

// =============================================
// selectors
// =============================================

const selectIsAuthChecked = (state: { auth: State }) => state.auth.isAuthChecked;
const selectIsAuthSucceeded = (state: { auth: State }) => state.auth.isAuthSucceeded;
const selectIsAuthFailed = (state: { auth: State }) => state.auth.isAuthFailed;

// =============================================
// exports
// =============================================

export const authStore = {
  actions: slice.actions,
  asyncActions: {
    refreshToken,
    obtainAuthToken,
    saveAuth,
    logout,
    requestGoogleAuthLink,
    authenticateWithGoogle,
    registerWithGoogle,
    requestMagicLink,
    authenticateWithMagicLink,
    registerWithMagicLink,
    requestInvitationGoogleAuthLink,
    requestInvitationMagicLink,
    authenticateWithInvitationMagicLink,
    registerWithInvitationMagicLink,
  },
  reducer: slice.reducer,
  selector: {
    selectIsAuthChecked,
    selectIsAuthSucceeded,
    selectIsAuthFailed,
  },
};

export type AuthState = State;
