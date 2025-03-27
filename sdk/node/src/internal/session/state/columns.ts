import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeColumns: WidgetType = 'columns';

export class ColumnsState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public columns: number = 1,
  ) {}

  getType(): WidgetType {
    return WidgetTypeColumns;
  }
}
