import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import {
  MultiSelectState,
  MultiSelectValue,
  WidgetTypeMultiSelect,
} from '../../session/state/multiselect';
import { MultiSelectInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  MultiSelect as MultiSelectProto,
  MultiSelectSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * MultiSelect options
 */
export interface MultiSelectOptions {
  /**
   * MultiSelect options
   */
  options?: string[];

  /**
   * Default values
   */
  defaultValue?: string[];

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
 * Add a multiselect to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options MultiSelect options
 * @returns The selected values
 */
export function multiSelect(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: MultiSelectOptions = {},
): MultiSelectValue | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const multiSelectOpts: MultiSelectInternalOptions = {
    label,
    options: options.options || [],
    defaultValue: options.defaultValue || null,
    placeholder: options.placeholder || '',
    required: options.required || false,
    disabled: options.disabled || false,
    formatFunc: options.formatFunc || ((v: string) => v),
  };

  // Find default value indexes
  const defaultVal: number[] = [];
  if (multiSelectOpts.defaultValue && multiSelectOpts.options.length > 0) {
    for (const defaultOption of multiSelectOpts.defaultValue) {
      for (let i = 0; i < multiSelectOpts.options.length; i++) {
        if (multiSelectOpts.options[i] === defaultOption) {
          defaultVal.push(i);
          break;
        }
      }
    }
  }

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(
    page.id,
    WidgetTypeMultiSelect,
    path,
  );

  let multiSelectState = session.state.getMultiSelect(widgetId);
  const formatFunc = multiSelectOpts.formatFunc || ((v: string) => v);
  if (!multiSelectState) {
    multiSelectState = new MultiSelectState(
      widgetId,
      defaultVal,
      multiSelectOpts.label,
      multiSelectOpts.options.map(formatFunc),
      multiSelectOpts.placeholder,
      defaultVal,
      multiSelectOpts.required,
      multiSelectOpts.disabled,
    );
  } else {
    const displayVals = multiSelectOpts.options.map((v, i) => formatFunc(v, i));

    multiSelectState.label = multiSelectOpts.label;
    multiSelectState.options = displayVals;
    multiSelectState.placeholder = multiSelectOpts.placeholder;
    multiSelectState.defaultValue = defaultVal;
    multiSelectState.required = multiSelectOpts.required;
    multiSelectState.disabled = multiSelectOpts.disabled;
  }
  session.state.set(widgetId, multiSelectState);

  const multiSelectProto = convertStateToMultiSelectProto(
    multiSelectState as MultiSelectState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'multiSelect',
        value: multiSelectProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  // Return the selected values
  let value: MultiSelectValue | null = null;
  if (multiSelectState.value.length > 0 && multiSelectOpts.options.length > 0) {
    value = {
      values: multiSelectState.value.map(
        (idx: number) => multiSelectOpts.options[idx],
      ),
      indexes: multiSelectState.value.map((idx: number) => idx),
    };
  }

  return value;
}

/**
 * Convert multiselect state to proto
 * @param state MultiSelect state
 * @returns MultiSelect proto
 */
export function convertStateToMultiSelectProto(
  state: MultiSelectState,
): MultiSelectProto {
  return fromJson(MultiSelectSchema, {
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
 * Convert multiselect proto to state
 * @param id Widget ID
 * @param data MultiSelect proto
 * @returns MultiSelect state
 */
export function convertMultiSelectProtoToState(
  id: string,
  data: MultiSelectProto | null,
): MultiSelectState | null {
  if (!data) {
    return null;
  }

  const d = toJson(MultiSelectSchema, data);

  return new MultiSelectState(
    id,
    d.value || [],
    d.label,
    d.options || [],
    d.placeholder,
    d.defaultValue || [],
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
