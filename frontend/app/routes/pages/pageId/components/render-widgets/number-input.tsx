import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';
import { useDebouncedCallback } from 'use-debounce';

const ExecuteNumberInput = ({
  widgetId,
  value,
  placeholder,
}: {
  widgetId: string;
  value: number | undefined;
  placeholder?: string;
}) => {
  const dispatch = useDispatch();
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleChangeDbounce = useDebouncedCallback((value: string) => {
    console.log('debounce', value);
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'numberInput',
        value: value === '' ? undefined : Number(value),
      }),
    );
  }, 1000);

  const handleChange = (value: string) => {
    if (isWidgetWaiting) {
      return;
    }
    handleChangeDbounce(value);
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'numberInput',
        value: value === '' ? undefined : Number(value),
      }),
    );
  };

  return (
    <Input
      disabled={isWidgetWaiting}
      value={value ?? undefined}
      type="number"
      onChange={(e) => handleChange(e.target.value)}
      placeholder={placeholder}
    />
  );
};

export const WidgetNumberInput: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );
  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );

  return (
    widget &&
    widget.widget?.numberInput &&
    state.type === 'numberInput' && (
      <div className="space-y-2">
        {widget.widget.numberInput.label && (
          <Label className="block">{widget.widget.numberInput.label}</Label>
        )}
        <ExecuteNumberInput
          widgetId={widgetId}
          value={
            Number.isFinite(state.value) ? (state.value as number) : undefined
          }
          placeholder={widget.widget.numberInput.placeholder}
        />
      </div>
    )
  );
};
