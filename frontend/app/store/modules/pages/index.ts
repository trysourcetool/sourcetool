import * as selector from './selector';
import * as asyncActions from './asyncActions';
import { slice } from './slice';
export type { State as PagesState } from './slice';

export const pagesStore = {
  actions: slice.actions,
  asyncActions,
  reducer: slice.reducer,
  selector,
};
