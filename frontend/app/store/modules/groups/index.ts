import * as selector from './selector';
import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as GroupsState } from './slice';

export const groupsStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
  selector,
};
