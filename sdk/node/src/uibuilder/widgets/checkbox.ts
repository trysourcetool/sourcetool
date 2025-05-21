import { v4 as uuidv4 } from 'uuid';
import { Cursor, uiBuilderGeneratePageId } from '../';
import {
  CheckboxState,
  WidgetTypeCheckbox,
} from '../../session/state/checkbox';
import { CheckboxOptions } from '../../types/options';
import {
  Checkbox as CheckboxProto,
  CheckboxSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * Checkbox component options
 */
export interface CheckboxComponentOptions {
  /**
   * Default value of the checkbox
   * @default false
   */
  defaultValue?: boolean;

  /**
   * Whether the checkbox is required
   * @default false
   */
  required?: boolean;

  /**
   * Whether the checkbox is disabled
   * @default false
   */
  disabled?: boolean;
}

/**
 * Checkbox component class
 */
export class Checkbox {
  /**
   * Set the default value of the checkbox
   * @param value Default value
   * @returns Checkbox options
   */
  static defaultValue(value: boolean): CheckboxComponentOptions {
    return { defaultValue: value };
  }

  /**
   * Make the checkbox required
   * @param required Whether the checkbox is required
   * @returns Checkbox options
   */
  static required(required: boolean): CheckboxComponentOptions {
    return { required };
  }

  /**
   * Disable the checkbox
   * @param disabled Whether the checkbox is disabled
   * @returns Checkbox options
   */
  static disabled(disabled: boolean): CheckboxComponentOptions {
    return { disabled };
  }
}

/**
 * Add a checkbox to the UI
 * @param builder The UI builder
 * @param label The checkbox label
 * @param options Checkbox options
 * @returns Whether the checkbox is checked
 */
export function checkbox(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: CheckboxComponentOptions = {},
): boolean {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return false;
  }

  const checkboxOpts: CheckboxOptions = {
    label,
    defaultValue: options.defaultValue ?? false,
    required: options.required ?? false,
    disabled: options.disabled ?? false,
  };

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(page.id, WidgetTypeCheckbox, path);

  let checkboxState = session.state.getCheckbox(widgetId);
  if (!checkboxState) {
    checkboxState = new CheckboxState(
      widgetId,
      checkboxOpts.label,
      checkboxOpts.defaultValue,
      checkboxOpts.defaultValue,
      checkboxOpts.required,
      checkboxOpts.disabled,
    );
  } else {
    checkboxState.label = checkboxOpts.label;
    checkboxState.defaultValue = checkboxOpts.defaultValue;
    checkboxState.required = checkboxOpts.required;
    checkboxState.disabled = checkboxOpts.disabled;
  }
  session.state.set(widgetId, checkboxState);

  const checkboxProto = convertStateToCheckboxProto(
    checkboxState as CheckboxState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'checkbox',
        value: checkboxProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  return checkboxState.value;
}

/**
 * Convert checkbox state to proto
 * @param state Checkbox state
 * @returns Checkbox proto
 */
export function convertStateToCheckboxProto(
  state: CheckboxState,
): CheckboxProto {
  return fromJson(CheckboxSchema, {
    value: state.value,
    label: state.label,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
  });
}

/**
 * Convert checkbox proto to state
 * @param id Widget ID
 * @param data Checkbox proto
 * @returns Checkbox state
 */
export function convertCheckboxProtoToState(
  id: string,
  data: CheckboxProto | null,
): CheckboxState | null {
  if (!data) {
    return null;
  }

  const d = toJson(CheckboxSchema, data);

  return new CheckboxState(
    id,
    d.label,
    d.value,
    d.defaultValue,
    d.required,
    d.disabled,
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
