import { Button } from '@/components/ui/button';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { Circle } from 'lucide-react';
import { type FC } from 'react';

export const WidgetButton: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const dispatch = useDispatch();
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleClick = () => {
    if (isWidgetWaiting) {
      return;
    }
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'button',
        value: true,
      }),
    );
  };

  return (
    widget &&
    widget.widget?.button && (
      <div>
        <Button
          disabled={widget.widget?.button?.disabled || isWidgetWaiting}
          onClick={handleClick}
        >
          <Circle className="size-4" />
          {widget.widget?.button?.label}
        </Button>
      </div>
    )
  );
};
