import { AppPreviewLayout } from '@/components/layout/app-preview-layout';
import { SidebarProvider } from '@/components/ui/sidebar';
import { Outlet } from 'react-router';

export default function AuthLayout() {
  return (
    <SidebarProvider>
      <AppPreviewLayout>
        <Outlet />
      </AppPreviewLayout>
    </SidebarProvider>
  );
}
