import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  RadioState,
  RadioValue,
  WidgetTypeRadio,
} from './internal/session/state/radio';
import { RadioOptions } from './internal/options';

/**
 * Radio component options
 */
export interface RadioComponentOptions {
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
 * Radio component class
 */
export class Radio {
  /**
   * Set the radio options
   * @param options Radio options
   * @returns Radio options
   */
  static options(...options: string[]): RadioComponentOptions {
    return { options };
  }

  /**
   * Set the default value
   * @param value Default value
   * @returns Radio options
   */
  static defaultValue(value: string): RadioComponentOptions {
    return { defaultValue: value };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns Radio options
   */
  static required(required: boolean): RadioComponentOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns Radio options
   */
  static disabled(disabled: boolean): RadioComponentOptions {
    return { disabled };
  }

  /**
   * Set the format function for option labels
   * @param formatFunc Format function
   * @returns Radio options
   */
  static formatFunc(
    formatFunc: (value: string, index: number) => string,
  ): RadioComponentOptions {
    return { formatFunc };
  }
}

/**
 * Add a radio input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options Radio options
 * @returns The selected radio value
 */
export function radio(
  builder: UIBuilder,
  label: string,
  options: RadioComponentOptions = {},
): RadioValue | null {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return null;
  }

  const radioOpts: RadioOptions = {
    label,
    options: options.options || [],
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    formatFunc: options.formatFunc || ((v: string, i: number) => v),
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
  const widgetID = builder.generatePageID(WidgetTypeRadio, path);

  let radioState = session.state.getRadio(widgetID);
  if (!radioState) {
    radioState = new RadioState(
      widgetID,
      defaultVal,
      radioOpts.label,
      [],
      defaultVal,
      radioOpts.required,
      radioOpts.disabled,
    );
  } else {
    // Apply format function to options
    const formatFunc = radioOpts.formatFunc || ((v: string, i: number) => v);
    const displayVals = radioOpts.options.map((v, i) => formatFunc(v, i));

    radioState.label = radioOpts.label;
    radioState.options = displayVals;
    radioState.defaultValue = defaultVal;
    radioState.required = radioOpts.required;
    radioState.disabled = radioOpts.disabled;
    session.state.set(widgetID, radioState);
  }

  const radioProto = convertStateToRadioProto(radioState as RadioState);

  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'Radio',
      radio: radioProto,
    },
  });

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
function convertStateToRadioProto(state: RadioState): any {
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
 * Convert radio proto to state
 * @param id Widget ID
 * @param data Radio proto
 * @returns Radio state
 */
export function convertRadioProtoToState(
  id: string,
  data: any,
): RadioState | null {
  if (!data) {
    return null;
  }

  return new RadioState(
    id,
    data.value,
    data.label,
    data.options,
    data.defaultValue,
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
