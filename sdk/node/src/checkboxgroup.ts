import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  CheckboxGroupState,
  CheckboxGroupValue,
  WidgetTypeCheckboxGroup,
} from './internal/session/state/checkboxgroup';
import { CheckboxGroupOptions } from './internal/options';

/**
 * CheckboxGroup component options
 */
export interface CheckboxGroupComponentOptions {
  /**
   * CheckboxGroup options
   */
  options?: string[];

  /**
   * Default values
   */
  defaultValue?: string[];

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
 * CheckboxGroup component class
 */
export class CheckboxGroup {
  /**
   * Set the checkboxgroup options
   * @param options CheckboxGroup options
   * @returns CheckboxGroup options
   */
  static options(...options: string[]): CheckboxGroupComponentOptions {
    return { options };
  }

  /**
   * Set the default values
   * @param values Default values
   * @returns CheckboxGroup options
   */
  static defaultValue(...values: string[]): CheckboxGroupComponentOptions {
    return { defaultValue: values };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns CheckboxGroup options
   */
  static required(required: boolean): CheckboxGroupComponentOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns CheckboxGroup options
   */
  static disabled(disabled: boolean): CheckboxGroupComponentOptions {
    return { disabled };
  }

  /**
   * Set the format function for option labels
   * @param formatFunc Format function
   * @returns CheckboxGroup options
   */
  static formatFunc(
    formatFunc: (value: string, index: number) => string,
  ): CheckboxGroupComponentOptions {
    return { formatFunc };
  }
}

/**
 * Add a checkboxgroup to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options CheckboxGroup options
 * @returns The selected values
 */
export function checkboxGroup(
  builder: UIBuilder,
  label: string,
  options: CheckboxGroupComponentOptions = {},
): CheckboxGroupValue | null {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return null;
  }

  const checkboxGroupOpts: CheckboxGroupOptions = {
    label,
    options: options.options || [],
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    formatFunc: options.formatFunc || ((v: string, i: number) => v),
  };

  // Find default value indexes
  const defaultVal: number[] = [];
  if (checkboxGroupOpts.defaultValue && checkboxGroupOpts.options.length > 0) {
    for (const defaultOption of checkboxGroupOpts.defaultValue) {
      for (let i = 0; i < checkboxGroupOpts.options.length; i++) {
        if (checkboxGroupOpts.options[i] === defaultOption) {
          defaultVal.push(i);
          break;
        }
      }
    }
  }

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeCheckboxGroup, path);

  let checkboxGroupState = session.state.getCheckboxGroup(widgetID);
  if (!checkboxGroupState) {
    checkboxGroupState = new CheckboxGroupState(
      widgetID,
      defaultVal,
      checkboxGroupOpts.label,
      [],
      defaultVal,
      checkboxGroupOpts.required,
      checkboxGroupOpts.disabled,
    );
  }

  // Apply format function to options
  const formatFunc =
    checkboxGroupOpts.formatFunc || ((v: string, i: number) => v);
  const displayVals = checkboxGroupOpts.options.map((v, i) => formatFunc(v, i));

  checkboxGroupState.label = checkboxGroupOpts.label;
  checkboxGroupState.options = displayVals;
  checkboxGroupState.defaultValue = defaultVal;
  checkboxGroupState.required = checkboxGroupOpts.required;
  checkboxGroupState.disabled = checkboxGroupOpts.disabled;
  session.state.set(widgetID, checkboxGroupState);

  const checkboxGroupProto =
    convertStateToCheckboxGroupProto(checkboxGroupState);
  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'CheckboxGroup',
      checkboxGroup: checkboxGroupProto,
    },
  });

  cursor.next();

  // Return the selected values
  let value: CheckboxGroupValue | null = null;
  if (
    checkboxGroupState.value.length > 0 &&
    checkboxGroupOpts.options.length > 0
  ) {
    value = {
      values: checkboxGroupState.value.map(
        (idx: number) => checkboxGroupOpts.options[idx],
      ),
      indexes: checkboxGroupState.value.map((idx: number) => idx),
    };
  }

  return value;
}

/**
 * Convert checkboxgroup state to proto
 * @param state CheckboxGroup state
 * @returns CheckboxGroup proto
 */
function convertStateToCheckboxGroupProto(state: CheckboxGroupState): any {
  return {
    label: state.label,
    value: state.value,
    options: state.options,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
  };
}

/**
 * Convert checkboxgroup proto to state
 * @param id Widget ID
 * @param data CheckboxGroup proto
 * @returns CheckboxGroup state
 */
export function convertCheckboxGroupProtoToState(
  id: string,
  data: any,
): CheckboxGroupState | null {
  if (!data) {
    return null;
  }

  return new CheckboxGroupState(
    id,
    data.value || [],
    data.label,
    data.options || [],
    data.defaultValue || [],
    data.required,
    data.disabled,
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
