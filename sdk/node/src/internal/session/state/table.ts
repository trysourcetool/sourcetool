import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeTable: WidgetType = 'table';

export interface TableStateValueSelection {
  row: number;
  rows: number[];
}

export interface TableStateValue {
  selection?: TableStateValueSelection;
}

export class TableState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public data: any = null,
    public header: string = '',
    public description: string = '',
    public height: number | null = null,
    public columnOrder: string[] = [],
    public onSelect: string = 'ignore',
    public rowSelection: string = 'single',
    public value: TableStateValue = {},
  ) {
    this.type = WidgetTypeTable;
  }

  getType(): WidgetType {
    return WidgetTypeTable;
  }

  public type: WidgetType = WidgetTypeTable;
}

export interface TableSelection {
  row: number;
  rows: number[];
}

export interface TableValue {
  selection?: TableSelection;
}

export enum SelectionBehavior {
  Ignore = 'ignore',
  Rerun = 'rerun',
}

export enum SelectionMode {
  Single = 'single',
  Multiple = 'multiple',
}
