import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import type { ErrorResponse } from '@/api/instance';
import { createSlice } from '@reduxjs/toolkit';

// =============================================
// asyncActions
// =============================================

const createOrganization = createAsyncThunk(
  'organizations/createOrganization',
  async (params: { data: { subdomain: string } }, { rejectWithValue }) => {
    try {
      const res = await api.organizations.createOrganization(params);
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const checkSubdomainAvailability = createAsyncThunk(
  'organizations/checkSubdomainAvailability',
  async (params: { subdomain: string }, { rejectWithValue }) => {
    try {
      const res = await api.organizations.checkSubdomainAvailability(params);
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
  isCreateOrganizationWaiting: boolean;
  isCheckSubdomainAvailabilityWaiting: boolean;
};

const initialState: State = {
  isCreateOrganizationWaiting: false,
  isCheckSubdomainAvailabilityWaiting: false,
};

// =============================================
// slice
// =============================================

export const slice = createSlice({
  name: 'organizations',
  initialState,
  reducers: {},
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
      });
  },
});

// =============================================
// exports
// =============================================

export const organizationsStore = {
  actions: slice.actions,
  asyncActions: {
    createOrganization,
    checkSubdomainAvailability,
  },
  reducer: slice.reducer,
};

export type OrganizationsState = State;
