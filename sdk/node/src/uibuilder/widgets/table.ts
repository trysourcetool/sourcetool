import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import {
  TableState,
  TableValue,
  WidgetTypeTable,
  TableOnSelect,
  TableRowSelection,
  TableStateValue,
  TableStateValueSelection,
  TableSelection,
} from '../../session/state/table';
import { TableInternalOptions } from '../../types/options';
import { create, toJson } from '@bufbuild/protobuf';
import {
  Table as TableProto,
  TableSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';

/**
 * Table options
 */
export interface TableOptions {
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
   * @default TableOnSelect.Ignore
   */
  onSelect?: TableOnSelect;

  /**
   * Row selection mode
   * @default TableRowSelection.Single
   */
  rowSelection?: TableRowSelection;
}

/**
 * Add a table to the UI
 * @param builder The UI builder
 * @param data The table data
 * @param options Table options
 * @returns The table value
 */
export function table(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  data: any,
  options: TableOptions = {},
): TableValue {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return {};
  }

  const tableOpts: TableInternalOptions = {
    header: options.header || '',
    description: options.description || '',
    height: options.height !== undefined ? options.height : null,
    columnOrder: options.columnOrder || [],
    onSelect:
      options.onSelect !== undefined
        ? options.onSelect.toString()
        : TableOnSelect.Ignore.toString(),
    rowSelection:
      options.rowSelection !== undefined
        ? options.rowSelection.toString()
        : TableRowSelection.Single.toString(),
  };

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(page.id, WidgetTypeTable, path);

  let tableState = session.state.getTable(widgetId);
  if (!tableState) {
    tableState = new TableState(
      widgetId,
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
  }

  session.state.set(widgetId, tableState);

  const tableProto = convertStateToTableProto(tableState as TableState);

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'table',
        value: tableProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

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
export function convertStateToTableProto(state: TableState): TableProto {
  const tableValue: any = {};
  const dataBytes = new TextEncoder().encode(JSON.stringify(state.data));

  if (state.value.selection) {
    tableValue.selection = {
      row: state.value.selection.row,
      rows: state.value.selection.rows,
    };
  }

  return create(TableSchema, {
    data: dataBytes,
    header: state.header,
    description: state.description,
    height: state.height || undefined,
    columnOrder: state.columnOrder,
    onSelect: state.onSelect,
    rowSelection: state.rowSelection,
    value: tableValue,
  });
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

  const d = toJson(TableSchema, data);

  const tableData =
    typeof d.data === 'string' ? JSON.parse(atob(d.data)) : d.data;
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
    d.header,
    d.description,
    d.height,
    d.columnOrder || [],
    d.onSelect,
    d.rowSelection,
    tableValue,
  );
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
export { TableOnSelect, TableRowSelection, TableValue, TableSelection };
