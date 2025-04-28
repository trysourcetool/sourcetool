import { AppPreviewLayout } from '@/components/layout/app-preview-layout';
import { SidebarProvider } from '@/components/ui/sidebar';
import { createFileRoute, Outlet } from '@tanstack/react-router';

export default function AuthLayout() {
  return (
    <SidebarProvider>
      <AppPreviewLayout>
        <Outlet />
      </AppPreviewLayout>
    </SidebarProvider>
  );
}

export const Route = createFileRoute('/_preview')({
  component: AuthLayout,
});
