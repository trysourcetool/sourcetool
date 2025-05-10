import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeColumnItem = 'columnItem' as const;

export class ColumnItemState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public weight: number = 1,
  ) {
    this.type = WidgetTypeColumnItem;
  }

  getType(): 'columnItem' {
    return WidgetTypeColumnItem;
  }

  public type: 'columnItem' = WidgetTypeColumnItem;
}
