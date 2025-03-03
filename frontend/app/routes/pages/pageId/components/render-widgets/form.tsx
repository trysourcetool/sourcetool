import { Button } from '@/components/ui/button';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';

export const WidgetForm: FC<{
  widgetId: string;
  children?: React.ReactNode;
}> = ({ widgetId, children }) => {
  const dispatch = useDispatch();
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'form',
        value: true,
      }),
    );
  };

  return (
    widget &&
    widget.widget?.form && (
      <form className="flex flex-col gap-6" onSubmit={handleSubmit}>
        {children}
        <Button
          type="submit"
          disabled={widget.widget.form.buttonDisabled || isWidgetWaiting}
        >
          {widget.widget.form.buttonLabel}
        </Button>
      </form>
    )
  );
};
