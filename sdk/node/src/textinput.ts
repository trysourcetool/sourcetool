import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  TextInputState,
  WidgetTypeTextInput,
} from './internal/session/state/textinput';
import { TextInputOptions } from './internal/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  TextInput as TextInputProto,
  TextInputSchema,
  WidgetSchema,
} from '@trysourcetool/proto/widget/v1/widget';
import { RenderWidgetSchema } from '@trysourcetool/proto/websocket/v1/message';
/**
 * TextInput component options
 */
export interface TextInputComponentOptions {
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
 * TextInput component class
 */
export class TextInput {
  /**
   * Set the placeholder text
   * @param placeholder Placeholder text
   * @returns TextInput options
   */
  static placeholder(placeholder: string): TextInputComponentOptions {
    return { placeholder };
  }

  /**
   * Set the default value
   * @param value Default value
   * @returns TextInput options
   */
  static defaultValue(value: string): TextInputComponentOptions {
    return { defaultValue: value };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns TextInput options
   */
  static required(required: boolean): TextInputComponentOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns TextInput options
   */
  static disabled(disabled: boolean): TextInputComponentOptions {
    return { disabled };
  }

  /**
   * Set the maximum length
   * @param length Maximum length
   * @returns TextInput options
   */
  static maxLength(length: number): TextInputComponentOptions {
    return { maxLength: length };
  }

  /**
   * Set the minimum length
   * @param length Minimum length
   * @returns TextInput options
   */
  static minLength(length: number): TextInputComponentOptions {
    return { minLength: length };
  }
}

/**
 * Add a text input to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options TextInput options
 * @returns The input value
 */
export function textInput(
  builder: UIBuilder,
  label: string,
  options: TextInputComponentOptions = {},
): string {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return '';
  }

  const textInputOpts: TextInputOptions = {
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
  const widgetID = builder.generatePageID(WidgetTypeTextInput, path);

  let textInputState = session.state.getTextInput(widgetID);
  if (!textInputState) {
    textInputState = new TextInputState(
      widgetID,
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
    session.state.set(widgetID, textInputState);
  }

  const textInputProto = convertStateToTextInputProto(
    textInputState as TextInputState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
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
function convertStateToTextInputProto(state: TextInputState): TextInputProto {
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
