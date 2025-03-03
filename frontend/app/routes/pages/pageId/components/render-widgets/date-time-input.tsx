import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
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

const ExecuteDateTimeInput = ({
  widgetId,
  value,
  placeholder,
  format,
}: {
  widgetId: string;
  value?: string;
  placeholder?: string;
  format?: string;
}) => {
  const dispatch = useDispatch();
  const hours = Array.from({ length: 24 }, (_, i) => i);
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const valueDate = value ? new Date(value) : undefined;

  const handleChangeDbounce = useDebouncedCallback(
    (value: Date | undefined) => {
      console.log('debounce', value);
      dispatch(
        widgetsStore.actions.setWidgetValue({
          widgetId,
          widgetType: 'dateInput',
          value: value ? dayjs(value).format('YYYY-MM-DD HH:MM:SS') : '',
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
        widgetType: 'dateTimeInput',
        value: value ? dayjs(value).format('YYYY-MM-DD HH:MM:SS') : '',
      }),
    );
  };

  const handleTimeChange = (
    type: 'hour' | 'minute' | 'second',
    value: string,
  ) => {
    if (value) {
      const newDate = new Date(value);
      if (type === 'hour') {
        newDate.setHours(parseInt(value));
      } else if (type === 'minute') {
        newDate.setMinutes(parseInt(value));
      } else if (type === 'second') {
        newDate.setSeconds(parseInt(value));
      }
      handleChange(newDate);
    }
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
          {valueDate ? (
            dayjs(valueDate).format(format || 'YYYY/MM/DD HH:MM:SS')
          ) : (
            <span>{placeholder || 'Pick a date and time'}</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <div className="sm:flex">
          <Calendar
            mode="single"
            selected={valueDate}
            onSelect={handleChange}
            initialFocus
            disabled={isWidgetWaiting}
          />
          <div className="flex flex-col divide-y sm:h-[300px] sm:flex-row sm:divide-x sm:divide-y-0">
            <ScrollArea className="w-64 sm:w-auto">
              <div className="flex p-2 sm:flex-col">
                {hours.reverse().map((hour) => (
                  <Button
                    key={hour}
                    size="icon"
                    variant={
                      valueDate && valueDate.getHours() === hour
                        ? 'default'
                        : 'ghost'
                    }
                    className="aspect-square shrink-0 sm:w-full"
                    onClick={() => handleTimeChange('hour', hour.toString())}
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
                    variant={
                      valueDate && valueDate.getMinutes() === minute
                        ? 'default'
                        : 'ghost'
                    }
                    className="aspect-square shrink-0 sm:w-full"
                    onClick={() =>
                      handleTimeChange('minute', minute.toString())
                    }
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
                    variant={
                      valueDate && valueDate.getSeconds() === second
                        ? 'default'
                        : 'ghost'
                    }
                    className="aspect-square shrink-0 sm:w-full"
                    onClick={() =>
                      handleTimeChange('second', second.toString())
                    }
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
            onClick={() => handleChange(undefined)}
          >
            Clear
          </Button>
        )}
      </PopoverContent>
    </Popover>
  );
};

export const WidgetDateTimeInput: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );

  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );

  const value = (() => {
    if (state.type === 'dateTimeInput') {
      if (state.value) {
        return state.value;
      }
    }
    return '';
  })();

  return (
    widget &&
    widget.widget?.dateTimeInput && (
      <div className="space-y-2">
        {widget.widget.dateTimeInput.label && (
          <Label className="block">{widget.widget.dateTimeInput.label}</Label>
        )}
        <ExecuteDateTimeInput
          widgetId={widgetId}
          value={value}
          placeholder={widget.widget.dateTimeInput.placeholder}
          format={widget.widget.dateTimeInput.format}
        />
      </div>
    )
  );
};
