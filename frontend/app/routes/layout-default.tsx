import { PlainNavbarLayout } from '@/components/layout/plain-navbar-layout';
import { createFileRoute, Outlet } from '@tanstack/react-router';

export default function DefaultLayout() {
  return (
    <PlainNavbarLayout>
      <Outlet />
    </PlainNavbarLayout>
  );
}

export const Route = createFileRoute('/_default')({
  component: DefaultLayout,
});
