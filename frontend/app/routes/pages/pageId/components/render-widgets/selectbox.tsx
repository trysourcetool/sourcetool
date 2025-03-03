import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { cn } from '@/lib/utils';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';

const ExecuteSelectbox = ({
  widgetId,
  value,
  options,
}: {
  widgetId: string;
  value?: string;
  options: {
    label: string;
    value: string;
  }[];
}) => {
  const dispatch = useDispatch();
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleChange = (value: string) => {
    console.log({ handleChange: value });
    if (isWidgetWaiting) {
      return;
    }
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'selectbox',
        value: value === 'clear' ? undefined : Number(value),
      }),
    );
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'selectbox',
        value: value === 'clear' ? undefined : Number(value),
      }),
    );
  };

  console.log({ value, options });

  return (
    <Select
      disabled={isWidgetWaiting}
      value={value ?? ''}
      onValueChange={(value) => handleChange(value)}
    >
      <SelectTrigger>
        <SelectValue placeholder="Select an option" />
      </SelectTrigger>
      <SelectContent>
        {value !== undefined && (
          <SelectItem
            value={'clear'}
            className="font-normal text-muted-foreground"
          >
            Clear selection
          </SelectItem>
        )}
        {options.map((option) => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export const WidgetSelectbox: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );
  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );

  const options =
    widget && widget.widget?.selectbox
      ? (widget.widget.selectbox.options?.map((option, index) => ({
          label: option,
          value: index.toString(),
        })) ?? [])
      : [];

  return (
    widget &&
    widget.widget?.selectbox &&
    state.type === 'selectbox' && (
      <div className="space-y-2">
        {widget.widget.selectbox.label && (
          <Label
            className={cn('block', state.error && 'text-destructive')}
            htmlFor={widgetId}
          >
            {widget.widget.selectbox.label}
          </Label>
        )}
        <ExecuteSelectbox
          widgetId={widgetId}
          value={state.value !== undefined ? state.value.toString() : undefined}
          options={options}
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
