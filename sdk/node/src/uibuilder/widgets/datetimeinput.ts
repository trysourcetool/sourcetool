import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import {
  DateTimeInputState,
  WidgetTypeDateTimeInput,
} from '../../session/state/datetimeinput';
import { DateTimeInputInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  DateTimeInput as DateTimeInputProto,
  DateTimeInputSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';

/**
 * DateTimeInput options
 */
export interface DateTimeInputOptions {
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
   * Date and time format
   * @default "YYYY/MM/DD HH:MM:SS"
   */
  format?: string;

  /**
   * Maximum date and time
   */
  maxValue?: Date;

  /**
   * Minimum date and time
   */
  minValue?: Date;

  /**
   * Timezone location
   * @default "local"
   */
  location?: string;
}

/**
 * Add a date and time input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options DateTimeInput options
 * @returns The input value
 */
export function dateTimeInput(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: DateTimeInputOptions = {},
): Date | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const dateTimeInputOpts: DateTimeInputInternalOptions = {
    label,
    placeholder: options.placeholder || '',
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    format: options.format || 'YYYY/MM/DD HH:MM:SS',
    maxValue: options.maxValue || null,
    minValue: options.minValue || null,
    location: options.location || 'local',
  };

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(
    page.id,
    WidgetTypeDateTimeInput,
    path,
  );

  let dateTimeInputState = session.state.getDateTimeInput(widgetId);
  if (!dateTimeInputState) {
    dateTimeInputState = new DateTimeInputState(
      widgetId,
      dateTimeInputOpts.defaultValue,
      dateTimeInputOpts.label,
      dateTimeInputOpts.placeholder,
      dateTimeInputOpts.defaultValue,
      dateTimeInputOpts.required,
      dateTimeInputOpts.disabled,
      dateTimeInputOpts.format,
      dateTimeInputOpts.maxValue,
      dateTimeInputOpts.minValue,
      dateTimeInputOpts.location,
    );
  } else {
    dateTimeInputState.label = dateTimeInputOpts.label;
    dateTimeInputState.placeholder = dateTimeInputOpts.placeholder;
    dateTimeInputState.defaultValue = dateTimeInputOpts.defaultValue;
    dateTimeInputState.required = dateTimeInputOpts.required;
    dateTimeInputState.disabled = dateTimeInputOpts.disabled;
    dateTimeInputState.format = dateTimeInputOpts.format;
    dateTimeInputState.maxValue = dateTimeInputOpts.maxValue;
    dateTimeInputState.minValue = dateTimeInputOpts.minValue;
    dateTimeInputState.location = dateTimeInputOpts.location;
  }
  session.state.set(widgetId, dateTimeInputState);

  const dateTimeInputProto = convertStateToDateTimeInputProto(
    dateTimeInputState as DateTimeInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'dateTimeInput',
        value: dateTimeInputProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return dateTimeInputState.value;
}

/**
 * Convert date and time input state to proto
 * @param state Date and time input state
 * @returns Date and time input proto
 */
export function convertStateToDateTimeInputProto(
  state: DateTimeInputState,
): DateTimeInputProto {
  const formatDateTime = (date: Date | null): string | null => {
    if (!date) {
      return null;
    }
    return date.toISOString(); // ISO format: YYYY-MM-DDTHH:mm:ss.sssZ
  };

  return fromJson(DateTimeInputSchema, {
    value: formatDateTime(state.value),
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: formatDateTime(state.defaultValue),
    required: state.required,
    disabled: state.disabled,
    format: state.format,
    maxValue: formatDateTime(state.maxValue),
    minValue: formatDateTime(state.minValue),
  });
}

/**
 * Convert date and time input proto to state
 * @param id Widget ID
 * @param data Date and time input proto
 * @returns Date and time input state
 */
export function convertDateTimeInputProtoToState(
  id: string,
  data: DateTimeInputProto | null,
  location: string = 'local',
): DateTimeInputState | null {
  if (!data) {
    return null;
  }

  const parseDateTime = (dateTimeStr: string | null): Date | null => {
    if (!dateTimeStr) {
      return null;
    }
    return new Date(dateTimeStr);
  };

  const d = toJson(DateTimeInputSchema, data);

  return new DateTimeInputState(
    id,
    parseDateTime(d.value || null),
    d.label,
    d.placeholder,
    parseDateTime(d.defaultValue || null),
    d.required,
    d.disabled,
    d.format,
    parseDateTime(d.maxValue || null),
    parseDateTime(d.minValue || null),
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
