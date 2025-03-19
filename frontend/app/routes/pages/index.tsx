import { PageHeader } from '@/components/common/page-header';

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useDispatch, useSelector } from '@/store';
import { pagesStore } from '@/store/modules/pages';
import { File } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { Link } from 'react-router';
import { useTranslation } from 'react-i18next';
import { CodeBlock } from '@/components/common/code-block';
import { usersStore } from '@/store/modules/users';

export default function Pages() {
  const isInitialLoading = useRef(false);
  const [isInitialLoaded, setIsInitialLoaded] = useState(false);
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');

  const pages = useSelector(pagesStore.selector.getPermissionPages);
  const user = useSelector(usersStore.selector.getMe);

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_pages') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        await dispatch(pagesStore.asyncActions.listPages());
        isInitialLoading.current = false;
        setIsInitialLoaded(true);
      })();
    }
  }, [dispatch]);

  return (
    <div>
      <PageHeader label={t('routes_pages_page_header')} />
      <div className="px-4 py-6 md:px-6">
        {isInitialLoaded && pages.length === 0 && (
          <div className="flex flex-col gap-6">
            <h2 className="text-xl font-bold">
              {t('routes_pages_placeholder_title')}
            </h2>
            <p className="font-normal text-sidebar-foreground">
              {t('routes_pages_placeholder_description')}
            </p>
            <CodeBlock
              code={`func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "your_api_key",
		Endpoint: "${user?.organization?.webSocketEndpoint}"
	})

	s.Page("Welcome to Sourcetool!", func(ui sourcetool.UIBuilder) error {
		ui.Markdown("## Hello {firstName}!")

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
            <p className="font-normal text-sidebar-foreground">
              {t('routes_pages_placeholder_restart_server')}
            </p>
            <p className="font-normal text-sidebar-foreground">
              {t('routes_pages_placeholder_page_added')}
            </p>
            <p className="font-normal text-sidebar-foreground">
              {t('routes_pages_placeholder_documentation')}
            </p>
          </div>
        )}
        {pages.length > 0 && (
          <div className="rounded-md border p-4">
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
                    <Link to={`/pages${page.route}`}>
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
