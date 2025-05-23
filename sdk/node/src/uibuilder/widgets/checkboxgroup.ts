import { v4 as uuidv4 } from 'uuid';
import { Cursor, generateWidgetId } from '../';
import {
  CheckboxGroupState,
  CheckboxGroupValue,
  WidgetTypeCheckboxGroup,
} from '../../session/state/checkboxgroup';
import { CheckboxGroupInternalOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  CheckboxGroup as CheckboxGroupProto,
  CheckboxGroupSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';

/**
 * CheckboxGroup options
 */
export interface CheckboxGroupOptions {
  /**
   * CheckboxGroup options
   */
  options?: string[];

  /**
   * Default values
   */
  defaultValue?: string[];

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
   * Format function for option labels
   */
  formatFunc?: (value: string, index: number) => string;
}

/**
 * Add a checkboxgroup to the UI
 * @param builder The UI builder
 * @param label The input label
 * @param options CheckboxGroup options
 * @returns The selected values
 */
export function checkboxGroup(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  label: string,
  options: CheckboxGroupOptions = {},
): CheckboxGroupValue | null {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return null;
  }

  const checkboxGroupOpts: CheckboxGroupInternalOptions = {
    label,
    options: options.options || [],
    defaultValue: options.defaultValue || null,
    required: options.required || false,
    disabled: options.disabled || false,
    formatFunc: options.formatFunc || ((v: string) => v),
  };

  // Find default value indexes
  const defaultVal: number[] = [];
  if (checkboxGroupOpts.defaultValue && checkboxGroupOpts.options.length > 0) {
    for (const defaultOption of checkboxGroupOpts.defaultValue) {
      for (let i = 0; i < checkboxGroupOpts.options.length; i++) {
        if (checkboxGroupOpts.options[i] === defaultOption) {
          defaultVal.push(i);
          break;
        }
      }
    }
  }

  const path = cursor.getPath();
  const widgetId = generateWidgetId(page.id, WidgetTypeCheckboxGroup, path);

  let checkboxGroupState = session.state.getCheckboxGroup(widgetId);
  const formatFunc = checkboxGroupOpts.formatFunc || ((v: string) => v);
  if (!checkboxGroupState) {
    checkboxGroupState = new CheckboxGroupState(
      widgetId,
      defaultVal,
      checkboxGroupOpts.label,
      checkboxGroupOpts.options.map(formatFunc),
      defaultVal,
      checkboxGroupOpts.required,
      checkboxGroupOpts.disabled,
    );
  } else {
    const displayVals = checkboxGroupOpts.options.map((v, i) =>
      formatFunc(v, i),
    );

    checkboxGroupState.label = checkboxGroupOpts.label;
    checkboxGroupState.options = displayVals;
    checkboxGroupState.defaultValue = defaultVal;
    checkboxGroupState.required = checkboxGroupOpts.required;
    checkboxGroupState.disabled = checkboxGroupOpts.disabled;
  }

  session.state.set(widgetId, checkboxGroupState);

  const checkboxGroupProto = convertStateToCheckboxGroupProto(
    checkboxGroupState as CheckboxGroupState,
  );

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'checkboxGroup',
        value: checkboxGroupProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  // Return the selected values
  let value: CheckboxGroupValue | null = null;
  if (
    checkboxGroupState.value.length > 0 &&
    checkboxGroupOpts.options.length > 0
  ) {
    value = {
      values: checkboxGroupState.value.map(
        (idx: number) => checkboxGroupOpts.options[idx],
      ),
      indexes: checkboxGroupState.value.map((idx: number) => idx),
    };
  }

  return value;
}

/**
 * Convert checkboxgroup state to proto
 * @param state CheckboxGroup state
 * @returns CheckboxGroup proto
 */
export function convertStateToCheckboxGroupProto(
  state: CheckboxGroupState,
): CheckboxGroupProto {
  return fromJson(CheckboxGroupSchema, {
    label: state.label,
    value: state.value,
    options: state.options,
    defaultValue: state.defaultValue,
    required: state.required,
    disabled: state.disabled,
  });
}

/**
 * Convert checkboxgroup proto to state
 * @param id Widget ID
 * @param data CheckboxGroup proto
 * @returns CheckboxGroup state
 */
export function convertCheckboxGroupProtoToState(
  id: string,
  data: CheckboxGroupProto | null,
): CheckboxGroupState | null {
  if (!data) {
    return null;
  }

  const d = toJson(CheckboxGroupSchema, data);

  return new CheckboxGroupState(
    id,
    d.value || [],
    d.label,
    d.options || [],
    d.defaultValue || [],
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
