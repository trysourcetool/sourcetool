import { AppAgentLayout } from '@/components/layout/app-agent-layout';
import { SidebarProvider } from '@/components/ui/sidebar';
import { createFileRoute, Outlet } from '@tanstack/react-router';

export default function AgentLayout() {
  return (
    <SidebarProvider>
      <AppAgentLayout>
        <Outlet />
      </AppAgentLayout>
    </SidebarProvider>
  );
}

export const Route = createFileRoute('/_agent')({
  component: AgentLayout,
});