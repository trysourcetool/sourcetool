import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as HostInstancesState } from './slice';

export const hostInstancesStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
};
