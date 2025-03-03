import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as ErrorState } from './slice';

export const errorStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
};
