import { useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';

export const WidgetColumnItem: FC<{
  widgetId: string;
  children?: React.ReactNode;
}> = ({ widgetId, children }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );

  return (
    widget &&
    widget.widget?.columnItem && (
      <div
        className="flex flex-1 flex-col gap-6"
        style={{ flexGrow: widget.widget.columnItem.weight }}
      >
        {children}
      </div>
    )
  );
};
