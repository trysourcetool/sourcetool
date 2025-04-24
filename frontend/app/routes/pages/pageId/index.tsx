import { createFileRoute, useParams } from '@tanstack/react-router';
import { RenderWidgets } from './components/render-widgets';
import { useSelector } from '@/store';
import { pagesStore } from '@/store/modules/pages';
import { PageHeader } from '@/components/common/page-header';
import { ExceptionView } from './components/exception/exception-view';

export default function Preview() {
  const { _splat: path } = useParams({ from: '/_preview/pages/$' });
  console.log({ path });
  const page = useSelector((state) =>
    pagesStore.selector.getPageFromPath(state, `/${path}`),
  );
  const exception = useSelector((state) => state.pages.exception);
  return (
    <div>
      <PageHeader label={page?.name ?? ''} />
      <div className="space-y-4 px-4 py-6 md:space-y-6 md:px-6 md:py-6">
        {!exception && <RenderWidgets parentPath={[]} />}
        {exception && <ExceptionView />}
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_preview/pages/$')({
  component: Preview,
});
