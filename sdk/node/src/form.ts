import { v4 as uuidv4 } from 'uuid';
import { UIBuilder, Cursor } from './uibuilder';
import { FormState, WidgetTypeForm } from './internal/session/state/form';
import { FormOptions } from './internal/options';
import { create, fromJson } from '@bufbuild/protobuf';
import {
  Form as FormProto,
  FormSchema,
  WidgetSchema,
} from '@trysourcetool/proto/widget/v1/widget';
import { RenderWidgetSchema } from '@trysourcetool/proto/websocket/v1/message';
/**
 * Form component options
 */
export interface FormComponentOptions {
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
 * Form component class
 */
export class Form {
  /**
   * Disable the button
   * @param disabled Whether the button is disabled
   * @returns Form options
   */
  static buttonDisabled(disabled: boolean): FormComponentOptions {
    return { buttonDisabled: disabled };
  }

  /**
   * Clear the form on submit
   * @param clear Whether to clear the form on submit
   * @returns Form options
   */
  static clearOnSubmit(clear: boolean): FormComponentOptions {
    return { clearOnSubmit: clear };
  }
}

/**
 * Add a form to the UI
 * @param builder The UI builder
 * @param buttonLabel The button label
 * @param options Form options
 * @returns A tuple containing the form builder and whether the form was submitted
 */
export function form(
  builder: UIBuilder,
  buttonLabel: string,
  options: FormComponentOptions = {},
): [UIBuilder, boolean] {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return [builder, false];
  }

  const formOpts: FormOptions = {
    buttonLabel,
    buttonDisabled: options.buttonDisabled || false,
    clearOnSubmit: options.clearOnSubmit || false,
  };

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeForm, path);

  let formState = session.state.getForm(widgetID);
  if (!formState) {
    formState = new FormState(
      widgetID,
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
  session.state.set(widgetID, formState);

  const formProto = convertStateToFormProto(formState as FormState);

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
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

  const childBuilder = new UIBuilder(runtime, session, page);
  childBuilder.cursor = childCursor;

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
