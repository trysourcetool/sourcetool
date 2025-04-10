import { api } from '@/api';
import { CodeBlock } from '@/components/common/code-block';
import { PageHeader } from '@/components/common/page-header';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { useToast } from '@/hooks/use-toast';
import { useDispatch, useSelector } from '@/store';
import { apiKeysStore } from '@/store/modules/apiKeys';
import { Loader2 } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router';
import { $path } from 'safe-routes';
import { usersStore } from '@/store/modules/users';

export default function Onboarding() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const { t } = useTranslation('common');
  const devKey = useSelector(apiKeysStore.selector.getDevKey);
  const user = useSelector(usersStore.selector.getMe);
  const [isCheckingConnection, setIsCheckingConnection] = useState(false);

  const handleCheckConnection = async () => {
    setIsCheckingConnection(true);
    try {
      await api.hostInstances.getHostInstancePing();
      navigate($path('/onboarding/complete'));
    } catch (error: any) {
      toast({
        title: t('routes_onboarding_toast_error_title'),
        variant: 'destructive',
        description: error.detail,
      });
    } finally {
      setIsCheckingConnection(false);
    }
  };

  useEffect(() => {
    if (!isInitialLoading.current) {
      (async () => {
        isInitialLoading.current = true;
        await dispatch(apiKeysStore.asyncActions.listApiKeys());
        isInitialLoading.current = false;
      })();
    }
  }, []);

  return (
    devKey && (
      <>
        <div className="m-auto flex w-full justify-center pb-14">
          <div className="w-full max-w-xl p-4 md:p-6">
            <div className="p-6">
              <PageHeader
                label={t('routes_onboarding_page_header')}
                border={false}
                description={t('routes_onboarding_page_description')}
              />
            </div>
            <Card className="flex flex-col gap-4 border-none px-4 py-0 md:px-6">
              <CardHeader className="p-0 pt-4">
                <CardTitle>{t('routes_onboarding_step1_title')}</CardTitle>
                <CardDescription>
                  {t('routes_onboarding_step1_description')}
                </CardDescription>
              </CardHeader>
              <CodeBlock
                code="go get github.com/trysourcetool/sourcetool-go"
                language="text"
              />

              <CardHeader className="p-0 pt-4">
                <CardTitle>{t('routes_onboarding_step2_title')}</CardTitle>
                <CardDescription>
                  {t('routes_onboarding_step2_description')}
                </CardDescription>
              </CardHeader>
              <CodeBlock
                code={`func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "${devKey.key}",
		Endpoint: "${user?.organization?.webSocketEndpoint}"
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
              <CardHeader className="p-0 pt-4">
                <CardTitle>{t('routes_onboarding_step3_title')}</CardTitle>
                <CardDescription>
                  {t('routes_onboarding_step3_description')}
                </CardDescription>
              </CardHeader>

              <CardHeader className="p-0 pt-4">
                <CardTitle>{t('routes_onboarding_verify_title')}</CardTitle>
                <CardDescription>
                  {t('routes_onboarding_verify_description')
                    .replace(
                      '<bold>',
                      '<span class="font-bold text-foreground">',
                    )
                    .replace('</bold>', '</span>')}
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </div>

        <div className="fixed inset-x-0 bottom-0 z-50 flex items-center justify-between border-t bg-background px-6 py-4">
          <Link
            to={$path('/')}
            className="text-sm font-normal text-muted-foreground underline"
          >
            {t('routes_onboarding_skip_button')}
          </Link>
          <Button onClick={handleCheckConnection}>
            {isCheckingConnection && (
              <Loader2 className="size-4 animate-spin" />
            )}
            {t('routes_onboarding_check_connection_button')}
          </Button>
        </div>
      </>
    )
  );
}
