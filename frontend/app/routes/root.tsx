import {
  Link,
  Outlet,
  createRootRoute,
  useNavigate,
  type ErrorComponentProps,
} from '@tanstack/react-router';
import i18n from '../i18n';
import '../tailwind.css';
import { ThemeProvider } from 'next-themes';
import { I18nextProvider } from 'react-i18next';
import { Provider } from 'react-redux';
import { configureStore } from '../store';
import { PlainNavbarLayout } from '../components/layout/plain-navbar-layout';
import { CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { AuthProvider } from '../hooks/use-auth';
import { Toaster } from '../components/ui/toaster';
import { WebSocketController } from '../components/common/websocket-controller';
import { DomainProvider } from '../components/common/domainProvider';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

const reduxStore = configureStore();

function Fallback(props: ErrorComponentProps) {
  const navigate = useNavigate();
  return (
    <PlainNavbarLayout>
      <div className="m-auto flex items-center justify-center">
        <div className="flex max-w-[374px] flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>Error: {props.error.name}</CardTitle>
            <CardDescription>{props.error.message}</CardDescription>
          </CardHeader>
          <Button
            onClick={() => {
              props.reset();
              navigate({ to: '/' });
            }}
          >
            Back to home
          </Button>
        </div>
      </div>
    </PlainNavbarLayout>
  );
}

export default function App() {
  return (
    <Provider store={reduxStore.store}>
      <ThemeProvider enableSystem={false} attribute="class">
        <AuthProvider>
          <DomainProvider>
            <I18nextProvider i18n={i18n}>
              <Outlet />
              <TanStackRouterDevtools position="bottom-right" />
            </I18nextProvider>
            <WebSocketController />
          </DomainProvider>
        </AuthProvider>
        <Toaster />
      </ThemeProvider>
    </Provider>
  );
}

export const Route = createRootRoute({
  component: App,
  errorComponent: Fallback,
  notFoundComponent: () => {
    return (
      <PlainNavbarLayout>
        <div className="m-auto flex items-center justify-center">
          <div className="flex max-w-[374px] flex-col gap-6 p-6">
            <CardHeader className="p-0">
              <CardTitle>Page not found</CardTitle>
              <CardDescription>
                The page you are looking for does not exist.
              </CardDescription>
            </CardHeader>
            <Button asChild>
              <Link to="/">Back to home</Link>
            </Button>
          </div>
        </div>
      </PlainNavbarLayout>
    );
  },
});
