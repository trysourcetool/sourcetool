import { v4 as uuidv4 } from 'uuid';
import { Cursor, generateWidgetId } from '../';
import {
  TextInputState,
  WidgetTypeTextInput,
} from '../../session/state/textinput';
import { TextInputInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  TextInput as TextInputProto,
  TextInputSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * TextInput options
 */
export interface TextInputOptions {
  /**
   * Placeholder text
   */
  placeholder?: string;

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
   * Maximum length of the input
   */
  maxLength?: number;

  /**
   * Minimum length of the input
   */
  minLength?: number;
}

/**
 * Add a text input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options TextInput options
 * @returns The input value
 */
export function textInput(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: TextInputOptions = {},
): string {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return '';
  }

  const textInputOpts: TextInputInternalOptions = {
    label,
    placeholder: options.placeholder || '',
    defaultValue:
      options.defaultValue !== undefined ? options.defaultValue : null,
    required: options.required || false,
    disabled: options.disabled || false,
    maxLength: options.maxLength !== undefined ? options.maxLength : null,
    minLength: options.minLength !== undefined ? options.minLength : null,
  };

  const path = cursor.getPath();
  const widgetId = generateWidgetId(page.id, WidgetTypeTextInput, path);

  let textInputState = session.state.getTextInput(widgetId);
  if (!textInputState) {
    textInputState = new TextInputState(
      widgetId,
      textInputOpts.defaultValue,
      textInputOpts.label,
      textInputOpts.placeholder,
      textInputOpts.defaultValue,
      textInputOpts.required,
      textInputOpts.disabled,
      textInputOpts.maxLength,
      textInputOpts.minLength,
    );
  } else {
    textInputState.label = textInputOpts.label;
    textInputState.placeholder = textInputOpts.placeholder;
    textInputState.defaultValue = textInputOpts.defaultValue;
    textInputState.required = textInputOpts.required;
    textInputState.disabled = textInputOpts.disabled;
    textInputState.maxLength = textInputOpts.maxLength;
    textInputState.minLength = textInputOpts.minLength;
  }
  session.state.set(widgetId, textInputState);

  const textInputProto = convertStateToTextInputProto(
    textInputState as TextInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'textInput',
        value: textInputProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return textInputState.value || '';
}

/**
 * Convert text input state to proto
 * @param state Text input state
 * @returns Text input proto
 */
export function convertStateToTextInputProto(
  state: TextInputState,
): TextInputProto {
  return fromJson(TextInputSchema, {
    value: state.value,
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
    maxLength: state.maxLength,
    minLength: state.minLength,
  });
}

/**
 * Convert text input proto to state
 * @param id Widget ID
 * @param data Text input proto
 * @returns Text input state
 */
export function convertTextInputProtoToState(
  id: string,
  data: TextInputProto | null,
): TextInputState | null {
  if (!data) {
    return null;
  }

  const d = toJson(TextInputSchema, data);

  return new TextInputState(
    id,
    d.value,
    d.label,
    d.placeholder,
    d.defaultValue,
    d.required,
    d.disabled,
    d.maxLength,
    d.minLength,
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
