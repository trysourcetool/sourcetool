import { useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { WidgetTextInput } from './text-input';
import { WidgetTable } from './table';
import { WidgetButton } from './button';
import { WidgetColumns } from './columns';
import { WidgetColumnItem } from './columnitem';
import { WidgetMarkdown } from './markdown';
import { WidgetNumberInput } from './number-input';
import { WidgetSelectbox } from './selectbox';
import { WidgetMultiSelect } from './multi-select';
import { WidgetDateInput } from './date-input';
import { WidgetDateTimeInput } from './date-time-input';
import { WidgetTimeInput } from './time-input';
import { WidgetTextarea } from './textarea';
import { WidgetCheckbox } from './checkbox';
import { WidgetForm } from './form';
import { WidgetCheckboxGroup } from './checkbox-group';
import { WidgetRadio } from './radio';

export const RenderWidgets = ({
  parentPath,
  parentWidgetId,
}: {
  parentPath: number[];
  parentWidgetId?: string;
}) => {
  const widgetIdsAndTypes = useSelector((state) =>
    widgetsStore.selector.getPathWidgetIdAndTypes(
      state,
      parentPath,
      parentWidgetId,
    ),
  );

  return widgetIdsAndTypes.map(({ id, widgetType }, index) => {
    if (widgetType === 'textInput') {
      return <WidgetTextInput key={id} widgetId={id} />;
    }
    if (widgetType === 'numberInput') {
      return <WidgetNumberInput key={id} widgetId={id} />;
    }
    if (widgetType === 'dateInput') {
      return <WidgetDateInput key={id} widgetId={id} />;
    }
    if (widgetType === 'dateTimeInput') {
      return <WidgetDateTimeInput key={id} widgetId={id} />;
    }
    if (widgetType === 'timeInput') {
      return <WidgetTimeInput key={id} widgetId={id} />;
    }
    if (widgetType === 'radio') {
      return <WidgetRadio key={id} widgetId={id} />;
    }
    if (widgetType === 'textArea') {
      return <WidgetTextarea key={id} widgetId={id} />;
    }
    if (widgetType === 'selectbox') {
      return <WidgetSelectbox key={id} widgetId={id} />;
    }
    if (widgetType === 'multiSelect') {
      return <WidgetMultiSelect key={id} widgetId={id} />;
    }
    if (widgetType === 'checkbox') {
      return <WidgetCheckbox key={id} widgetId={id} />;
    }
    if (widgetType === 'checkboxGroup') {
      return <WidgetCheckboxGroup key={id} widgetId={id} />;
    }
    if (widgetType === 'table') {
      return <WidgetTable key={id} widgetId={id} />;
    }
    if (widgetType === 'button') {
      return <WidgetButton key={id} widgetId={id} />;
    }
    if (widgetType === 'markdown') {
      return <WidgetMarkdown key={id} widgetId={id} />;
    }
    if (widgetType === 'columns') {
      return (
        <WidgetColumns key={id} widgetId={id}>
          <RenderWidgets
            parentPath={[...parentPath, index]}
            parentWidgetId={id}
          />
        </WidgetColumns>
      );
    }
    if (widgetType === 'columnItem') {
      return (
        <WidgetColumnItem key={id} widgetId={id}>
          <RenderWidgets
            parentPath={[...parentPath, index]}
            parentWidgetId={id}
          />
        </WidgetColumnItem>
      );
    }
    if (widgetType === 'form') {
      return (
        <WidgetForm key={id} widgetId={id}>
          <RenderWidgets
            parentPath={[...parentPath, index]}
            parentWidgetId={id}
          />
        </WidgetForm>
      );
    }
    return null;
  });
};
