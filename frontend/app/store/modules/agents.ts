import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../';
import { createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '@/api';
import type { ErrorResponse } from '@/api/instance';
import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import type { Agent } from '@/api/modules/agents';

// =============================================
// asyncActions
// =============================================
const listAgents = createAsyncThunk(
  'agents/listAgents',
  async ({ environmentId }: { environmentId: string }, { rejectWithValue }) => {
    try {
      const res = await api.agents.listAgents({ environmentId });
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

const getAgent = createAsyncThunk(
  'agents/getAgent',
  async ({ agentId }: { agentId: string }, { rejectWithValue }) => {
    try {
      const res = await api.agents.getAgent({ agentId });
      return res;
    } catch (error: any) {
      return rejectWithValue(error as ErrorResponse);
    }
  },
);

// =============================================
// slice
// =============================================
// =============================================
// schema

const agentsAdapter = createEntityAdapter<Agent, string>({
  selectId: (agent) => agent.id,
});

// =============================================
// State

export type State = {
  agents: EntityState<Agent, string>;
  isListAgentsWaiting: boolean;
  isGetAgentWaiting: boolean;
};

const initialState: State = {
  agents: agentsAdapter.getInitialState(),
  isListAgentsWaiting: false,
  isGetAgentWaiting: false,
};

// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // listAgents
      .addCase(listAgents.pending, (state) => {
        state.isListAgentsWaiting = true;
      })
      .addCase(listAgents.fulfilled, (state, action) => {
        state.isListAgentsWaiting = false;
        agentsAdapter.setAll(state.agents, action.payload.agents);
      })
      .addCase(listAgents.rejected, (state) => {
        state.isListAgentsWaiting = false;
      })

      // getAgent
      .addCase(getAgent.pending, (state) => {
        state.isGetAgentWaiting = true;
      })
      .addCase(getAgent.fulfilled, (state, action) => {
        state.isGetAgentWaiting = false;

        if (state.agents.entities[action.payload.agent.id]) {
          agentsAdapter.updateOne(state.agents, {
            id: action.payload.agent.id,
            changes: action.payload.agent,
          });
        } else {
          agentsAdapter.addOne(state.agents, action.payload.agent);
        }
      })
      .addCase(getAgent.rejected, (state) => {
        state.isGetAgentWaiting = false;
      });
  },
  initialState,
  name: 'agents',
  reducers: {},
});

// =============================================
// selectors
// =============================================
const getAgentIds = createSelector(
  (state: RootState) => state.agents,
  (values) => values.agents.ids,
);

const getAgentEntities = createSelector(
  (state: RootState) => state.agents,
  (values) => values.agents.entities,
);

const getAgents = createSelector(
  (state: RootState) => state.agents,
  (values) => values.agents.ids.map((id) => values.agents.entities[id]),
);

const getAgentValue = createSelector(
  (state: RootState, agentId: string) =>
    state.agents.agents.entities[agentId],
  (values) => values,
);

// =============================================
// exports
// =============================================

export const agentsStore = {
  actions: slice.actions,
  asyncActions: {
    listAgents,
    getAgent,
  },
  reducer: slice.reducer,
  selector: {
    getAgentIds,
    getAgentEntities,
    getAgents,
    getAgent: getAgentValue,
  },
};

export type AgentsState = State;