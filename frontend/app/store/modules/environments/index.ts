import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';
import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import { errorStore } from '../error';
import type { ErrorResponse } from '@/api/instance';
import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import type { Environment } from '@/api/modules/environments';

// =============================================
// asyncActions
// =============================================
const listEnvironments = createAsyncThunk(
  'environments/listEnvironments',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.environments.listEnvironments();

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const getEnvironment = createAsyncThunk(
  'environments/getEnvironment',
  async (params: { environmentId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.environments.getEnvironment(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const createEnvironment = createAsyncThunk(
  'environments/createEnvironment',
  async (
    params: { data: { color: string; name: string; slug: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.environments.createEnvironment(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const updateEnvironment = createAsyncThunk(
  'environments/updateEnvironment',
  async (
    params: { environmentId: string; data: { color: string; name: string } },
    { dispatch, rejectWithValue },
  ) => {
    try {
      const res = await api.environments.updateEnvironment(params);

      return res;
    } catch (error: any) {
      dispatch(errorStore.asyncActions.handleError(error));
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const deleteEnvironment = createAsyncThunk(
  'environments/deleteEnvironment',
  async (params: { environmentId: string }, { dispatch, rejectWithValue }) => {
    try {
      const res = await api.environments.deleteEnvironment(params);

      return res;
    } catch (error: any) {
      console.error({ error });
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
      .addCase(listEnvironments.pending, (state) => {
        state.isListEnvironmentsWaiting = true;
      })
      .addCase(listEnvironments.fulfilled, (state, action) => {
        state.isListEnvironmentsWaiting = false;
        environmentsAdapter.setAll(
          state.environments,
          action.payload.environments,
        );
      })
      .addCase(listEnvironments.rejected, (state) => {
        state.isListEnvironmentsWaiting = false;
      })

      // getEnvironment
      .addCase(getEnvironment.pending, (state) => {
        state.isGetEnvironmentWaiting = true;
      })
      .addCase(getEnvironment.fulfilled, (state, action) => {
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
      .addCase(getEnvironment.rejected, (state) => {
        state.isGetEnvironmentWaiting = false;
      })

      // createEnvironment
      .addCase(createEnvironment.pending, (state) => {
        state.isCreateEnvironmentWaiting = true;
      })
      .addCase(createEnvironment.fulfilled, (state, action) => {
        state.isCreateEnvironmentWaiting = false;
        environmentsAdapter.addOne(
          state.environments,
          action.payload.environment,
        );
      })
      .addCase(createEnvironment.rejected, (state) => {
        state.isCreateEnvironmentWaiting = false;
      })

      // updateEnvironment
      .addCase(updateEnvironment.pending, (state) => {
        state.isUpdateEnvironmentWaiting = true;
      })
      .addCase(updateEnvironment.fulfilled, (state, action) => {
        state.isUpdateEnvironmentWaiting = false;
        environmentsAdapter.updateOne(state.environments, {
          id: action.payload.environment.id,
          changes: action.payload.environment,
        });
      })
      .addCase(updateEnvironment.rejected, (state) => {
        state.isUpdateEnvironmentWaiting = false;
      })

      // deleteEnvironment
      .addCase(deleteEnvironment.pending, (state) => {
        state.isDeleteEnvironmentWaiting = true;
      })
      .addCase(deleteEnvironment.fulfilled, (state, action) => {
        state.isDeleteEnvironmentWaiting = false;
        environmentsAdapter.removeOne(
          state.environments,
          action.payload.environment.id,
        );
      })
      .addCase(deleteEnvironment.rejected, (state) => {
        state.isDeleteEnvironmentWaiting = false;
      });
  },
  initialState,
  name: 'environments',
  reducers: {},
});

// =============================================
// selectors
// =============================================
const getEnvironmentIds = createSelector(
  (state: RootState) => state.environments,
  (values) => values.environments.ids,
);

const getEnvironmentEntities = createSelector(
  (state: RootState) => state.environments,
  (values) => values.environments.entities,
);

const getEnvironments = createSelector(
  (state: RootState) => state.environments,
  (values) =>
    values.environments.ids.map((id) => values.environments.entities[id]),
);

const getEnvironmentValue = createSelector(
  (state: RootState, environmentId: string) =>
    state.environments.environments.entities[environmentId],
  (values) => values,
);

// =============================================
// exports
// =============================================

export const environmentsStore = {
  actions: slice.actions,
  asyncActions: {
    listEnvironments,
    getEnvironment,
    createEnvironment,
    updateEnvironment,
    deleteEnvironment,
  },
  reducer: slice.reducer,
  selector: {
    getEnvironmentIds,
    getEnvironmentEntities,
    getEnvironments,
    getEnvironment: getEnvironmentValue,
  },
};

export type EnvironmentsState = State;
