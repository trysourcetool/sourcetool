import { v4 as uuidv4 } from 'uuid';
import { UIBuilder, Cursor } from './uibuilder';
import {
  ColumnsState,
  WidgetTypeColumns,
} from './internal/session/state/columns';
import {
  ColumnItemState,
  WidgetTypeColumnItem,
} from './internal/session/state/columnitem';
import { ColumnsOptions } from './internal/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  ColumnItem,
  ColumnItemSchema,
  Columns as ColumnsProto,
  ColumnsSchema,
  WidgetSchema,
} from '@trysourcetool/proto/widget/v1/widget';
import { RenderWidgetSchema } from '@trysourcetool/proto/websocket/v1/message';

/**
 * Columns component options
 */
export interface ColumnsComponentOptions {
  /**
   * Column weights
   * @description Relative weights for each column
   */
  weight?: number[];
}

/**
 * Columns component class
 */
export class Columns {
  /**
   * Set column weights
   * @param weights Column weights
   * @returns Columns options
   */
  static weight(...weights: number[]): ColumnsComponentOptions {
    return { weight: weights };
  }
}

/**
 * Add columns to the UI
 * @param builder The UI builder
 * @param cols Number of columns
 * @param options Columns options
 * @returns Array of UI builders for each column
 */
export function columns(
  builder: UIBuilder,
  cols: number,
  options: ColumnsComponentOptions = {},
): UIBuilder[] {
  if (cols < 1) {
    return [];
  }

  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return [];
  }

  const columnsOpts: ColumnsOptions = {
    columns: cols,
    weight: options.weight,
  };

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeColumns, path);

  // Calculate weights
  let weights = columnsOpts.weight || [];
  if (weights.length === 0 || weights.length !== cols) {
    weights = Array(cols).fill(1);
  }

  // Validate weights
  for (const w of weights) {
    if (w <= 0) {
      weights = Array(cols).fill(1);
      break;
    }
  }

  // Calculate total weight
  const totalWeight = weights.reduce((sum, w) => sum + w, 0);

  // Create columns state
  const columnsState = new ColumnsState(widgetID, cols);
  session.state.set(widgetID, columnsState);

  // Send columns widget
  const columnsProto = convertStateToColumnsProto(columnsState);
  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'Columns',
      columns: columnsProto,
    },
  });

  // Create builders for each column
  const builders: UIBuilder[] = [];
  for (let i = 0; i < cols; i++) {
    // Create column cursor
    const columnCursor = new Cursor();
    columnCursor.parentPath = [...path, i];

    // Create column path
    const columnPath = [...path, i];

    // Create column item state
    const columnItemID = builder.generatePageID(
      WidgetTypeColumnItem,
      columnPath,
    );
    const columnItemState = new ColumnItemState(
      columnItemID,
      weights[i] / totalWeight,
    );
    session.state.set(columnItemID, columnItemState);

    // Send column item widget
    const columnItemProto = convertStateToColumnItemProto(
      columnItemState as ColumnItemState,
    );

    const renderWidget = create(RenderWidgetSchema, {
      sessionId: session.id,
      pageId: page.id,
      path: convertPathToInt32Array(columnPath),
      widget: create(WidgetSchema, {
        id: columnItemID,
        type: {
          case: 'columnItem',
          value: columnItemProto,
        },
      }),
    });

    runtime.wsClient.enqueue(uuidv4(), renderWidget);

    // Create builder for this column
    const columnBuilder = new UIBuilder(runtime, session, page);
    columnBuilder.cursor = columnCursor;
    builders.push(columnBuilder);
  }

  cursor.next();

  return builders;
}

/**
 * Convert columns state to proto
 * @param state Columns state
 * @returns Columns proto
 */
function convertStateToColumnsProto(state: ColumnsState): ColumnsProto {
  return fromJson(ColumnsSchema, {
    columns: state.columns,
  });
}

/**
 * Convert columns proto to state
 * @param id Widget ID
 * @param data Columns proto
 * @returns Columns state
 */
export function convertColumnsProtoToState(
  id: string,
  data: ColumnsProto | null,
): ColumnsState | null {
  if (!data) {
    return null;
  }

  const d = toJson(ColumnsSchema, data);

  return new ColumnsState(id, d.columns);
}

/**
 * Convert column item state to proto
 * @param state Column item state
 * @returns Column item proto
 */
function convertStateToColumnItemProto(state: ColumnItemState): ColumnItem {
  return fromJson(ColumnItemSchema, {
    weight: state.weight,
  });
}

/**
 * Convert column item proto to state
 * @param id Widget ID
 * @param data Column item proto
 * @returns Column item state
 */
export function convertColumnItemProtoToState(
  id: string,
  data: ColumnItem | null,
): ColumnItemState | null {
  if (!data) {
    return null;
  }

  const d = toJson(ColumnItemSchema, data);

  return new ColumnItemState(id, d.weight as number);
}

/**
 * Convert path to int32 array
 * @param path Path
 * @returns Int32 array
 */
function convertPathToInt32Array(path: number[]): number[] {
  return path.map((v) => v);
}
