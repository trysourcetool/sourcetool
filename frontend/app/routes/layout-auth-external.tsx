import { AppExternalLayout } from '@/components/layout/app-external-layout';
import { SidebarProvider } from '@/components/ui/sidebar';
import { BreadcrumbsProvider } from '@/hooks/use-breadcrumbs';
import { Outlet } from 'react-router';

export default function AuthLayout() {
  return (
    <BreadcrumbsProvider>
      <SidebarProvider>
        <AppExternalLayout>
          <Outlet />
        </AppExternalLayout>
      </SidebarProvider>
    </BreadcrumbsProvider>
  );
}
