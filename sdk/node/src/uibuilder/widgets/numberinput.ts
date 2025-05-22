import { v4 as uuidv4 } from 'uuid';
import { Cursor, generateWidgetId } from '../';
import {
  NumberInputState,
  WidgetTypeNumberInput,
} from '../../session/state/numberinput';
import { NumberInputInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  NumberInput as NumberInputProto,
  NumberInputSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * NumberInput options
 */
export interface NumberInputOptions {
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
 * Add a number input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options NumberInput options
 * @returns The input value
 */
export function numberInput(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: NumberInputOptions = {},
): number | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const numberInputOpts: NumberInputInternalOptions = {
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
  const widgetId = generateWidgetId(page.id, WidgetTypeNumberInput, path);

  let numberInputState = session.state.getNumberInput(widgetId);
  if (!numberInputState) {
    numberInputState = new NumberInputState(
      widgetId,
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

  session.state.set(widgetId, numberInputState);

  const numberInputProto = convertStateToNumberInputProto(
    numberInputState as NumberInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
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
