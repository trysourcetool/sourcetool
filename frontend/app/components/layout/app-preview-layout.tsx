import { useEffect, useRef, type PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '../ui/sidebar';

import { ModeToggle } from '../common/mode-toggle';
import { Link, useParams } from 'react-router';
import { $path } from 'safe-routes';
import { useAuth } from '@/hooks/use-auth';
import { Loader2 } from 'lucide-react';
import { useDispatch, useSelector } from '@/store';
import { pagesStore } from '@/store/modules/pages';

export function AppPreviewLayout(props: PropsWithChildren) {
  const isInitialLoading = useRef(false);
  const { pageId } = useParams();
  const { subDomainMatched, handleNoAuthRoute } = useAuth();
  const dispatch = useDispatch();
  const pages = useSelector(pagesStore.selector.getPermissionPages);
  const { t } = useTranslation('common');

  useEffect(() => {
    if (subDomainMatched.status === 'checked' && !subDomainMatched.isMatched) {
      handleNoAuthRoute();
    }
  }, [subDomainMatched, handleNoAuthRoute]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        await dispatch(pagesStore.asyncActions.listPages());
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch]);

  return subDomainMatched.status === 'checked' && subDomainMatched.isMatched ? (
    <>
      <Sidebar collapsible="icon">
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                size="lg"
                className="w-full cursor-default data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              >
                <div className="flex flex-1 items-center gap-2 data-[state=open]:px-2 data-[state=open]:py-1">
                  <Link to={$path('/')} className="size-8">
                    <img
                      src="/images/logo-sidebar.png"
                      alt="Sourcetool"
                      className="size-full"
                    />
                  </Link>
                  <div className="flex flex-1 flex-col gap-0.5">
                    <p className="text-sm font-semibold text-sidebar-foreground">
                      {t('components_layout_app_name')}
                    </p>
                    <p className="text-xs font-normal text-sidebar-foreground">
                      {t('components_layout_app_version')}
                    </p>
                  </div>
                  <ModeToggle />
                </div>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            {pages.map((page) => (
              <SidebarMenu key={page.id}>
                <SidebarMenuButton asChild isActive={pageId === page.id}>
                  <Link to={`/pages${page.route}`}>
                    <span>{page.name}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenu>
            ))}
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter></SidebarFooter>
      </Sidebar>
      <SidebarInset>
        <div>{props.children}</div>
      </SidebarInset>
    </>
  ) : (
    <div className="flex h-screen w-full items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
