import * as selector from './selector';
import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as UsersState } from './slice';

export const usersStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
  selector,
};
