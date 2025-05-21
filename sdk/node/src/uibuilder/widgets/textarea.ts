import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import {
  TextAreaState,
  WidgetTypeTextArea,
} from '../../session/state/textarea';
import { TextAreaInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  TextArea as TextAreaProto,
  TextAreaSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';

/**
 * TextArea options
 */
export interface TextAreaOptions {
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

  /**
   * Maximum number of lines
   */
  maxLines?: number;

  /**
   * Minimum number of lines
   * @default 2
   */
  minLines?: number;

  /**
   * Whether to auto-resize the textarea
   * @default true
   */
  autoResize?: boolean;
}

/**
 * TextArea component class
 */
export class TextArea {
  /**
   * Set the placeholder text
   * @param placeholder Placeholder text
   * @returns TextArea options
   */
  static placeholder(placeholder: string): TextAreaOptions {
    return { placeholder };
  }

  /**
   * Set the default value
   * @param value Default value
   * @returns TextArea options
   */
  static defaultValue(value: string): TextAreaOptions {
    return { defaultValue: value };
  }

  /**
   * Make the input required
   * @param required Whether the input is required
   * @returns TextArea options
   */
  static required(required: boolean): TextAreaOptions {
    return { required };
  }

  /**
   * Disable the input
   * @param disabled Whether the input is disabled
   * @returns TextArea options
   */
  static disabled(disabled: boolean): TextAreaOptions {
    return { disabled };
  }

  /**
   * Set the maximum length of the input
   * @param length Maximum length
   * @returns TextArea options
   */
  static maxLength(length: number): TextAreaOptions {
    return { maxLength: length };
  }

  /**
   * Set the minimum length of the input
   * @param length Minimum length
   * @returns TextArea options
   */
  static minLength(length: number): TextAreaOptions {
    return { minLength: length };
  }

  /**
   * Set the maximum number of lines
   * @param lines Maximum number of lines
   * @returns TextArea options
   */
  static maxLines(lines: number): TextAreaOptions {
    return { maxLines: lines };
  }

  /**
   * Set the minimum number of lines
   * @param lines Minimum number of lines
   * @returns TextArea options
   */
  static minLines(lines: number): TextAreaOptions {
    return { minLines: lines };
  }

  /**
   * Set whether to auto-resize the textarea
   * @param autoResize Whether to auto-resize
   * @returns TextArea options
   */
  static autoResize(autoResize: boolean): TextAreaOptions {
    return { autoResize };
  }
}

/**
 * Add a textarea to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options TextArea options
 * @returns The input value
 */
export function textArea(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: TextAreaOptions = {},
): string {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return '';
  }

  // Set default minLines
  const defaultMinLines = 2;

  const textAreaOpts: TextAreaInternalOptions = {
    label,
    placeholder: options.placeholder || '',
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    maxLength: options.maxLength !== undefined ? options.maxLength : null,
    minLength: options.minLength !== undefined ? options.minLength : null,
    minLines:
      options.minLines !== undefined ? options.minLines : defaultMinLines,
    maxLines: options.maxLines !== undefined ? options.maxLines : null,
    autoResize: options.autoResize !== undefined ? options.autoResize : true,
  };

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(page.id, WidgetTypeTextArea, path);

  let textAreaState = session.state.getTextArea(widgetId);
  if (!textAreaState) {
    textAreaState = new TextAreaState(
      widgetId,
      textAreaOpts.defaultValue,
      textAreaOpts.label,
      textAreaOpts.placeholder,
      textAreaOpts.defaultValue,
      textAreaOpts.required,
      textAreaOpts.disabled,
      textAreaOpts.maxLength,
      textAreaOpts.minLength,
      textAreaOpts.maxLines,
      textAreaOpts.minLines,
      textAreaOpts.autoResize,
    );
  } else {
    textAreaState.label = textAreaOpts.label;
    textAreaState.placeholder = textAreaOpts.placeholder;
    textAreaState.defaultValue = textAreaOpts.defaultValue;
    textAreaState.required = textAreaOpts.required;
    textAreaState.disabled = textAreaOpts.disabled;
    textAreaState.maxLength = textAreaOpts.maxLength;
    textAreaState.minLength = textAreaOpts.minLength;
    textAreaState.maxLines = textAreaOpts.maxLines;
    textAreaState.minLines = textAreaOpts.minLines;
    textAreaState.autoResize = textAreaOpts.autoResize;
  }

  session.state.set(widgetId, textAreaState);

  const textAreaProto = convertStateToTextAreaProto(
    textAreaState as TextAreaState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'textArea',
        value: textAreaProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return textAreaState.value || '';
}

/**
 * Convert textarea state to proto
 * @param state TextArea state
 * @returns TextArea proto
 */
export function convertStateToTextAreaProto(
  state: TextAreaState,
): TextAreaProto {
  return fromJson(TextAreaSchema, {
    value: state.value,
    label: state.label,
    placeholder: state.placeholder,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
    maxLength: state.maxLength,
    minLength: state.minLength,
    maxLines: state.maxLines,
    minLines: state.minLines,
    autoResize: state.autoResize,
  });
}

/**
 * Convert textarea proto to state
 * @param id Widget ID
 * @param data TextArea proto
 * @returns TextArea state
 */
export function convertTextAreaProtoToState(
  id: string,
  data: TextAreaProto | null,
): TextAreaState | null {
  if (!data) {
    return null;
  }

  const d = toJson(TextAreaSchema, data);

  return new TextAreaState(
    id,
    d.value,
    d.label,
    d.placeholder,
    d.defaultValue,
    d.required,
    d.disabled,
    d.maxLength,
    d.minLength,
    d.maxLines,
    d.minLines,
    d.autoResize,
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
