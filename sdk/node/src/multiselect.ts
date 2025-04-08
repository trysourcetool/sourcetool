import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  MultiSelectState,
  MultiSelectValue,
  WidgetTypeMultiSelect,
} from './internal/session/state/multiselect';
import { MultiSelectOptions } from './internal/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  MultiSelect as MultiSelectProto,
  MultiSelectSchema,
  WidgetSchema,
} from '@trysourcetool/proto/widget/v1/widget';
import { RenderWidgetSchema } from '@trysourcetool/proto/websocket/v1/message';
/**
 * MultiSelect component options
 */
export interface MultiSelectComponentOptions {
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
 * MultiSelect component class
 */
export class MultiSelect {
  /**
   * Set the multiselect options
   * @param options MultiSelect options
   * @returns MultiSelect options
   */
  static options(...options: string[]): MultiSelectComponentOptions {
    return { options };
  }

  /**
   * Set the default values
   * @param values Default values
   * @returns MultiSelect options
   */
  static defaultValue(...values: string[]): MultiSelectComponentOptions {
    return { defaultValue: values };
  }

  /**
   * Set the placeholder text
   * @param placeholder Placeholder text
   * @returns MultiSelect options
   */
  static placeholder(placeholder: string): MultiSelectComponentOptions {
    return { placeholder };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns MultiSelect options
   */
  static required(required: boolean): MultiSelectComponentOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns MultiSelect options
   */
  static disabled(disabled: boolean): MultiSelectComponentOptions {
    return { disabled };
  }

  /**
   * Set the format function for option labels
   * @param formatFunc Format function
   * @returns MultiSelect options
   */
  static formatFunc(
    formatFunc: (value: string, index: number) => string,
  ): MultiSelectComponentOptions {
    return { formatFunc };
  }
}

/**
 * Add a multiselect to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options MultiSelect options
 * @returns The selected values
 */
export function multiSelect(
  builder: UIBuilder,
  label: string,
  options: MultiSelectComponentOptions = {},
): MultiSelectValue | null {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return null;
  }

  const multiSelectOpts: MultiSelectOptions = {
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
  const widgetID = builder.generatePageID(WidgetTypeMultiSelect, path);

  let multiSelectState = session.state.getMultiSelect(widgetID);
  const formatFunc = multiSelectOpts.formatFunc || ((v: string) => v);
  if (!multiSelectState) {
    multiSelectState = new MultiSelectState(
      widgetID,
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
  session.state.set(widgetID, multiSelectState);

  const multiSelectProto = convertStateToMultiSelectProto(
    multiSelectState as MultiSelectState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
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
