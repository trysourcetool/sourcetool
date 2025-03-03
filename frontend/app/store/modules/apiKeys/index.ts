import * as selector from './selector';
import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as ApiKeysState } from './slice';

export const apiKeysStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
  selector,
};
