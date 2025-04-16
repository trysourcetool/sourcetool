import { Fragment, useEffect, type PropsWithChildren } from 'react';
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
  SidebarTrigger,
  useSidebar,
} from '../ui/sidebar';
import { Separator } from '../ui/separator';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '../ui/breadcrumb';
import { ModeToggle } from '../common/mode-toggle';
import { Link, useLocation } from 'react-router';
import {
  CheckCheck,
  ChevronsUpDown,
  FileText,
  KeyRound,
  Loader2,
  LogOut,
  Settings2,
  Split,
  Users,
} from 'lucide-react';
import { $path } from 'safe-routes';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { DropdownMenu } from '@radix-ui/react-dropdown-menu';
import {
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu';
import { Avatar, AvatarFallback } from '../ui/avatar';
import { useAuth } from '@/hooks/use-auth';
import { usersStore } from '@/store/modules/users';
import { useDispatch, useSelector } from '@/store';
import { ENVIRONMENTS } from '@/environments';
import { authStore } from '@/store/modules/auth';

export function AppExternalLayout(props: PropsWithChildren) {
  const dispatch = useDispatch();
  const { pathname } = useLocation();
  const { isMobile, setOpenMobile } = useSidebar();
  const { breadcrumbsState } = useBreadcrumbs();
  const {
    isSubDomainMatched,
    isAuthChecked,
    handleNoAuthRoute,
  } = useAuth();
  const user = useSelector(usersStore.selector.getUserMe);
  const { t } = useTranslation('common');

  const handleSignout = async () => {
    await dispatch(authStore.asyncActions.logout());
  };

  const handleSidebarClose = () => {
    console.log('isMobile', isMobile);
    if (isMobile) {
      setOpenMobile(false);
    }
  };

  useEffect(() => {
    if (
      isAuthChecked === 'checked' &&
      ENVIRONMENTS.IS_CLOUD_EDITION &&
      !isSubDomainMatched
    ) {
      handleNoAuthRoute();
    } else if (isAuthChecked === 'checked' && !ENVIRONMENTS.IS_CLOUD_EDITION && !user) {
      handleNoAuthRoute();
    }
  }, [isSubDomainMatched, isAuthChecked, handleNoAuthRoute]);

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
                  <Link
                    to={$path('/')}
                    className="size-8"
                    onClick={handleSidebarClose}
                  >
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
            <SidebarMenu>
              <SidebarMenuButton asChild isActive={pathname === '/'}>
                <Link to={$path('/')} onClick={handleSidebarClose}>
                  <FileText />
                  <span>{t('components_layout_sidebar_pages')}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenu>
            <SidebarMenu>
              <SidebarMenuButton asChild isActive={pathname.includes('/users')}>
                <Link to={$path('/users')} onClick={handleSidebarClose}>
                  <Users />
                  <span>{t('components_layout_sidebar_users')}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenu>
            <SidebarMenu>
              <SidebarMenuButton
                asChild
                isActive={pathname.includes('/environments')}
              >
                <Link to={$path('/environments')} onClick={handleSidebarClose}>
                  <Split />
                  <span>{t('components_layout_sidebar_environments')}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenu>
            <SidebarMenu>
              <SidebarMenuButton
                asChild
                isActive={pathname.includes('/apiKeys')}
              >
                <Link to={$path('/apiKeys')} onClick={handleSidebarClose}>
                  <KeyRound />
                  <span>{t('components_layout_sidebar_api_keys')}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenu>
            <SidebarMenu>
              <SidebarMenuButton
                asChild
                isActive={pathname.includes('/groups')}
              >
                <Link to={$path('/groups')} onClick={handleSidebarClose}>
                  <CheckCheck />
                  <span>{t('components_layout_sidebar_permissions')}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenu>
            <SidebarMenu>
              <SidebarMenuButton
                asChild
                isActive={pathname.includes('/settings')}
              >
                <Link to={$path('/settings')} onClick={handleSidebarClose}>
                  <Settings2 />
                  <span>{t('components_layout_sidebar_settings')}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenu>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <SidebarMenuButton
                    size="lg"
                    className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                  >
                    <Avatar className="size-8 rounded-lg">
                      <AvatarFallback className="rounded-lg">
                        {user?.firstName[0]}
                        {user?.lastName[0]}
                      </AvatarFallback>
                    </Avatar>
                    <div className="grid flex-1 text-left text-sm leading-tight">
                      <span className="truncate font-semibold">
                        {user?.firstName} {user?.lastName}
                      </span>
                      <span className="truncate text-xs">{user?.email}</span>
                    </div>
                    <ChevronsUpDown className="ml-auto size-4" />
                  </SidebarMenuButton>
                </DropdownMenuTrigger>
                <DropdownMenuContent
                  className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                  side={isMobile ? 'bottom' : 'right'}
                  align="end"
                  sideOffset={4}
                >
                  <DropdownMenuLabel className="p-0 font-normal">
                    <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                      <Avatar className="size-8 rounded-lg">
                        <AvatarFallback className="rounded-lg">
                          {user?.firstName[0]}
                          {user?.lastName[0]}
                        </AvatarFallback>
                      </Avatar>
                      <div className="grid flex-1 text-left text-sm leading-tight">
                        <span className="truncate font-semibold">
                          {user?.firstName} {user?.lastName}
                        </span>
                        <span className="truncate text-xs">{user?.email}</span>
                      </div>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={handleSignout}
                    className="cursor-pointer"
                  >
                    <LogOut />
                    {t('components_layout_sidebar_logout')}
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sidebar>
      <SidebarInset>
        <header className="sticky top-0 z-10 flex h-16 shrink-0 items-center gap-2 border-b bg-background transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <Breadcrumb>
              <BreadcrumbList>
                {breadcrumbsState?.map((breadcrumb, index) => (
                  <Fragment key={breadcrumb.label}>
                    {!!index && <BreadcrumbSeparator />}
                    <BreadcrumbItem>
                      {breadcrumb.to ? (
                        <BreadcrumbLink asChild>
                          <Link to={breadcrumb.to}>{breadcrumb.label}</Link>
                        </BreadcrumbLink>
                      ) : (
                        <BreadcrumbPage>{breadcrumb.label}</BreadcrumbPage>
                      )}
                    </BreadcrumbItem>
                  </Fragment>
                ))}
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <main>{props.children}</main>
      </SidebarInset>
    </>
  ) : (
    <div className="flex h-screen w-full items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
