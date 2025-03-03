import { Checkbox } from '@/components/ui/checkbox';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { useId, type FC } from 'react';

export const WidgetCheckbox: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const id = useId();
  const dispatch = useDispatch();
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );
  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleClick = (value: boolean) => {
    if (isWidgetWaiting) {
      return;
    }
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'checkbox',
        value,
      }),
    );
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'checkbox',
        value,
      }),
    );
  };

  return (
    widget &&
    widget.widget?.checkbox &&
    state.type === 'checkbox' && (
      <div className="items-top flex space-x-2">
        <Checkbox
          checked={state.value}
          onCheckedChange={() => handleClick(!state.value)}
          disabled={isWidgetWaiting}
          id={id}
        />
        {widget.widget?.checkbox?.label && (
          <div className="grid gap-1.5 leading-none">
            <label
              className="text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              htmlFor={id}
            >
              {widget.widget?.checkbox?.label}
            </label>
          </div>
        )}
      </div>
    )
  );
};
