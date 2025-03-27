import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  TableState,
  TableValue,
  WidgetTypeTable,
  SelectionBehavior,
  SelectionMode,
  TableStateValue,
  TableStateValueSelection,
  TableSelection,
} from './internal/session/state/table';
import { TableOptions } from './internal/options';

/**
 * Table component options
 */
export interface TableComponentOptions {
  /**
   * Table header
   */
  header?: string;

  /**
   * Table description
   */
  description?: string;

  /**
   * Table height
   */
  height?: number;

  /**
   * Column order
   */
  columnOrder?: string[];

  /**
   * Selection behavior
   * @default SelectionBehavior.Ignore
   */
  onSelect?: SelectionBehavior;

  /**
   * Row selection mode
   * @default SelectionMode.Single
   */
  rowSelection?: SelectionMode;
}

/**
 * Table component class
 */
export class Table {
  /**
   * Set the table header
   * @param header Table header
   * @returns Table options
   */
  static header(header: string): TableComponentOptions {
    return { header };
  }

  /**
   * Set the table description
   * @param description Table description
   * @returns Table options
   */
  static description(description: string): TableComponentOptions {
    return { description };
  }

  /**
   * Set the table height
   * @param height Table height
   * @returns Table options
   */
  static height(height: number): TableComponentOptions {
    return { height };
  }

  /**
   * Set the column order
   * @param columns Column order
   * @returns Table options
   */
  static columnOrder(...columns: string[]): TableComponentOptions {
    return { columnOrder: columns };
  }

  /**
   * Set the selection behavior
   * @param behavior Selection behavior
   * @returns Table options
   */
  static onSelect(behavior: SelectionBehavior): TableComponentOptions {
    return { onSelect: behavior };
  }

  /**
   * Set the row selection mode
   * @param mode Row selection mode
   * @returns Table options
   */
  static rowSelection(mode: SelectionMode): TableComponentOptions {
    return { rowSelection: mode };
  }
}

/**
 * Add a table to the UI
 * @param builder The UI builder
 * @param data The table data
 * @param options Table options
 * @returns The table value
 */
export function table(
  builder: UIBuilder,
  data: any,
  options: TableComponentOptions = {},
): TableValue {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return {};
  }

  const tableOpts: TableOptions = {
    header: options.header || '',
    description: options.description || '',
    height: options.height !== undefined ? options.height : null,
    columnOrder: options.columnOrder || [],
    onSelect:
      options.onSelect !== undefined
        ? options.onSelect.toString()
        : SelectionBehavior.Ignore.toString(),
    rowSelection:
      options.rowSelection !== undefined
        ? options.rowSelection.toString()
        : SelectionMode.Single.toString(),
  };

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeTable, path);

  let tableState = session.state.getTable(widgetID);
  if (!tableState) {
    tableState = new TableState(
      widgetID,
      data,
      tableOpts.header,
      tableOpts.description,
      tableOpts.height,
      tableOpts.columnOrder,
      tableOpts.onSelect,
      tableOpts.rowSelection,
      {},
    );
  } else {
    tableState.data = data;
    tableState.header = tableOpts.header;
    tableState.description = tableOpts.description;
    tableState.height = tableOpts.height;
    tableState.columnOrder = tableOpts.columnOrder;
    tableState.onSelect = tableOpts.onSelect;
    tableState.rowSelection = tableOpts.rowSelection;
    session.state.set(widgetID, tableState);
  }

  const tableProto = convertStateToTableProto(tableState as TableState);
  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'Table',
      table: tableProto,
    },
  });

  cursor.next();

  // Return the table value
  const value: TableValue = {};
  if (tableState.value.selection) {
    const rows: number[] = [];
    for (const row of tableState.value.selection.rows) {
      rows.push(row);
    }
    value.selection = {
      row: tableState.value.selection.row,
      rows: rows,
    };
  }

  return value;
}

/**
 * Convert table state to proto
 * @param state Table state
 * @returns Table proto
 */
function convertStateToTableProto(state: TableState): any {
  try {
    const dataBytes = JSON.stringify(state.data);
    const tableValue: any = {};

    if (state.value.selection) {
      tableValue.selection = {
        row: state.value.selection.row,
        rows: state.value.selection.rows,
      };
    }

    return {
      data: dataBytes,
      header: state.header,
      description: state.description,
      height: state.height,
      columnOrder: state.columnOrder,
      onSelect: state.onSelect,
      rowSelection: state.rowSelection,
      value: tableValue,
    };
  } catch (error) {
    console.error('Error converting table state to proto:', error);
    return {};
  }
}

/**
 * Convert table proto to state
 * @param id Widget ID
 * @param data Table proto
 * @returns Table state
 */
export function convertTableProtoToState(
  id: string,
  data: any,
): TableState | null {
  if (!data) {
    return null;
  }

  try {
    const tableData =
      typeof data.data === 'string' ? JSON.parse(data.data) : data.data;
    const tableValue: TableStateValue = {};

    if (data.value && data.value.selection) {
      const selection: TableStateValueSelection = {
        row: data.value.selection.row,
        rows: data.value.selection.rows || [],
      };
      tableValue.selection = selection;
    }

    return new TableState(
      id,
      tableData,
      data.header,
      data.description,
      data.height,
      data.columnOrder || [],
      data.onSelect,
      data.rowSelection,
      tableValue,
    );
  } catch (error) {
    console.error('Error converting table proto to state:', error);
    return null;
  }
}

/**
 * Convert path to int32 array
 * @param path Path
 * @returns Int32 array
 */
export function convertPathToInt32Array(path: number[]): number[] {
  return path.map((v) => v);
}

// Re-export types from table state
export { SelectionBehavior, SelectionMode, TableValue, TableSelection };
