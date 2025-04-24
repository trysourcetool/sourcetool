import * as React from 'react';
import { Outlet, createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/users')({
  component: Users,
});

function Users() {
  return <Outlet />;
}
