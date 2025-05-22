import { v4 as uuidv4 } from 'uuid';
import { Cursor, generateWidgetId } from '../';
import {
  TimeInputState,
  WidgetTypeTimeInput,
} from '../../session/state/timeinput';
import { TimeInputInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  TimeInput as TimeInputProto,
  TimeInputSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * TimeInput options
 */
export interface TimeInputOptions {
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
   * Timezone location
   * @default "local"
   */
  location?: string;
}

/**
 * Add a time input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options TimeInput options
 * @returns The input value
 */
export function timeInput(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: TimeInputOptions = {},
): Date | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const timeInputOpts: TimeInputInternalOptions = {
    label,
    placeholder: options.placeholder || '',
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    location: options.location || 'local',
  };

  const path = cursor.getPath();
  const widgetId = generateWidgetId(page.id, WidgetTypeTimeInput, path);

  let timeInputState = session.state.getTimeInput(widgetId);
  if (!timeInputState) {
    timeInputState = new TimeInputState(
      widgetId,
      timeInputOpts.defaultValue,
      timeInputOpts.label,
      timeInputOpts.placeholder,
      timeInputOpts.defaultValue,
      timeInputOpts.required,
      timeInputOpts.disabled,
      timeInputOpts.location,
    );
    session.state.set(widgetId, timeInputState);
  } else {
    timeInputState.label = timeInputOpts.label;
    timeInputState.placeholder = timeInputOpts.placeholder;
    timeInputState.defaultValue = timeInputOpts.defaultValue;
    timeInputState.required = timeInputOpts.required;
    timeInputState.disabled = timeInputOpts.disabled;
    timeInputState.location = timeInputOpts.location;
    session.state.set(widgetId, timeInputState);
  }

  const timeInputProto = convertStateToTimeInputProto(
    timeInputState as TimeInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'timeInput',
        value: timeInputProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return timeInputState.value;
}

/**
 * Convert time input state to proto
 * @param state Time input state
 * @returns Time input proto
 */
export function convertStateToTimeInputProto(
  state: TimeInputState,
): TimeInputProto {
  const formatTime = (date: Date | null): string | null => {
    if (!date) {
      return null;
    }

    // Format as HH:MM:SS
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');

    return `${hours}:${minutes}:${seconds}`;
  };

  return fromJson(TimeInputSchema, {
    value: formatTime(state.value),
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: formatTime(state.defaultValue),
    required: state.required,
    disabled: state.disabled,
  });
}

/**
 * Convert time input proto to state
 * @param id Widget ID
 * @param data Time input proto
 * @returns Time input state
 */
export function convertTimeInputProtoToState(
  id: string,
  data: TimeInputProto | null,
  location: string = 'local',
): TimeInputState | null {
  if (!data) {
    return null;
  }

  const parseTime = (timeStr: string | null): Date | null => {
    if (!timeStr) {
      return null;
    }

    // Parse HH:MM:SS format
    const [hours, minutes, seconds] = timeStr.split(':').map(Number);
    const date = new Date();
    date.setHours(hours);
    date.setMinutes(minutes);
    date.setSeconds(seconds || 0);

    return date;
  };

  const d = toJson(TimeInputSchema, data);

  return new TimeInputState(
    id,
    parseTime(d.value || null),
    d.label,
    d.placeholder,
    parseTime(d.defaultValue || null),
    d.required,
    d.disabled,
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
