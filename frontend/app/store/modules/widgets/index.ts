import * as selector from './selector';
import { slice } from './slice';
export type { State as WidgetsState } from './slice';

export const widgetsStore = {
  actions: slice.actions,
  reducer: slice.reducer,
  selector,
};
