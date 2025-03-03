import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';
import type { Environment } from '@/api/modules/environments';

// =============================================
// schema

const environmentsAdapter = createEntityAdapter<Environment, string>({
  selectId: (environment) => environment.id,
});

// =============================================
// State

export type State = {
  environments: EntityState<Environment, string>;
  isListEnvironmentsWaiting: boolean;
  isGetEnvironmentWaiting: boolean;
  isCreateEnvironmentWaiting: boolean;
  isUpdateEnvironmentWaiting: boolean;
  isDeleteEnvironmentWaiting: boolean;
};

const initialState: State = {
  environments: environmentsAdapter.getInitialState(),
  isListEnvironmentsWaiting: false,
  isGetEnvironmentWaiting: false,
  isCreateEnvironmentWaiting: false,
  isUpdateEnvironmentWaiting: false,
  isDeleteEnvironmentWaiting: false,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // listEnvironments
      .addCase(asyncActions.listEnvironments.pending, (state) => {
        state.isListEnvironmentsWaiting = true;
      })
      .addCase(asyncActions.listEnvironments.fulfilled, (state, action) => {
        state.isListEnvironmentsWaiting = false;
        environmentsAdapter.setAll(
          state.environments,
          action.payload.environments,
        );
      })
      .addCase(asyncActions.listEnvironments.rejected, (state) => {
        state.isListEnvironmentsWaiting = false;
      })

      // getEnvironment
      .addCase(asyncActions.getEnvironment.pending, (state) => {
        state.isGetEnvironmentWaiting = true;
      })
      .addCase(asyncActions.getEnvironment.fulfilled, (state, action) => {
        state.isGetEnvironmentWaiting = false;

        if (state.environments.entities[action.payload.environment.id]) {
          environmentsAdapter.updateOne(state.environments, {
            id: action.payload.environment.id,
            changes: action.payload.environment,
          });
        } else {
          environmentsAdapter.addOne(
            state.environments,
            action.payload.environment,
          );
        }
      })
      .addCase(asyncActions.getEnvironment.rejected, (state) => {
        state.isGetEnvironmentWaiting = false;
      })

      // createEnvironment
      .addCase(asyncActions.createEnvironment.pending, (state) => {
        state.isCreateEnvironmentWaiting = true;
      })
      .addCase(asyncActions.createEnvironment.fulfilled, (state, action) => {
        state.isCreateEnvironmentWaiting = false;
        environmentsAdapter.addOne(
          state.environments,
          action.payload.environment,
        );
      })
      .addCase(asyncActions.createEnvironment.rejected, (state) => {
        state.isCreateEnvironmentWaiting = false;
      })

      // updateEnvironment
      .addCase(asyncActions.updateEnvironment.pending, (state) => {
        state.isUpdateEnvironmentWaiting = true;
      })
      .addCase(asyncActions.updateEnvironment.fulfilled, (state, action) => {
        state.isUpdateEnvironmentWaiting = false;
        environmentsAdapter.updateOne(state.environments, {
          id: action.payload.environment.id,
          changes: action.payload.environment,
        });
      })
      .addCase(asyncActions.updateEnvironment.rejected, (state) => {
        state.isUpdateEnvironmentWaiting = false;
      })

      // deleteEnvironment
      .addCase(asyncActions.deleteEnvironment.pending, (state) => {
        state.isDeleteEnvironmentWaiting = true;
      })
      .addCase(asyncActions.deleteEnvironment.fulfilled, (state, action) => {
        state.isDeleteEnvironmentWaiting = false;
        environmentsAdapter.removeOne(
          state.environments,
          action.payload.environment.id,
        );
      })
      .addCase(asyncActions.deleteEnvironment.rejected, (state) => {
        state.isDeleteEnvironmentWaiting = false;
      });
  },
  initialState,
  name: 'environments',
  reducers: {},
});
