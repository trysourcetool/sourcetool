import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  NumberInputState,
  WidgetTypeNumberInput,
} from './internal/session/state/numberinput';
import { NumberInputOptions } from './internal/options';

/**
 * NumberInput component options
 */
export interface NumberInputComponentOptions {
  /**
   * Placeholder text
   */
  placeholder?: string;

  /**
   * Default value
   */
  defaultValue?: number;

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
   * Maximum value
   */
  maxValue?: number;

  /**
   * Minimum value
   */
  minValue?: number;
}

/**
 * NumberInput component class
 */
export class NumberInput {
  /**
   * Set the placeholder text
   * @param placeholder Placeholder text
   * @returns NumberInput options
   */
  static placeholder(placeholder: string): NumberInputComponentOptions {
    return { placeholder };
  }

  /**
   * Set the default value
   * @param value Default value
   * @returns NumberInput options
   */
  static defaultValue(value: number): NumberInputComponentOptions {
    return { defaultValue: value };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns NumberInput options
   */
  static required(required: boolean): NumberInputComponentOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns NumberInput options
   */
  static disabled(disabled: boolean): NumberInputComponentOptions {
    return { disabled };
  }

  /**
   * Set the maximum value
   * @param value Maximum value
   * @returns NumberInput options
   */
  static maxValue(value: number): NumberInputComponentOptions {
    return { maxValue: value };
  }

  /**
   * Set the minimum value
   * @param value Minimum value
   * @returns NumberInput options
   */
  static minValue(value: number): NumberInputComponentOptions {
    return { minValue: value };
  }
}

/**
 * Add a number input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options NumberInput options
 * @returns The input value
 */
export function numberInput(
  builder: UIBuilder,
  label: string,
  options: NumberInputComponentOptions = {},
): number | null {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return null;
  }

  const numberInputOpts: NumberInputOptions = {
    label,
    placeholder: options.placeholder || '',
    defaultValue:
      options.defaultValue !== undefined ? options.defaultValue : null,
    required: options.required || false,
    disabled: options.disabled || false,
    maxValue: options.maxValue !== undefined ? options.maxValue : null,
    minValue: options.minValue !== undefined ? options.minValue : null,
  };

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeNumberInput, path);

  let numberInputState = session.state.getNumberInput(widgetID);
  if (!numberInputState) {
    numberInputState = new NumberInputState(
      widgetID,
      numberInputOpts.defaultValue,
      numberInputOpts.label,
      numberInputOpts.placeholder,
      numberInputOpts.defaultValue,
      numberInputOpts.required,
      numberInputOpts.disabled,
      numberInputOpts.maxValue,
      numberInputOpts.minValue,
    );
  } else {
    numberInputState.label = numberInputOpts.label;
    numberInputState.placeholder = numberInputOpts.placeholder;
    numberInputState.defaultValue = numberInputOpts.defaultValue;
    numberInputState.required = numberInputOpts.required;
    numberInputState.disabled = numberInputOpts.disabled;
    numberInputState.maxValue = numberInputOpts.maxValue;
    numberInputState.minValue = numberInputOpts.minValue;
    session.state.set(widgetID, numberInputState);
  }

  const numberInputProto = convertStateToNumberInputProto(
    numberInputState as NumberInputState,
  );
  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'NumberInput',
      numberInput: numberInputProto,
    },
  });

  cursor.next();

  return numberInputState.value;
}

/**
 * Convert number input state to proto
 * @param state Number input state
 * @returns Number input proto
 */
function convertStateToNumberInputProto(state: NumberInputState): any {
  return {
    value: state.value,
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
    maxValue: state.maxValue,
    minValue: state.minValue,
  };
}

/**
 * Convert number input proto to state
 * @param id Widget ID
 * @param data Number input proto
 * @returns Number input state
 */
export function convertNumberInputProtoToState(
  id: string,
  data: any,
): NumberInputState | null {
  if (!data) {
    return null;
  }

  return new NumberInputState(
    id,
    data.value,
    data.label,
    data.placeholder,
    data.defaultValue,
    data.required,
    data.disabled,
    data.maxValue,
    data.minValue,
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
