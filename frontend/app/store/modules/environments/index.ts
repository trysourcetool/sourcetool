import * as selector from './selector';
import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as EnvironmentsState } from './slice';

export const environmentsStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
  selector,
};
