import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import { Label } from '@/components/ui/label';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import dayjs from 'dayjs';
import { CalendarIcon } from 'lucide-react';
import { type FC } from 'react';
import { useDebouncedCallback } from 'use-debounce';

const ExecuteDateInput = ({
  widgetId,
  value,
  placeholder,
  format,
}: {
  widgetId: string;
  value: string;
  placeholder?: string;
  format?: string;
}) => {
  const dispatch = useDispatch();
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleChangeDbounce = useDebouncedCallback(
    (value: Date | undefined) => {
      console.log('debounce', value);
      dispatch(
        widgetsStore.actions.setWidgetValue({
          widgetId,
          widgetType: 'dateInput',
          value: value ? dayjs(value).format('YYYY-MM-DD') : '',
        }),
      );
    },
    300,
  );

  const handleChange = (value: Date | undefined) => {
    if (isWidgetWaiting) {
      return;
    }
    handleChangeDbounce(value);
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'dateInput',
        value: value ? dayjs(value).format('YYYY-MM-DD') : '',
      }),
    );
  };

  return (
    <Popover>
      <PopoverTrigger asChild disabled={isWidgetWaiting}>
        <Button
          variant={'outline'}
          className={cn(
            'w-full justify-start text-left font-normal',
            !value && 'text-muted-foreground',
          )}
        >
          <CalendarIcon className="mr-2 size-4" />
          {value ? (
            dayjs(value).format(format || 'YYYY/MM/DD')
          ) : (
            <span>{placeholder || 'Pick a date'}</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <Calendar
          mode="single"
          selected={new Date(value)}
          onSelect={handleChange}
          initialFocus
          disabled={isWidgetWaiting}
        />
        {value && (
          <Button
            variant={'ghost'}
            className="w-full cursor-pointer font-normal text-muted-foreground"
            type="button"
            onClick={() => handleChange(undefined)}
          >
            Clear
          </Button>
        )}
      </PopoverContent>
    </Popover>
  );
};

export const WidgetDateInput: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );

  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );

  const value = (() => {
    if (state.type === 'dateInput') {
      if (state.value) {
        return state.value;
      }
    }
    return '';
  })();

  return (
    widget &&
    widget.widget?.dateInput && (
      <div className="space-y-2">
        {widget.widget.dateInput.label && (
          <Label className="block">{widget.widget.dateInput.label}</Label>
        )}
        <ExecuteDateInput
          widgetId={widgetId}
          value={value}
          placeholder={widget.widget.dateInput.placeholder}
          format={widget.widget.dateInput.format}
        />
      </div>
    )
  );
};
