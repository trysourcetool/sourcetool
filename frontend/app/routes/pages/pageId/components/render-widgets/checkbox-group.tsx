import { Checkbox } from '@/components/ui/checkbox';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { useId, type FC } from 'react';

const WidgetCheckbox = ({
  checked,
  onChange,
  label,
  disabled,
}: {
  checked?: boolean;
  onChange: () => void;
  label: string;
  disabled: boolean;
}) => {
  const id = useId();
  return (
    <div className="items-top flex space-x-2">
      <Checkbox
        checked={checked}
        onCheckedChange={onChange}
        id={id}
        disabled={disabled}
      />
      <div className="grid gap-1.5 leading-none">
        <label
          className="cursor-pointer text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
          htmlFor={id}
        >
          {label}
        </label>
      </div>
    </div>
  );
};

export const WidgetCheckboxGroup: FC<{
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

  const handleClick = (index: number) => {
    if (isWidgetWaiting || !widget || !widget.widget?.checkboxGroup) {
      return;
    }
    const newValues = widget.widget.checkboxGroup.value?.includes(index)
      ? widget.widget.checkboxGroup.value?.filter((i) => i !== index)
      : [...(widget.widget.checkboxGroup.value || []), index];
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'checkboxGroup',
        value: newValues,
      }),
    );
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'checkboxGroup',
        value: newValues,
      }),
    );
  };

  return (
    widget &&
    widget.widget?.checkboxGroup &&
    state.type === 'checkboxGroup' && (
      <div className="flex flex-wrap gap-4">
        {widget.widget.checkboxGroup.options?.map((option, index) => (
          <WidgetCheckbox
            disabled={isWidgetWaiting}
            key={index}
            checked={state?.value?.includes(index)}
            onChange={() => handleClick(index)}
            label={option}
          />
        ))}
      </div>
    )
  );
};
