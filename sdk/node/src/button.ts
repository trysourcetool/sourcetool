import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import { ButtonState, WidgetTypeButton } from './internal/session/state/button';
import { ButtonOptions } from './internal/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  ButtonSchema,
  Button as ButtonProto,
  WidgetSchema,
} from '@trysourcetool/proto/widget/v1/widget';
import { RenderWidgetSchema } from '@trysourcetool/proto/websocket/v1/message';

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

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
      type: {
        case: 'button',
        value: buttonProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return buttonState.value;
}

/**
 * Convert button state to proto
 * @param state Button state
 * @returns Button proto
 */
function convertStateToButtonProto(state: ButtonState): ButtonProto {
  return fromJson(ButtonSchema, {
    value: state.value,
    label: state.label,
    disabled: state.disabled,
  });
}

/**
 * Convert button proto to state
 * @param id Widget ID
 * @param data Button proto
 * @returns Button state
 */
export function convertButtonProtoToState(
  id: string,
  data: ButtonProto | null,
): ButtonState | null {
  if (!data) {
    return null;
  }

  const d = toJson(ButtonSchema, data);

  return new ButtonState(id, d.value, d.label, d.disabled);
}

/**
 * Convert path to int32 array
 * @param path Path
 * @returns Int32 array
 */
export function convertPathToInt32Array(path: number[]): number[] {
  return path.map((v) => v);
}
