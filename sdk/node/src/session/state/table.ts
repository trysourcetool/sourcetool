import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeTable = 'table' as const;

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

  getType(): 'table' {
    return WidgetTypeTable;
  }

  public type: 'table' = WidgetTypeTable;
}

export interface TableSelection {
  row: number;
  rows: number[];
}

export interface TableValue {
  selection?: TableSelection;
}

export const TableOnSelect = {
  Ignore: 'ignore',
  Rerun: 'rerun',
} as const;

export type TableOnSelect =
  (typeof TableOnSelect)[keyof typeof TableOnSelect];

export const TableRowSelection = {
  Single: 'single',
  Multiple: 'multiple',
} as const;

export type TableRowSelection =
  (typeof TableRowSelection)[keyof typeof TableRowSelection];

