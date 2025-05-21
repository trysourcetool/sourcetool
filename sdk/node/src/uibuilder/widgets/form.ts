import { v4 as uuidv4 } from 'uuid';
import {
  UIBuilder,
  Cursor,
  uiBuilderGeneratePageId,
  UIBuilderImpl,
} from '../index';
import { FormState, WidgetTypeForm } from '../../session/state/form';
import { FormInternalOptions } from '../../types/options';
import { create, fromJson } from '@bufbuild/protobuf';
import {
  Form as FormProto,
  FormSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
import { Runtime } from '../../runtime';
import { Session } from '../../session';
import { Page } from '../../page';
/**
 * Form options
 */
export interface FormOptions {
  /**
   * Whether the button is disabled
   * @default false
   */
  buttonDisabled?: boolean;

  /**
   * Whether to clear the form on submit
   * @default false
   */
  clearOnSubmit?: boolean;
}

/**
 * Add a form to the UI
 * @param builder The UI builder
 * @param buttonLabel The button label
 * @param options Form options
 * @returns A tuple containing the form builder and whether the form was submitted
 */
export function form(
  context: {
    runtime: Runtime;
    session: Session;
    page: Page;
    cursor: Cursor;
  },
  buttonLabel: string,
  options: FormOptions = {},
): [UIBuilder, boolean] {
  const { runtime, session, page, cursor } = context;

  if (!session || !page || !cursor) {
    return [new UIBuilderImpl(runtime, session, page), false];
  }

  const formOpts: FormInternalOptions = {
    buttonLabel,
    buttonDisabled: options.buttonDisabled || false,
    clearOnSubmit: options.clearOnSubmit || false,
  };

  const path = cursor.getPath();
  const widgetId = uiBuilderGeneratePageId(page.id, WidgetTypeForm, path);

  let formState = session.state.getForm(widgetId);
  if (!formState) {
    formState = new FormState(
      widgetId,
      false,
      formOpts.buttonLabel,
      formOpts.buttonDisabled,
      formOpts.clearOnSubmit,
    );
  } else {
    formState.buttonLabel = formOpts.buttonLabel;
    formState.buttonDisabled = formOpts.buttonDisabled;
    formState.clearOnSubmit = formOpts.clearOnSubmit;
  }
  session.state.set(widgetId, formState);

  const formProto = convertStateToFormProto(formState as FormState);

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetId,
      type: {
        case: 'form',
        value: formProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();

  // Create a child builder with a new cursor
  const childCursor = new Cursor();
  childCursor.parentPath = path;

  const childBuilder = new UIBuilderImpl(runtime, session, page, childCursor);

  return [childBuilder, formState.value];
}

/**
 * Convert form state to proto
 * @param state Form state
 * @returns Form proto
 */
export function convertStateToFormProto(state: FormState): FormProto {
  return fromJson(FormSchema, {
    value: state.value,
    buttonLabel: state.buttonLabel,
    buttonDisabled: state.buttonDisabled,
    clearOnSubmit: state.clearOnSubmit,
  });
}

/**
 * Convert form proto to state
 * @param id Widget ID
 * @param data Form proto
 * @returns Form state
 */
export function convertFormProtoToState(
  id: string,
  data: FormProto | null,
): FormState | null {
  if (!data) {
    return null;
  }

  return new FormState(
    id,
    data.value,
    data.buttonLabel,
    data.buttonDisabled,
    data.clearOnSubmit,
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
