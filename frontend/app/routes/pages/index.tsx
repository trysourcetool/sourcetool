import { PageHeader } from '@/components/common/page-header';

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useDispatch, useSelector } from '@/store';
import { pagesStore } from '@/store/modules/pages';
import { File } from 'lucide-react';
import { useEffect, useRef } from 'react';
import { Link } from 'react-router';
import { useTranslation } from 'react-i18next';

export default function Pages() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');

  const pages = useSelector(pagesStore.selector.getPermissionPages);

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_pages') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        await dispatch(pagesStore.asyncActions.listPages());
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch]);

  return (
    <div>
      <PageHeader label="Pages" />
      <div className="p-6">
        <div className="rounded-md border p-4">
          <SidebarMenu>
            {/* TODO: Recursive processing and folder support */}
            {/* <Collapsible className="group/collapsible">
              <SidebarMenuItem>
                <CollapsibleTrigger asChild>
                  <SidebarMenuSubButton />
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    <SidebarMenuSubItem />
                  </SidebarMenuSub>
                </CollapsibleContent>
              </SidebarMenuItem>
            </Collapsible> */}
            {pages.map((page) => (
              <SidebarMenuItem key={page.id}>
                <SidebarMenuButton asChild>
                  <Link to={`/pages${page.route}`}>
                    <File className="size-4" />
                    {page.name}
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </div>
      </div>
    </div>
  );
}
