import { PageHeader } from '@/components/common/page-header';

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useDispatch, useSelector } from '@/store';
import { pagesStore } from '@/store/modules/pages';
import { AlertCircle, File } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { Link, createFileRoute } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import { CodeBlock } from '@/components/common/code-block';
import { usersStore } from '@/store/modules/users';
import { apiKeysStore } from '@/store/modules/apiKeys';
import { environmentsStore } from '@/store/modules/environments';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from '@/components/ui/select';
import { SelectValue } from '@radix-ui/react-select';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';

export default function Pages() {
  const isInitialLoading = useRef(false);
  const [isInitialLoaded, setIsInitialLoaded] = useState(false);
  const [selectedEnvironmentId, setSelectedEnvironmentId] = useState<
    string | null
  >(null);
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const pages = useSelector(pagesStore.selector.getPermissionPages);
  const user = useSelector(usersStore.selector.getUserMe);
  const devKey = useSelector(apiKeysStore.selector.getDevKey);
  const apiKeys = useSelector(apiKeysStore.selector.getApiKeys);
  const environments = useSelector(environmentsStore.selector.getEnvironments);

  const selectedApiKey = useMemo(() => {
    if (!selectedEnvironmentId) {
      return null;
    }
    if (
      environments.find((e) => e.id === selectedEnvironmentId)?.slug ===
      'development'
    ) {
      return devKey;
    }
    return (
      apiKeys.find(
        (apiKey) => apiKey.environment.id === selectedEnvironmentId,
      ) ?? null
    );
  }, [apiKeys, devKey, environments, selectedEnvironmentId]);

  console.log({ selectedApiKey });

  // TODO: Consider using redux-persist if localStorage is used frequently
  const setLocalStorageSelectedEnvironmentId = (environmentId: string) => {
    localStorage.setItem('selectedEnvironmentId', environmentId);
  };

  const getLocalStorageSelectedEnvironmentId = (): string | null => {
    const environmentId = localStorage.getItem('selectedEnvironmentId');
    return environmentId || null;
  };

  const handleSelectEnvironment = async (environmentId: string) => {
    setSelectedEnvironmentId(environmentId);
    setLocalStorageSelectedEnvironmentId(environmentId);

    await dispatch(pagesStore.asyncActions.listPages({ environmentId }));
  };

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_pages') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        const resultActions = await Promise.all([
          dispatch(environmentsStore.asyncActions.listEnvironments()),
          dispatch(apiKeysStore.asyncActions.listApiKeys()),
        ]);
        if (
          environmentsStore.asyncActions.listEnvironments.fulfilled.match(
            resultActions[0],
          )
        ) {
          const localStorageEnvironmentId =
            getLocalStorageSelectedEnvironmentId();
          let environmentId =
            localStorageEnvironmentId ||
            resultActions[0].payload.environments[0].id;
          const hasEnvironmentId = resultActions[0].payload.environments.some(
            (e) => e.id === environmentId,
          );
          if (!hasEnvironmentId) {
            environmentId = resultActions[0].payload.environments[0].id;
          }
          console.log({ environmentId });
          setSelectedEnvironmentId(environmentId);
          if (
            !localStorageEnvironmentId ||
            localStorageEnvironmentId !== environmentId
          ) {
            setLocalStorageSelectedEnvironmentId(environmentId);
          }

          await Promise.all([
            dispatch(
              pagesStore.asyncActions.listPages({
                environmentId,
              }),
            ),
            dispatch(apiKeysStore.asyncActions.listApiKeys()),
          ]);
        }
        isInitialLoading.current = false;
        setIsInitialLoaded(true);
      })();
    }
  }, [dispatch]);

  console.log({ selectedEnvironmentId, environments });

  return (
    <div>
      <PageHeader label={t('routes_pages_page_header')} />
      <div className="flex flex-col gap-2.5 px-4 py-6 md:px-6">
        <div className="flex flex-col items-start justify-between gap-3 md:flex-row md:items-center">
          <div className="flex gap-2 text-lg font-bold">
            {pages.length} Pages in
            <div className="flex items-center gap-2">
              <div
                className="size-3 rounded-full"
                style={{
                  backgroundColor: environments.find(
                    (e) => e.id === selectedEnvironmentId,
                  )?.color,
                }}
              />
              {environments.find((e) => e.id === selectedEnvironmentId)?.name}
            </div>
          </div>

          <div className="w-full md:max-w-72">
            {selectedEnvironmentId && (
              <Select
                value={selectedEnvironmentId ?? ''}
                onValueChange={handleSelectEnvironment}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a environment" />
                </SelectTrigger>
                <SelectContent>
                  {environments.map((environment) => (
                    <SelectItem key={environment.id} value={environment.id}>
                      {environment.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
        </div>
        {isInitialLoaded && pages.length === 0 && (
          <div className="flex w-full flex-col gap-4 rounded-md md:border md:p-6">
            <h2 className="text-xl font-bold">
              {t('routes_pages_placeholder_title')}
            </h2>
            <p className="text-sidebar-foreground font-normal">
              {t('routes_pages_placeholder_description')}
            </p>
            {!selectedApiKey?.key && (
              <Alert
                variant="destructive"
                className="border-destructive bg-destructive/10"
              >
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>No API Key Found for This Environment</AlertTitle>
                <AlertDescription>
                  <p>
                    Create an API key for this environment to access the setup
                    code and get started with your integration.
                  </p>
                  <Button variant="destructive" asChild>
                    <Link to={'/apiKeys/new'}>Edit API Key</Link>
                  </Button>
                </AlertDescription>
              </Alert>
            )}
            <CodeBlock
              code={`func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "${selectedApiKey?.key ?? 'your_api_key'}",
		Endpoint: "${user?.organization?.webSocketEndpoint}",
	})

	s.Page("/welcome", "Welcome to Sourcetool!", func(ui sourcetool.UIBuilder) error {
		ui.Markdown("## Hello ${user?.firstName}!")

		// Example:
		// name := ui.TextInput("Name")
		// email := ui.Email("Email")
		//
		// users, err := listUsers(ui.Context(), name, email)
		// if err != nil {
		//   return err
		// }
		//
		// ui.Table(users)

		return nil
	})

	if err := s.Listen(); err != nil {
		log.Fatal(err)
	}
}`}
              language="go"
            />

            <p className="text-sidebar-foreground font-normal">
              {t('routes_pages_placeholder_restart_server')}
            </p>
            <p className="text-sidebar-foreground font-normal">
              {t('routes_pages_placeholder_page_added')}
            </p>
            <p className="text-sidebar-foreground font-normal">
              {t('routes_pages_placeholder_documentation')}
            </p>
          </div>
        )}
        {pages.length > 0 && (
          <div className="rounded-md md:border md:p-4">
            <SidebarMenu>
              {/* TODO: Recursive processing and folder support */}
              {/* <Collapsible className="group/collapsible">
              <SidebarMenuItem>
                <CollapsibleTrigger asChild>
                  <SidebarMenuSubButton />
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    <SidebarMenuSubItem />
                  </SidebarMenuSub>
                </CollapsibleContent>
              </SidebarMenuItem>
            </Collapsible> */}
              {pages.map((page) => (
                <SidebarMenuItem key={page.id}>
                  <SidebarMenuButton asChild>
                    <Link to={`/pages/$`} params={{ _splat: page.route }}>
                      <File className="size-4" />
                      {page.name}
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </div>
        )}
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_auth/')({
  component: Pages,
});
