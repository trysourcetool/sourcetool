import { useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { type FC } from 'react';
import Markdown from 'react-markdown';

export const WidgetMarkdown: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );

  return (
    widget &&
    widget.widget?.markdown && (
      <div className="WidgetMarkdown">
        <Markdown>{widget.widget.markdown.body}</Markdown>
      </div>
    )
  );
};
