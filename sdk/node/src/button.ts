import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import { ButtonState, WidgetTypeButton } from './internal/session/state/button';
import { ButtonOptions } from './internal/options';

/**
 * Button component options
 */
export interface ButtonComponentOptions {
  /**
   * Whether the button is disabled
   * @default false
   */
  disabled?: boolean;
}

/**
 * Button component class
 */
export class Button {
  /**
   * Create a disabled button
   * @param disabled Whether the button is disabled
   * @returns Button options
   */
  static disabled(disabled: boolean): ButtonComponentOptions {
    return { disabled };
  }
}

/**
 * Add a button to the UI
 * @param builder The UI builder
 * @param label The button label
 * @param options Button options
 * @returns Whether the button was clicked
 */
export function button(
  builder: UIBuilder,
  label: string,
  options: ButtonComponentOptions = {},
): boolean {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return false;
  }

  const buttonOpts: ButtonOptions = {
    label,
    disabled: options.disabled ?? false,
  };

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeButton, path);

  let buttonState = session.state.getButton(widgetID);
  if (!buttonState) {
    buttonState = new ButtonState(
      widgetID,
      false,
      buttonOpts.label,
      buttonOpts.disabled,
    );
  } else {
    buttonState.label = buttonOpts.label;
    buttonState.disabled = buttonOpts.disabled;
  }
  session.state.set(widgetID, buttonState);

  const buttonProto = convertStateToButtonProto(buttonState as ButtonState);
  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'Button',
      button: buttonProto,
    },
  });

  cursor.next();

  return buttonState.value;
}

/**
 * Convert button state to proto
 * @param state Button state
 * @returns Button proto
 */
function convertStateToButtonProto(state: ButtonState): any {
  return {
    value: state.value,
    label: state.label,
    disabled: state.disabled,
  };
}

/**
 * Convert button proto to state
 * @param id Widget ID
 * @param data Button proto
 * @returns Button state
 */
export function convertButtonProtoToState(
  id: string,
  data: any,
): ButtonState | null {
  if (!data) {
    return null;
  }

  return new ButtonState(id, data.value, data.label, data.disabled);
}

/**
 * Convert path to int32 array
 * @param path Path
 * @returns Int32 array
 */
export function convertPathToInt32Array(path: number[]): number[] {
  return path.map((v) => v);
}
