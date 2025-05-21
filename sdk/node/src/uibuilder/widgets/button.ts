import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import { ButtonState, WidgetTypeButton } from '../../session/state/button';
import { ButtonInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  ButtonSchema,
  Button as ButtonProto,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';

/**
 * Button options
 */
export interface ButtonOptions {
  /**
   * Whether the button is disabled
   * @default false
   */
  disabled?: boolean;
}

/**
 * Add a button to the UI
 * @param builder The UI builder
 * @param label The button label
 * @param options Button options
 * @returns Whether the button was clicked
 */
export function button(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: ButtonOptions = {},
): boolean {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return false;
  }

  const buttonOpts: ButtonInternalOptions = {
    label,
    disabled: options.disabled ?? false,
  };

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(page.id, WidgetTypeButton, path);

  let buttonState = session.state.getButton(widgetId);
  if (!buttonState) {
    buttonState = new ButtonState(
      widgetId,
      false,
      buttonOpts.label,
      buttonOpts.disabled,
    );
  } else {
    buttonState.label = buttonOpts.label;
    buttonState.disabled = buttonOpts.disabled;
  }
  session.state.set(widgetId, buttonState);

  const buttonProto = convertStateToButtonProto(buttonState as ButtonState);

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
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
export function convertStateToButtonProto(state: ButtonState): ButtonProto {
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
