import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { cn } from '@/lib/utils';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { useId, type FC } from 'react';

const RadioComponent = ({ label, value }: { label: string; value: string }) => {
  const id = useId();
  return (
    <div className="flex items-center space-x-2">
      <RadioGroupItem value={value} id={id} />
      <Label htmlFor={id}>{label}</Label>
    </div>
  );
};

export const WidgetRadio: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const dispatch = useDispatch();
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );
  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );

  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleClick = (value: number) => {
    console.log('value', value);
    if (isWidgetWaiting) {
      return;
    }
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'radio',
        value,
      }),
    );
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'radio',
        value: value,
      }),
    );
  };

  return (
    widget &&
    widget.widget?.radio &&
    state.type === 'radio' && (
      <>
        <RadioGroup
          onValueChange={(value) => handleClick(Number(value))}
          defaultValue={state.value?.toString()}
          disabled={isWidgetWaiting}
        >
          {widget.widget.radio.options?.map((option, index) => (
            <RadioComponent
              key={index}
              value={index.toString()}
              label={option}
            />
          ))}
        </RadioGroup>
        {state.error && (
          <p className={cn('text-sm font-medium text-destructive')}>
            {state.error.message}
          </p>
        )}
      </>
    )
  );
};
