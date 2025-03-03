import { useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';

export const WidgetColumns: FC<{
  widgetId: string;
  children?: React.ReactNode;
}> = ({ widgetId, children }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );

  return (
    widget &&
    widget.widget?.columns && <div className="flex gap-4">{children}</div>
  );
};
