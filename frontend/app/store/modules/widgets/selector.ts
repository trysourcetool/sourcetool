import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';
import { type WidgetType } from './slice';

export const getWidgetIds = createSelector(
  (state: RootState) => state.widgets,
  (values) => values.widgets.ids,
);

export const getWidgetEntities = createSelector(
  (state: RootState) => state.widgets,
  (values) => values.widgets.entities,
);

export const getPathWidgetIdAndTypes = createSelector(
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

export const getWidgets = createSelector(
  (state: RootState) => state.widgets,
  (values) => values.widgets.ids.map((id) => values.widgets.entities[id]),
);

export const getWidget = createSelector(
  (state: RootState, widgetId: string) =>
    state.widgets.widgets.entities[widgetId],
  (values) => values,
);

export const getWidgetState = createSelector(
  (state: RootState, widgetId: string) =>
    state.widgets.widgetStates.entities[widgetId],
  (values) => values,
);

export const getWidgetType = createSelector(
  (state: RootState, widgetId: string) =>
    state.widgets.widgets.entities[widgetId],
  (values) =>
    Object.keys(values?.widget ?? {}).find((key) => key !== 'id') as WidgetType,
);
