import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeColumnItem: WidgetType = 'columnItem';

export class ColumnItemState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public weight: number = 1,
  ) {}

  getType(): WidgetType {
    return WidgetTypeColumnItem;
  }
}
