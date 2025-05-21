import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import {
  DateInputState,
  WidgetTypeDateInput,
} from '../../session/state/dateinput';
import { DateInputInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  DateInput as DateInputProto,
  DateInputSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';

/**
 * DateInput options
 */
export interface DateInputOptions {
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
 * Add a date input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options DateInput options
 * @returns The input value
 */
export function dateInput(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: DateInputOptions = {},
): Date | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const dateInputOpts: DateInputInternalOptions = {
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
  const widgetId = uiBuilderGeneratePageId(page.id, WidgetTypeDateInput, path);

  let dateInputState = session.state.getDateInput(widgetId);
  if (!dateInputState) {
    dateInputState = new DateInputState(
      widgetId,
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
  session.state.set(widgetId, dateInputState);

  const dateInputProto = convertStateToDateInputProto(
    dateInputState as DateInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
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
