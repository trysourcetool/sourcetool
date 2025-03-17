import {
  isRouteErrorResponse,
  Link,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useRouteError,
} from 'react-router';
import type { LinksFunction } from 'react-router';
import i18n from './i18n';

import styles from './tailwind.css?url';
import { ThemeProvider } from 'next-themes';
import { I18nextProvider, useTranslation } from 'react-i18next';
import { Provider } from 'react-redux';
import { configureStore } from './store';
import { Loader2 } from 'lucide-react';
import { PlainNavbarLayout } from './components/layout/plain-navbar-layout';
import { CardDescription, CardHeader, CardTitle } from './components/ui/card';
import { Button } from './components/ui/button';
import { $path } from 'safe-routes';
import { AuthProvider } from './hooks/use-auth';
import { Toaster } from './components/ui/toaster';
import { WebSocketController } from './components/common/websocket-controller';
import { DomainProvider } from './components/common/domainProvider';
import { NuqsAdapter } from 'nuqs/adapters/react-router/v7';

export const links: LinksFunction = () => [{ rel: 'stylesheet', href: styles }];

const reduxStore = configureStore();

export function Layout({ children }: { children: React.ReactNode }) {
  const { i18n } = useTranslation('common');
  console.log({ i18n });
  return (
    <html lang={i18n.language}>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        <Provider store={reduxStore.store}>
          <ThemeProvider enableSystem={false} attribute="class">
            <AuthProvider>
              <DomainProvider>
                {children}
                <WebSocketController />
              </DomainProvider>
            </AuthProvider>
            <Toaster />
          </ThemeProvider>
        </Provider>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export function ErrorBoundary() {
  const { i18n } = useTranslation('common');
  const error = useRouteError();
  console.log({ error });
  return (
    <html lang={i18n.language}>
      <head>
        <title>Oops!</title>
        <Meta />
        <Links />
      </head>
      <body>
        <PlainNavbarLayout>
          <div className="m-auto flex items-center justify-center">
            <div className="flex max-w-[374px] flex-col gap-6 p-6">
              <CardHeader className="p-0">
                <CardTitle>
                  {isRouteErrorResponse(error)
                    ? `${error.status}: ${error.statusText}`
                    : error instanceof Error
                      ? error.message
                      : 'Unknown Error'}
                </CardTitle>
                <CardDescription>{`We couldn't find the page you're looking for.
It might have been moved, deleted, or never existed in the first place.`}</CardDescription>
              </CardHeader>
              <Link to={$path('/')}>
                <Button>Back to home</Button>
              </Link>
            </div>
          </div>
        </PlainNavbarLayout>
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return (
    <NuqsAdapter>
      <I18nextProvider i18n={i18n}>
        <Outlet />
      </I18nextProvider>
    </NuqsAdapter>
  );
}

export function HydrateFallback() {
  return (
    <div className="flex h-screen items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
