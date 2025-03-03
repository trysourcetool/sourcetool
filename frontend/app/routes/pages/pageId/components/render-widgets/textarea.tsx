import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';
import { useDebouncedCallback } from 'use-debounce';

const ExecuteTextarea = ({
  widgetId,
  value,
  placeholder,
  defaultValue,
}: {
  widgetId: string;
  value?: string;
  placeholder?: string;
  defaultValue?: string;
}) => {
  const dispatch = useDispatch();
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleChangeDbounce = useDebouncedCallback((value: string) => {
    console.log('debounce', value);
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'textArea',
        value,
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
        widgetType: 'textArea',
        value,
      }),
    );
  };

  return (
    <Textarea
      className="h-24 resize-none"
      disabled={isWidgetWaiting}
      value={value}
      onChange={(e) => handleChange(e.target.value)}
      placeholder={placeholder}
      defaultValue={defaultValue}
    />
  );
};

export const WidgetTextarea: FC<{
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
    widget.widget?.textArea &&
    state.type === 'textArea' && (
      <div className="space-y-2">
        {widget.widget.textArea.label && (
          <Label className="block">{widget.widget.textArea.label}</Label>
        )}
        <ExecuteTextarea
          widgetId={widgetId}
          value={state.value}
          placeholder={widget.widget.textArea.placeholder}
          defaultValue={widget.widget.textArea.defaultValue}
        />
      </div>
    )
  );
};
