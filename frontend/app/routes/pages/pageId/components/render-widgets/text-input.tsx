import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';
import { useDebouncedCallback } from 'use-debounce';

const ExecuteTextInput = ({
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

  const handleChangeDebounce = useDebouncedCallback((value: string) => {
    console.log('debounce', value);
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'textInput',
        value,
      }),
    );
  }, 1000);

  const handleChange = (value: string) => {
    if (isWidgetWaiting) {
      return;
    }
    handleChangeDebounce(value);
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'textInput',
        value,
      }),
    );
  };

  return (
    <Input
      disabled={isWidgetWaiting}
      value={value}
      onChange={(e) => handleChange(e.target.value)}
      placeholder={placeholder}
      defaultValue={defaultValue}
    />
  );
};

export const WidgetTextInput: FC<{
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
    widget.widget?.textInput &&
    state.type === 'textInput' && (
      <div className="space-y-2">
        {widget.widget.textInput.label && (
          <Label
            className={cn('block', state.error && 'text-destructive')}
            htmlFor={widgetId}
          >
            {widget.widget.textInput.label}
          </Label>
        )}
        <ExecuteTextInput
          widgetId={widgetId}
          value={state.value}
          placeholder={widget.widget.textInput.placeholder}
          defaultValue={widget.widget.textInput.defaultValue}
        />
        {state.error && (
          <p className={cn('text-sm font-medium text-destructive')}>
            {state.error.message}
          </p>
        )}
      </div>
    )
  );
};
