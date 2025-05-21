import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeColumns = 'columns' as const;

export class ColumnsState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public columns: number = 1,
  ) {
    this.type = WidgetTypeColumns;
  }

  getType(): 'columns' {
    return WidgetTypeColumns;
  }

  public type: 'columns' = WidgetTypeColumns;
}
