import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from '../';
import {
  NumberInputState,
  WidgetTypeNumberInput,
} from '../../session/state/numberinput';
import { NumberInputOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  NumberInput as NumberInputProto,
  NumberInputSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
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
  }

  session.state.set(widgetID, numberInputState);

  const numberInputProto = convertStateToNumberInputProto(
    numberInputState as NumberInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
      type: {
        case: 'numberInput',
        value: numberInputProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return numberInputState.value;
}

/**
 * Convert number input state to proto
 * @param state Number input state
 * @returns Number input proto
 */
export function convertStateToNumberInputProto(
  state: NumberInputState,
): NumberInputProto {
  return fromJson(NumberInputSchema, {
    value: state.value,
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
    maxValue: state.maxValue,
    minValue: state.minValue,
  });
}

/**
 * Convert number input proto to state
 * @param id Widget ID
 * @param data Number input proto
 * @returns Number input state
 */
export function convertNumberInputProtoToState(
  id: string,
  data: NumberInputProto | null,
): NumberInputState | null {
  if (!data) {
    return null;
  }

  const d = toJson(NumberInputSchema, data);

  return new NumberInputState(
    id,
    d.value as number,
    d.label,
    d.placeholder,
    d.defaultValue as number,
    d.required,
    d.disabled,
    d.maxValue as number,
    d.minValue as number,
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
