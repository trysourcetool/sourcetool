import * as React from 'react';
import { Outlet, createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/agents')({
  component: Agents,
});

function Agents() {
  return <Outlet />;
}