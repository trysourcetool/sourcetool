import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  DateInputState,
  WidgetTypeDateInput,
} from './internal/session/state/dateinput';
import { DateInputOptions } from './internal/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  DateInput as DateInputProto,
  DateInputSchema,
  WidgetSchema,
} from '@trysourcetool/proto/widget/v1/widget';
import { RenderWidgetSchema } from '@trysourcetool/proto/websocket/v1/message';

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
  } else {
    dateInputState.label = dateInputOpts.label;
    dateInputState.placeholder = dateInputOpts.placeholder;
    dateInputState.defaultValue = dateInputOpts.defaultValue;
    dateInputState.required = dateInputOpts.required;
    dateInputState.disabled = dateInputOpts.disabled;
    dateInputState.format = dateInputOpts.format;
    dateInputState.maxValue = dateInputOpts.maxValue;
    dateInputState.minValue = dateInputOpts.minValue;
    dateInputState.location = dateInputOpts.location;
  }
  session.state.set(widgetID, dateInputState);

  const dateInputProto = convertStateToDateInputProto(
    dateInputState as DateInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
      type: {
        case: 'dateInput',
        value: dateInputProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return dateInputState.value;
}

/**
 * Convert date input state to proto
 * @param state Date input state
 * @returns Date input proto
 */
export function convertStateToDateInputProto(
  state: DateInputState,
): DateInputProto {
  const formatDate = (date: Date | null): string | null => {
    if (!date) {
      return null;
    }
    return date.toISOString().split('T')[0]; // YYYY-MM-DD format
  };

  return fromJson(DateInputSchema, {
    value: formatDate(state.value),
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: formatDate(state.defaultValue),
    required: state.required,
    disabled: state.disabled,
    format: state.format,
    maxValue: formatDate(state.maxValue),
    minValue: formatDate(state.minValue),
  });
}

/**
 * Convert date input proto to state
 * @param id Widget ID
 * @param data Date input proto
 * @returns Date input state
 */
export function convertDateInputProtoToState(
  id: string,
  data: DateInputProto | null,
  location: string = 'local',
): DateInputState | null {
  if (!data) {
    return null;
  }

  const d = toJson(DateInputSchema, data);

  const parseDate = (dateStr: string | null): Date | null => {
    if (!dateStr) {
      return null;
    }
    return new Date(dateStr);
  };

  return new DateInputState(
    id,
    parseDate(d.value || null),
    d.label,
    d.placeholder,
    parseDate(d.defaultValue || null),
    d.required,
    d.disabled,
    d.format,
    parseDate(d.maxValue || null),
    parseDate(d.minValue || null),
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
