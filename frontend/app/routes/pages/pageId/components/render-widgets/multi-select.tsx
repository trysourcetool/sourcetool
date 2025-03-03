import { Label } from '@/components/ui/label';
import { MultiSelect } from '@/components/ui/multi-select';

import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';

const ExecuteMultiSelect = ({
  widgetId,
  value,
  options,
  defaultValue,
}: {
  widgetId: string;
  value?: number[];
  options: {
    label: string;
    value: string;
  }[];
  defaultValue: string[];
}) => {
  const dispatch = useDispatch();
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleChange = (value: string[]) => {
    if (isWidgetWaiting) {
      return;
    }
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'multiSelect',
        value: value.map((v) => Number(v)),
      }),
    );
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'multiSelect',
        value: value.map((v) => Number(v)),
      }),
    );
  };

  return (
    <MultiSelect
      options={options}
      disabled={isWidgetWaiting}
      value={value ? value.map((v) => v.toString()) : undefined}
      onValueChange={(value) => handleChange(value)}
      defaultValue={defaultValue}
    />
  );
};

export const WidgetMultiSelect: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );

  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );

  const options =
    widget && widget.widget?.multiSelect
      ? (widget.widget.multiSelect.options?.map((option, index) => ({
          label: option,
          value: index.toString(),
        })) ?? [])
      : [];

  return (
    widget &&
    widget.widget?.multiSelect &&
    state.type === 'multiSelect' && (
      <div className="space-y-2">
        {widget.widget.multiSelect.label && (
          <Label className="block">{widget.widget.multiSelect.label}</Label>
        )}
        <ExecuteMultiSelect
          widgetId={widgetId}
          value={state.value}
          options={options}
          defaultValue={
            widget.widget.multiSelect.defaultValue
              ? widget.widget.multiSelect.defaultValue.map((v) => v.toString())
              : []
          }
        />
      </div>
    )
  );
};
