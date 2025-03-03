import { createSlice } from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';

// =============================================
// schema

// =============================================
// State

export type State = {
  isCreateOrganizationWaiting: boolean;
  isCheckSubdomainAvailabilityWaiting: boolean;
  isUpdateOrganizationUserWaiting: boolean;
};

const initialState: State = {
  isCreateOrganizationWaiting: false,
  isCheckSubdomainAvailabilityWaiting: false,
  isUpdateOrganizationUserWaiting: false,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // createOrganization
      .addCase(asyncActions.createOrganization.pending, (state) => {
        state.isCreateOrganizationWaiting = true;
      })
      .addCase(asyncActions.createOrganization.fulfilled, (state) => {
        state.isCreateOrganizationWaiting = false;
      })
      .addCase(asyncActions.createOrganization.rejected, (state) => {
        state.isCreateOrganizationWaiting = false;
      })
      // checkSubdomainAvailability
      .addCase(asyncActions.checkSubdomainAvailability.pending, (state) => {
        state.isCheckSubdomainAvailabilityWaiting = true;
      })
      .addCase(asyncActions.checkSubdomainAvailability.fulfilled, (state) => {
        state.isCheckSubdomainAvailabilityWaiting = false;
      })
      .addCase(asyncActions.checkSubdomainAvailability.rejected, (state) => {
        state.isCheckSubdomainAvailabilityWaiting = false;
      })
      // updateOrganizationUser
      .addCase(asyncActions.updateOrganizationUser.pending, (state) => {
        state.isUpdateOrganizationUserWaiting = true;
      })
      .addCase(asyncActions.updateOrganizationUser.fulfilled, (state) => {
        state.isUpdateOrganizationUserWaiting = false;
      })
      .addCase(asyncActions.updateOrganizationUser.rejected, (state) => {
        state.isUpdateOrganizationUserWaiting = false;
      });
  },
  initialState,
  name: 'organizations',
  reducers: {},
});
