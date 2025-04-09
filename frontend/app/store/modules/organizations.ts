import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from './error';
import type { ErrorResponse } from '@/api/instance';
import type { UserRole } from '@/api/modules/users';
import { createSlice } from '@reduxjs/toolkit';

// =============================================
// asyncActions
// =============================================

const createOrganization = createAsyncThunk(
  'organizations/createOrganization',
  async (
    params: { data: { subdomain: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.organizations.createOrganization(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const checkSubdomainAvailability = createAsyncThunk(
  'organizations/checkSubdomainAvailability',
  async (params: { subdomain: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.organizations.checkSubdomainAvailability(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateOrganizationUser = createAsyncThunk(
  'organizations/updateOrganizationUser',
  async (
    params: { userId: string; data: { groupIds: string[]; role: UserRole } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.organizations.updateOrganizationUser(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const deleteOrganizationUser = createAsyncThunk(
  'organizations/deleteOrganizationUser',
  async (params: { userId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.organizations.deleteOrganizationUser(params);

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

// =============================================
// State

export type State = {
  isCreateOrganizationWaiting: boolean;
  isCheckSubdomainAvailabilityWaiting: boolean;
  isUpdateOrganizationUserWaiting: boolean;
  isDeleteOrganizationUserWaiting: boolean;
};

const initialState: State = {
  isCreateOrganizationWaiting: false,
  isCheckSubdomainAvailabilityWaiting: false,
  isUpdateOrganizationUserWaiting: false,
  isDeleteOrganizationUserWaiting: false,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // createOrganization
      .addCase(createOrganization.pending, (state) => {
        state.isCreateOrganizationWaiting = true;
      })
      .addCase(createOrganization.fulfilled, (state) => {
        state.isCreateOrganizationWaiting = false;
      })
      .addCase(createOrganization.rejected, (state) => {
        state.isCreateOrganizationWaiting = false;
      })
      // checkSubdomainAvailability
      .addCase(checkSubdomainAvailability.pending, (state) => {
        state.isCheckSubdomainAvailabilityWaiting = true;
      })
      .addCase(checkSubdomainAvailability.fulfilled, (state) => {
        state.isCheckSubdomainAvailabilityWaiting = false;
      })
      .addCase(checkSubdomainAvailability.rejected, (state) => {
        state.isCheckSubdomainAvailabilityWaiting = false;
      })
      // updateOrganizationUser
      .addCase(updateOrganizationUser.pending, (state) => {
        state.isUpdateOrganizationUserWaiting = true;
      })
      .addCase(updateOrganizationUser.fulfilled, (state) => {
        state.isUpdateOrganizationUserWaiting = false;
      })
      .addCase(updateOrganizationUser.rejected, (state) => {
        state.isUpdateOrganizationUserWaiting = false;
      })
      // deleteOrganizationUser
      .addCase(deleteOrganizationUser.pending, (state) => {
        state.isDeleteOrganizationUserWaiting = true;
      })
      .addCase(deleteOrganizationUser.fulfilled, (state) => {
        state.isDeleteOrganizationUserWaiting = false;
      })
      .addCase(deleteOrganizationUser.rejected, (state) => {
        state.isDeleteOrganizationUserWaiting = false;
      });
  },
  initialState,
  name: 'organizations',
  reducers: {},
});

// =============================================
// exports
// =============================================

export const organizationsStore = {
  actions: slice.actions,
  asyncActions: {
    createOrganization,
    checkSubdomainAvailability,
    updateOrganizationUser,
    deleteOrganizationUser,
  },
  reducer: slice.reducer,
};

export type OrganizationsState = State;
