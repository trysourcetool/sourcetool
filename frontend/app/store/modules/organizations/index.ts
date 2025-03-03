import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as OrganizationsState } from './slice';

export const organizationsStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
};
