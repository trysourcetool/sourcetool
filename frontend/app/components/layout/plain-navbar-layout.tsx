import type { PropsWithChildren } from 'react';
import { Logo } from '../common/logo';
import { ModeToggle } from '../common/mode-toggle';

export function PlainNavbarLayout(props: PropsWithChildren) {
  return (
    <div className="relative h-svh">
      <header className="sticky inset-x-0 top-0 z-50 flex h-16 items-center justify-between border-b bg-background px-6 shadow-xs">
        <Logo />
        <ModeToggle />
      </header>
      <main className="flex min-h-[calc(100svh-64px)] flex-col px-4 py-6 md:px-6 md:py-6">
        {props.children}
      </main>
    </div>
  );
}
