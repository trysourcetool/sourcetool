import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';
import type { UserRole } from '@/api/modules/users';

export const createOrganization = createAsyncThunk(
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

export const checkSubdomainAvailability = createAsyncThunk(
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

export const updateOrganizationUser = createAsyncThunk(
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
