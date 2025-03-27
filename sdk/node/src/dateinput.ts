import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  DateInputState,
  WidgetTypeDateInput,
} from './internal/session/state/dateinput';
import { DateInputOptions } from './internal/options';

/**
 * DateInput component options
 */
export interface DateInputComponentOptions {
  /**
   * Placeholder text
   */
  placeholder?: string;

  /**
   * Default value
   */
  defaultValue?: Date;

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
   * Date format
   * @default "YYYY/MM/DD"
   */
  format?: string;

  /**
   * Maximum date
   */
  maxValue?: Date;

  /**
   * Minimum date
   */
  minValue?: Date;

  /**
   * Timezone location
   * @default "local"
   */
  location?: string;
}

/**
 * DateInput component class
 */
export class DateInput {
  /**
   * Set the placeholder text
   * @param placeholder Placeholder text
   * @returns DateInput options
   */
  static placeholder(placeholder: string): DateInputComponentOptions {
    return { placeholder };
  }

  /**
   * Set the default value
   * @param value Default value
   * @returns DateInput options
   */
  static defaultValue(value: Date): DateInputComponentOptions {
    return { defaultValue: value };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns DateInput options
   */
  static required(required: boolean): DateInputComponentOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns DateInput options
   */
  static disabled(disabled: boolean): DateInputComponentOptions {
    return { disabled };
  }

  /**
   * Set the date format
   * @param format Date format
   * @returns DateInput options
   */
  static format(format: string): DateInputComponentOptions {
    return { format };
  }

  /**
   * Set the maximum date
   * @param value Maximum date
   * @returns DateInput options
   */
  static maxValue(value: Date): DateInputComponentOptions {
    return { maxValue: value };
  }

  /**
   * Set the minimum date
   * @param value Minimum date
   * @returns DateInput options
   */
  static minValue(value: Date): DateInputComponentOptions {
    return { minValue: value };
  }

  /**
   * Set the timezone location
   * @param location Timezone location
   * @returns DateInput options
   */
  static location(location: string): DateInputComponentOptions {
    return { location };
  }
}

/**
 * Add a date input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options DateInput options
 * @returns The input value
 */
export function dateInput(
  builder: UIBuilder,
  label: string,
  options: DateInputComponentOptions = {},
): Date | null {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return null;
  }

  const dateInputOpts: DateInputOptions = {
    label,
    placeholder: options.placeholder || '',
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    format: options.format || 'YYYY/MM/DD',
    maxValue: options.maxValue || null,
    minValue: options.minValue || null,
    location: options.location || 'local',
  };

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeDateInput, path);

  let dateInputState = session.state.getDateInput(widgetID);
  if (!dateInputState) {
    dateInputState = new DateInputState(
      widgetID,
      dateInputOpts.defaultValue,
      dateInputOpts.label,
      dateInputOpts.placeholder,
      dateInputOpts.defaultValue,
      dateInputOpts.required,
      dateInputOpts.disabled,
      dateInputOpts.format,
      dateInputOpts.maxValue,
      dateInputOpts.minValue,
      dateInputOpts.location,
    );
  }

  dateInputState.label = dateInputOpts.label;
  dateInputState.placeholder = dateInputOpts.placeholder;
  dateInputState.defaultValue = dateInputOpts.defaultValue;
  dateInputState.required = dateInputOpts.required;
  dateInputState.disabled = dateInputOpts.disabled;
  dateInputState.format = dateInputOpts.format;
  dateInputState.maxValue = dateInputOpts.maxValue;
  dateInputState.minValue = dateInputOpts.minValue;
  dateInputState.location = dateInputOpts.location;
  session.state.set(widgetID, dateInputState);

  const dateInputProto = convertStateToDateInputProto(dateInputState);
  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'DateInput',
      dateInput: dateInputProto,
    },
  });

  cursor.next();

  return dateInputState.value;
}

/**
 * Convert date input state to proto
 * @param state Date input state
 * @returns Date input proto
 */
function convertStateToDateInputProto(state: DateInputState): any {
  const formatDate = (date: Date | null): string | null => {
    if (!date) {
      return null;
    }
    return date.toISOString().split('T')[0]; // YYYY-MM-DD format
  };

  return {
    value: formatDate(state.value),
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: formatDate(state.defaultValue),
    required: state.required,
    disabled: state.disabled,
    format: state.format,
    maxValue: formatDate(state.maxValue),
    minValue: formatDate(state.minValue),
  };
}

/**
 * Convert date input proto to state
 * @param id Widget ID
 * @param data Date input proto
 * @returns Date input state
 */
export function convertDateInputProtoToState(
  id: string,
  data: any,
  location: string = 'local',
): DateInputState | null {
  if (!data) {
    return null;
  }

  const parseDate = (dateStr: string | null): Date | null => {
    if (!dateStr) {
      return null;
    }
    return new Date(dateStr);
  };

  return new DateInputState(
    id,
    parseDate(data.value),
    data.label,
    data.placeholder,
    parseDate(data.defaultValue),
    data.required,
    data.disabled,
    data.format,
    parseDate(data.maxValue),
    parseDate(data.minValue),
    location,
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
