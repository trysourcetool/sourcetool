import { PlainNavbarLayout } from '@/components/layout/plain-navbar-layout';
import { Outlet } from 'react-router';

export default function AuthLayout() {
  console.log('AuthLayout');
  return (
    <PlainNavbarLayout>
      <Outlet />
    </PlainNavbarLayout>
  );
}
