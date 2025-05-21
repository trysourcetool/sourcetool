import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import {
  SelectboxState,
  SelectboxValue,
  WidgetTypeSelectbox,
} from '../../session/state/selectbox';
import { SelectboxInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  Selectbox as SelectboxProto,
  SelectboxSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * Selectbox options
 */
export interface SelectboxOptions {
  /**
   * Selectbox options
   */
  options?: string[];

  /**
   * Default value
   */
  defaultValue?: string;

  /**
   * Placeholder text
   */
  placeholder?: string;

  /**
   * Whether the input is required
   * @default false
   */
  required?: boolean;

  /**
   * Whether the input is disabled
   * @default false
   */
  disabled?: boolean;

  /**
   * Format function for option labels
   */
  formatFunc?: (value: string, index: number) => string;
}

/**
 * Add a selectbox to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options Selectbox options
 * @returns The selected value
 */
export function selectbox(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: SelectboxOptions = {},
): SelectboxValue | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const selectboxOpts: SelectboxInternalOptions = {
    label,
    options: options.options || [],
    defaultValue: options.defaultValue || null,
    placeholder: options.placeholder || '',
    required: options.required || false,
    disabled: options.disabled || false,
    formatFunc: options.formatFunc || ((v: string) => v),
  };

  // Find default value index
  let defaultVal: number | null = null;
  if (options.defaultValue && selectboxOpts.options.length > 0) {
    for (let i = 0; i < selectboxOpts.options.length; i++) {
      if (selectboxOpts.options[i] === options.defaultValue) {
        defaultVal = i;
        break;
      }
    }
  }

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(page.id, WidgetTypeSelectbox, path);

  let selectboxState = session.state.getSelectbox(widgetId);
  const formatFunc = selectboxOpts.formatFunc || ((v: string) => v);

  if (!selectboxState) {
    selectboxState = new SelectboxState(
      widgetId,
      defaultVal,
      selectboxOpts.label,
      selectboxOpts.options.map(formatFunc),
      selectboxOpts.placeholder,
      defaultVal,
      selectboxOpts.required,
      selectboxOpts.disabled,
    );
  } else {
    const displayVals = selectboxOpts.options.map((v, i) => formatFunc(v, i));

    selectboxState.label = selectboxOpts.label;
    selectboxState.options = displayVals;
    selectboxState.placeholder = selectboxOpts.placeholder;
    selectboxState.defaultValue = defaultVal;
    selectboxState.required = selectboxOpts.required;
    selectboxState.disabled = selectboxOpts.disabled;
  }

  session.state.set(widgetId, selectboxState);

  const selectboxProto = convertStateToSelectboxProto(
    selectboxState as SelectboxState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'selectbox',
        value: selectboxProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  // Return the selected value
  let value: SelectboxValue | null = null;
  if (selectboxState.value !== null && selectboxOpts.options.length > 0) {
    value = {
      value: selectboxOpts.options[selectboxState.value],
      index: selectboxState.value,
    };
  }

  return value;
}

/**
 * Convert selectbox state to proto
 * @param state Selectbox state
 * @returns Selectbox proto
 */
export function convertStateToSelectboxProto(
  state: SelectboxState,
): SelectboxProto {
  return fromJson(SelectboxSchema, {
    label: state.label,
    value: state.value,
    options: state.options,
    placeholder: state.placeholder,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
  });
}

/**
 * Convert selectbox proto to state
 * @param id Widget ID
 * @param data Selectbox proto
 * @returns Selectbox state
 */
export function convertSelectboxProtoToState(
  id: string,
  data: SelectboxProto | null,
): SelectboxState | null {
  if (!data) {
    return null;
  }

  const d = toJson(SelectboxSchema, data);

  return new SelectboxState(
    id,
    d.value as number,
    d.label,
    d.options,
    d.placeholder,
    d.defaultValue as number,
    d.required,
    d.disabled,
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
