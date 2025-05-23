import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../';
import { createWidgetState, validateWidgetValue } from '@/lib/widgetState';
import {
  createEntityAdapter,
  createSlice,
  current,
  type EntityState,
  type PayloadAction,
} from '@reduxjs/toolkit';
import type { RenderWidgetJson } from '@/pb/ts/websocket/v1/message_pb';
import type {
  ButtonJson,
  CheckboxGroupJson,
  CheckboxJson,
  DateInputJson,
  DateTimeInputJson,
  FormJson,
  MultiSelectJson,
  NumberInputJson,
  RadioJson,
  SelectboxJson,
  TableJson,
  TextAreaJson,
  TextInputJson,
  TimeInputJson,
  WidgetJson,
} from '@/pb/ts/widget/v1/widget_pb';
import dayjs from 'dayjs';

// =============================================
// slice
// =============================================
export type WidgetType = Exclude<keyof WidgetJson, 'id'>;

export const inputWidgetTypes = [
  'textInput',
  'numberInput',
  'dateInput',
  'dateTimeInput',
  'timeInput',
  'selectbox',
  'textArea',
  'multiSelect',
  'checkbox',
  'radio',
  'checkboxGroup',
] as const;

export type Widget = RenderWidgetJson;

export const getChildFormItemWidgetIds = (
  widgets: RenderWidgetJson[],
  path: number[],
) =>
  widgets
    .filter((widget) => {
      if (
        !widget ||
        widget.widget?.form ||
        !inputWidgetTypes.some((type) => widget.widget && type in widget.widget)
      ) {
        return false;
      }
      const widgetPath = widget.path;
      if (!widgetPath || widgetPath.length <= path.length) {
        return false;
      }
      const slicedPath = widgetPath.slice(0, path.length);

      return slicedPath.every((p, index) => p === path[index]);
    })
    .map((widget) => widget.widget?.id) as string[];

export const checkParentForm = (
  widgets: RenderWidgetJson[],
  path: number[],
) => {
  const forms = widgets.filter((widget) => widget.widget?.form);

  return forms.some((form) =>
    form.path?.every((p, index) => p === path[index]),
  );
};

export type SetWidgetStatePayload = {
  widgetId: string;
} & (
  | {
      widgetType: Extract<WidgetType, 'textInput'>;
      value: TextInputJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'numberInput'>;
      value: NumberInputJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'dateInput'>;
      value: DateInputJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'dateTimeInput'>;
      value: DateTimeInputJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'timeInput'>;
      value: TimeInputJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'selectbox'>;
      value: SelectboxJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'textArea'>;
      value: TextAreaJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'multiSelect'>;
      value: MultiSelectJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'checkbox'>;
      value: CheckboxJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'radio'>;
      value: RadioJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'checkboxGroup'>;
      value: CheckboxGroupJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'button'>;
      value: ButtonJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'table'>;
      value: TableJson['value'];
    }
  | {
      widgetType: Extract<WidgetType, 'form'>;
      value: FormJson['value'];
    }
);

export type WidgetState = { id: string } & (
  | {
      type: Extract<WidgetType, 'textInput'>;
      value: TextInputJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'numberInput'>;
      value: NumberInputJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'dateInput'>;
      value: DateInputJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'dateTimeInput'>;
      value: DateTimeInputJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'timeInput'>;
      value: TimeInputJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'selectbox'>;
      value: SelectboxJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'textArea'>;
      value: TextAreaJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'multiSelect'>;
      value: MultiSelectJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'checkbox'>;
      value: CheckboxJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'radio'>;
      value: RadioJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'checkboxGroup'>;
      value: CheckboxGroupJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'form'>;
      value: FormJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'button'>;
      value: ButtonJson['value'];
      error: {
        message: string;
      } | null;
    }
  | {
      type: Extract<WidgetType, 'table'>;
      value: TableJson['value'];
      error: {
        message: string;
      } | null;
    }
);

// =============================================
// schema

const widgetsAdapter = createEntityAdapter<Widget, string>({
  selectId: (widget) => widget.widget?.id ?? '',
});

const widgetStateAdapter = createEntityAdapter<WidgetState, string>({
  selectId: (widget) => widget.id ?? '',
});

// =============================================
// State

export type State = {
  widgets: EntityState<Widget, string>;
  widgetStates: EntityState<WidgetState, string>;
  updateAt: number | null;
  isWidgetWaiting: boolean;
};

const initialState: State = {
  widgets: widgetsAdapter.getInitialState(),
  widgetStates: widgetStateAdapter.getInitialState(),
  updateAt: null,
  isWidgetWaiting: false,
};

// =============================================
// slice

const slice = createSlice({
  extraReducers: () => {},
  initialState,
  name: 'widgets',
  reducers: {
    setWidgetData: (state, action: PayloadAction<Widget>) => {
      if (!state.widgets.ids.includes(action.payload.widget?.id ?? '')) {
        widgetsAdapter.addOne(state.widgets, action.payload);
        if (action.payload.widget) {
          const widgetState = createWidgetState(action.payload.widget);
          if (widgetState) {
            widgetStateAdapter.addOne(state.widgetStates, widgetState);
          }
        }
      } else {
        const widgetState =
          state.widgetStates.entities[action.payload.widget?.id ?? ''];
        const payloadWidget = action.payload.widget;
        const widget =
          state.widgets.entities[action.payload.widget?.id ?? ''].widget;

        if (widget) {
          inputWidgetTypes.forEach((type) => {
            if (widget[type] && payloadWidget?.[type]) {
              if (
                payloadWidget[type].defaultValue !== widget[type].defaultValue
              ) {
                payloadWidget[type].value = payloadWidget[type].defaultValue;
                widgetState.value = payloadWidget[type].defaultValue;
              }
            }
          });
        }

        widgetsAdapter.updateOne(state.widgets, {
          id: action.payload.widget?.id ?? '',
          changes: { ...action.payload, widget: { ...payloadWidget } },
        });
      }
    },
    renderWidgetCompleted: (state) => {
      const widgets = state.widgets.ids.map((id) => state.widgets.entities[id]);
      const formIds = state.widgets.ids.filter((id) => {
        const widget = state.widgets.entities[id];
        return widget?.widget?.form;
      });
      let hasClearOnSubmit = false;
      formIds.forEach((id) => {
        const widget = state.widgets.entities[id];
        if (
          widget?.widget?.form?.value &&
          widget?.widget?.form?.clearOnSubmit
        ) {
          hasClearOnSubmit = true;
          const childFormItemWidgetIds = getChildFormItemWidgetIds(
            widgets,
            widget.path ?? [],
          );

          childFormItemWidgetIds.forEach((childId) => {
            const childWidget = state.widgets.entities[childId];
            const childWidgetState = state.widgetStates.entities[childId];
            if (childWidget?.widget) {
              const childWidgetType = Object.keys(childWidget.widget).filter(
                (key) => key !== 'id',
              )[0] as WidgetType;

              if ('value' in (childWidget?.widget?.[childWidgetType] ?? {})) {
                (childWidget.widget[childWidgetType] as any).value = (
                  childWidget.widget[childWidgetType] as any
                ).defaultValue;
                if (childWidgetState) {
                  childWidgetState.value = (
                    childWidget.widget[childWidgetType] as any
                  ).defaultValue;
                  childWidgetState.error = null;
                }
              }
            }
          });

          widget.widget.form.value = false;
        }
      });

      if (!hasClearOnSubmit) {
        state.isWidgetWaiting = false;
      } else {
        state.updateAt = dayjs().valueOf();
      }
    },
    clearWidgets: (state) => {
      widgetsAdapter.removeAll(state.widgets);
      widgetStateAdapter.removeAll(state.widgetStates);
      state.updateAt = null;
      state.isWidgetWaiting = false;
    },
    setWidgetState: (state, action: PayloadAction<SetWidgetStatePayload>) => {
      const widget = state.widgets.entities[action.payload.widgetId];
      if (widget?.widget) {
        const validateResult = validateWidgetValue(
          current(widget.widget),
          action.payload.widgetType,
          action.payload.value,
        );

        const widgetState =
          state.widgetStates.entities[action.payload.widgetId];
        if (widgetState) {
          widgetState.value = action.payload.value;
          if (validateResult?.error) {
            widgetState.error = {
              message: validateResult.error,
            };
          } else {
            widgetState.error = null;
          }
        }
      }
    },
    setWidgetValue: (state, action: PayloadAction<SetWidgetStatePayload>) => {
      const widget = state.widgets.entities[action.payload.widgetId];
      const widgets = state.widgets.ids.map((id) => state.widgets.entities[id]);
      if (widget.widget) {
        if (widget.widget.form) {
          const childFormItemWidgetIds = getChildFormItemWidgetIds(
            widgets,
            widget.path ?? [],
          );

          let hasError = false;
          childFormItemWidgetIds.forEach((id) => {
            const childWidget = state.widgets.entities[id];
            if (childWidget?.widget) {
              const widgetState = state.widgetStates.entities[id];
              const validateResult = validateWidgetValue(
                current(childWidget.widget),
                widgetState?.type,
                widgetState?.value,
              );
              if (validateResult?.error) {
                hasError = true;
                widgetState.error = {
                  message: validateResult.error,
                };
              } else {
                widgetState.error = null;
              }
            }
          });

          if (!hasError) {
            widget.widget.form.value = true;
            state.updateAt = dayjs().valueOf();
            state.isWidgetWaiting = true;
          }
        } else {
          const validateResult = validateWidgetValue(
            current(widget.widget),
            action.payload.widgetType,
            action.payload.value,
          );

          if (validateResult.success) {
            const hasParentForm = checkParentForm(widgets, widget.path ?? []);

            if (
              widget.widget &&
              Object.keys(widget.widget).includes(action.payload.widgetType)
            ) {
              const target = widget.widget[action.payload.widgetType];
              if (target) {
                target.value = action.payload.value;
              }
            }
            if (!hasParentForm) {
              state.updateAt = dayjs().valueOf();
              state.isWidgetWaiting = true;
            }
          }
        }
      }
    },
  },
});

// =============================================
// Selectors
// =============================================

const getWidgetIds = createSelector(
  (state: RootState) => state.widgets,
  (values) => values.widgets.ids,
);

const getWidgetEntities = createSelector(
  (state: RootState) => state.widgets,
  (values) => values.widgets.entities,
);

const getPathWidgetIdAndTypes = createSelector(
  (state: RootState, parentPath: number[], parentWidgetId?: string) => ({
    widgets: state.widgets,
    parentPath,
    parentWidgetId,
  }),
  (values) => {
    const { widgets, parentPath, parentWidgetId } = values;
    if (parentPath.length === 0) {
      return widgets.widgets.ids
        .filter((id) => widgets.widgets.entities[id]?.path?.length === 1)
        .map((id) => ({
          id,
          widgetType: Object.keys(
            widgets.widgets.entities[id]?.widget ?? {},
          ).find((key) => key !== 'id') as WidgetType,
        }));
    }

    return widgets.widgets.ids
      .filter((id) => {
        const path = widgets.widgets.entities[id]?.path;
        const slicedPath = path?.slice(0, parentPath.length);

        if (!parentWidgetId) {
          return false;
        }
        if (parentWidgetId && parentWidgetId === id) {
          return false;
        }

        return (
          slicedPath &&
          slicedPath.every((p, index) => p === parentPath[index]) &&
          parentPath.length + 1 === (path?.length ?? 0)
        );
      })
      .map((id) => ({
        id,
        widgetType: Object.keys(
          widgets.widgets.entities[id]?.widget ?? {},
        ).find((key) => key !== 'id') as WidgetType,
      }));
  },
);

const getWidgets = createSelector(
  (state: RootState) => state.widgets,
  (values) => values.widgets.ids.map((id) => values.widgets.entities[id]),
);

const getWidget = createSelector(
  (state: RootState, widgetId: string) =>
    state.widgets.widgets.entities[widgetId],
  (values) => values,
);

const getWidgetState = createSelector(
  (state: RootState, widgetId: string) =>
    state.widgets.widgetStates.entities[widgetId],
  (values) => values,
);

const getWidgetType = createSelector(
  (state: RootState, widgetId: string) =>
    state.widgets.widgets.entities[widgetId],
  (values) =>
    Object.keys(values?.widget ?? {}).find((key) => key !== 'id') as WidgetType,
);

// =============================================
// exports
// =============================================

export const widgetsStore = {
  actions: slice.actions,
  reducer: slice.reducer,
  selector: {
    getWidgetIds,
    getWidgetEntities,
    getPathWidgetIdAndTypes,
    getWidgets,
    getWidget,
    getWidgetState,
    getWidgetType,
  },
};

export type WidgetsState = State;
