import { v4 as uuidv4 } from 'uuid';
import { Cursor, generateWidgetId } from '../';
import {
  RadioState,
  RadioValue,
  WidgetTypeRadio,
} from '../../session/state/radio';
import { RadioInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import { Radio as RadioProto, RadioSchema } from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * Radio options
 */
export interface RadioOptions {
  /**
   * Radio options
   */
  options?: string[];

  /**
   * Default value
   */
  defaultValue?: string;

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
 * Add a radio input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options Radio options
 * @returns The selected radio value
 */
export function radio(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: RadioOptions = {},
): RadioValue | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const radioOpts: RadioInternalOptions = {
    label,
    options: options.options || [],
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    formatFunc: options.formatFunc || ((v: string) => v),
  };

  // Find default value index
  let defaultVal: number | null = null;
  if (options.defaultValue && radioOpts.options.length > 0) {
    for (let i = 0; i < radioOpts.options.length; i++) {
      if (radioOpts.options[i] === options.defaultValue) {
        defaultVal = i;
        break;
      }
    }
  }

  const path = cursor.getPath();
  const widgetId = generateWidgetId(page.id, WidgetTypeRadio, path);

  let radioState = session.state.getRadio(widgetId);
  const formatFunc = radioOpts.formatFunc || ((v: string) => v);

  if (!radioState) {
    radioState = new RadioState(
      widgetId,
      defaultVal,
      radioOpts.label,
      radioOpts.options.map(formatFunc),
      defaultVal,
      radioOpts.required,
      radioOpts.disabled,
    );
  } else {
    const displayVals = radioOpts.options.map((v, i) => formatFunc(v, i));

    radioState.label = radioOpts.label;
    radioState.options = displayVals;
    radioState.defaultValue = defaultVal;
    radioState.required = radioOpts.required;
    radioState.disabled = radioOpts.disabled;
  }

  session.state.set(widgetId, radioState);

  const radioProto = convertStateToRadioProto(radioState as RadioState);

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetId,
      type: {
        case: 'radio',
        value: radioProto,
      },
    },
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  // Return the selected value
  let value: RadioValue | null = null;
  if (radioState.value !== null && radioOpts.options.length > 0) {
    value = {
      value: radioOpts.options[radioState.value],
      index: radioState.value,
    };
  }

  return value;
}

/**
 * Convert radio state to proto
 * @param state Radio state
 * @returns Radio proto
 */
export function convertStateToRadioProto(state: RadioState): RadioProto {
  return fromJson(RadioSchema, {
    label: state.label,
    value: state.value,
    options: state.options,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
  });
}

/**
 * Convert radio proto to state
 * @param id Widget ID
 * @param data Radio proto
 * @returns Radio state
 */
export function convertRadioProtoToState(
  id: string,
  data: RadioProto | null,
): RadioState | null {
  if (!data) {
    return null;
  }

  const d = toJson(RadioSchema, data);

  return new RadioState(
    id,
    d.value as number,
    d.label,
    d.options,
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
