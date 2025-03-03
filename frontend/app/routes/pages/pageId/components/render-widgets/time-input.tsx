import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area';
import { cn } from '@/lib/utils';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import dayjs from 'dayjs';
import { CalendarIcon } from 'lucide-react';
import { type FC } from 'react';
import { useDebouncedCallback } from 'use-debounce';

const ExecuteTimeInput = ({
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
  const hours = Array.from({ length: 24 }, (_, i) => i);
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleChangeDebounce = useDebouncedCallback((value: string) => {
    console.log('debounce', value);
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'timeInput',
        value: value,
      }),
    );
  }, 300);

  const handleChange = (type: 'hour' | 'minute' | 'second', value: string) => {
    if (isWidgetWaiting) {
      return;
    }
    const newTime = (value ?? '00:00:00').split(':');
    if (type === 'hour') {
      newTime[0] = value;
    } else if (type === 'minute') {
      newTime[1] = value;
    } else if (type === 'second') {
      newTime[2] = value;
    }
    handleChangeDebounce(newTime.join(':'));
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'timeInput',
        value: newTime.join(':'),
      }),
    );
  };

  const handleClear = () => {
    handleChangeDebounce('');
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'timeInput',
        value: '',
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
            dayjs(value).format(format || 'YYYY/MM/DD HH:MM:SS')
          ) : (
            <span>{placeholder || 'Pick a date and time'}</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <div className="sm:flex">
          <div className="flex flex-col divide-y sm:h-[300px] sm:flex-row sm:divide-x sm:divide-y-0">
            <ScrollArea className="w-64 sm:w-auto">
              <div className="flex p-2 sm:flex-col">
                {hours.reverse().map((hour) => (
                  <Button
                    key={hour}
                    size="icon"
                    variant={value ? 'default' : 'ghost'}
                    className="aspect-square shrink-0 sm:w-full"
                    onClick={() => handleChange('hour', hour.toString())}
                  >
                    {hour}
                  </Button>
                ))}
              </div>
              <ScrollBar orientation="horizontal" className="sm:hidden" />
            </ScrollArea>
            <ScrollArea className="w-64 sm:w-auto">
              <div className="flex p-2 sm:flex-col">
                {Array.from({ length: 12 }, (_, i) => i * 5).map((minute) => (
                  <Button
                    key={minute}
                    size="icon"
                    variant={value ? 'default' : 'ghost'}
                    className="aspect-square shrink-0 sm:w-full"
                    onClick={() => handleChange('minute', minute.toString())}
                  >
                    {minute.toString().padStart(2, '0')}
                  </Button>
                ))}
              </div>
              <ScrollBar orientation="horizontal" className="sm:hidden" />
            </ScrollArea>
            <ScrollArea className="w-64 sm:w-auto">
              <div className="flex p-2 sm:flex-col">
                {Array.from({ length: 12 }, (_, i) => i * 5).map((second) => (
                  <Button
                    key={second}
                    size="icon"
                    variant={value ? 'default' : 'ghost'}
                    className="aspect-square shrink-0 sm:w-full"
                    onClick={() => handleChange('second', second.toString())}
                  >
                    {second.toString().padStart(2, '0')}
                  </Button>
                ))}
              </div>
              <ScrollBar orientation="horizontal" className="sm:hidden" />
            </ScrollArea>
          </div>
        </div>
        {value && (
          <Button
            variant={'ghost'}
            className="w-full cursor-pointer font-normal text-muted-foreground"
            type="button"
            onClick={() => handleClear()}
          >
            Clear
          </Button>
        )}
      </PopoverContent>
    </Popover>
  );
};

export const WidgetTimeInput: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );

  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );

  const value = (() => {
    if (state.type === 'timeInput') {
      if (state.value) {
        return state.value;
      }
    }
    return '';
  })();

  return (
    widget &&
    widget.widget?.dateTimeInput &&
    state.type === 'timeInput' && (
      <div className="space-y-2">
        {widget.widget.dateTimeInput.label && (
          <Label
            className={cn('block', state.error && 'text-destructive')}
            htmlFor={widgetId}
          >
            {widget.widget.dateTimeInput.label}
          </Label>
        )}
        <ExecuteTimeInput
          widgetId={widgetId}
          value={value}
          placeholder={widget.widget.dateTimeInput.placeholder}
          format={widget.widget.dateTimeInput.format}
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
