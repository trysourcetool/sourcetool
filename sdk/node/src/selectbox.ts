import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  SelectboxState,
  SelectboxValue,
  WidgetTypeSelectbox,
} from './internal/session/state/selectbox';
import { SelectboxOptions } from './internal/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  Selectbox as SelectboxProto,
  SelectboxSchema,
  WidgetSchema,
} from '@trysourcetool/proto/widget/v1/widget';
import { RenderWidgetSchema } from '@trysourcetool/proto/websocket/v1/message';
/**
 * Selectbox component options
 */
export interface SelectboxComponentOptions {
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
 * Selectbox component class
 */
export class Selectbox {
  /**
   * Set the selectbox options
   * @param options Selectbox options
   * @returns Selectbox options
   */
  static options(...options: string[]): SelectboxComponentOptions {
    return { options };
  }

  /**
   * Set the default value
   * @param value Default value
   * @returns Selectbox options
   */
  static defaultValue(value: string): SelectboxComponentOptions {
    return { defaultValue: value };
  }

  /**
   * Set the placeholder text
   * @param placeholder Placeholder text
   * @returns Selectbox options
   */
  static placeholder(placeholder: string): SelectboxComponentOptions {
    return { placeholder };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns Selectbox options
   */
  static required(required: boolean): SelectboxComponentOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns Selectbox options
   */
  static disabled(disabled: boolean): SelectboxComponentOptions {
    return { disabled };
  }

  /**
   * Set the format function for option labels
   * @param formatFunc Format function
   * @returns Selectbox options
   */
  static formatFunc(
    formatFunc: (value: string, index: number) => string,
  ): SelectboxComponentOptions {
    return { formatFunc };
  }
}

/**
 * Add a selectbox to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options Selectbox options
 * @returns The selected value
 */
export function selectbox(
  builder: UIBuilder,
  label: string,
  options: SelectboxComponentOptions = {},
): SelectboxValue | null {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return null;
  }

  const selectboxOpts: SelectboxOptions = {
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
  const widgetID = builder.generatePageID(WidgetTypeSelectbox, path);

  let selectboxState = session.state.getSelectbox(widgetID);
  const formatFunc = selectboxOpts.formatFunc || ((v: string) => v);

  if (!selectboxState) {
    selectboxState = new SelectboxState(
      widgetID,
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

  session.state.set(widgetID, selectboxState);

  const selectboxProto = convertStateToSelectboxProto(
    selectboxState as SelectboxState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
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
