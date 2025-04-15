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
import { Link, useNavigate, useParams } from 'react-router';
import { $path } from 'safe-routes';
import { useAuth } from '@/hooks/use-auth';
import { Loader2 } from 'lucide-react';
import { useDispatch, useSelector } from '@/store';
import { pagesStore } from '@/store/modules/pages';
import { usersStore } from '@/store/modules/users';
import { environmentsStore } from '@/store/modules/environments';
import { ENVIRONMENTS } from '@/environments';

export function AppPreviewLayout(props: PropsWithChildren) {
  const isInitialLoading = useRef(false);
  const { '*': path } = useParams();
  const {
    isSubDomainMatched,
    isAuthChecked,
    handleNoAuthRoute,
  } = useAuth();
  const dispatch = useDispatch();
  const user = useSelector(usersStore.selector.getUserMe);
  const pages = useSelector(pagesStore.selector.getPermissionPages);
  const { t } = useTranslation('common');
  const navigate = useNavigate();

  useEffect(() => {
    if (
      isAuthChecked === 'checked' &&
      ENVIRONMENTS.IS_CLOUD_EDITION &&
      !isSubDomainMatched
    ) {
      handleNoAuthRoute();
    }
  }, [isSubDomainMatched, isAuthChecked, handleNoAuthRoute]);

  const getLocalStorageSelectedEnvironmentId = (): string | null => {
    const environmentId = localStorage.getItem('selectedEnvironmentId');
    return environmentId || null;
  };

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        const resultAction = await dispatch(
          environmentsStore.asyncActions.listEnvironments(),
        );
        if (
          environmentsStore.asyncActions.listEnvironments.fulfilled.match(
            resultAction,
          )
        ) {
          const localStorageEnvironmentId =
            getLocalStorageSelectedEnvironmentId();
          console.log({ localStorageEnvironmentId });
          if (!localStorageEnvironmentId) {
            navigate($path('/'));
            return;
          }
          const hasEnvironmentId = resultAction.payload.environments.some(
            (e) => e.id === localStorageEnvironmentId,
          );
          console.log({ hasEnvironmentId }, resultAction.payload.environments);
          if (!hasEnvironmentId) {
            navigate($path('/'));
            return;
          }
          const resultActionPages = await dispatch(
            pagesStore.asyncActions.listPages({
              environmentId: localStorageEnvironmentId,
            }),
          );
          if (
            pagesStore.asyncActions.listPages.fulfilled.match(resultActionPages)
          ) {
            const hasPage = resultActionPages.payload.pages.some(
              (p) => p.route === `/${path}`,
            );
            console.log({ hasPage, path }, resultActionPages.payload.pages);
            if (!hasPage) {
              navigate($path('/'));
            }
          }
        }
        isInitialLoading.current = false;
      })();
    } else if (isAuthChecked === 'checked' && !ENVIRONMENTS.IS_CLOUD_EDITION && !user) {
      handleNoAuthRoute();
    }
  }, [dispatch]);

  return (isAuthChecked === 'checked' &&
    ENVIRONMENTS.IS_CLOUD_EDITION &&
    isSubDomainMatched) ||
    (isAuthChecked === 'checked' && !ENVIRONMENTS.IS_CLOUD_EDITION && user) ? (
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
                <SidebarMenuButton asChild isActive={`/${path}` === page.route}>
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
