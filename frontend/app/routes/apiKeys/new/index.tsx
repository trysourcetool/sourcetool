import { PageHeader } from '@/components/common/page-header';
import { object, string } from 'zod';
import type { z } from 'zod';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from '@/store';
import { environmentsStore } from '@/store/modules/environments';
import { useEffect } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Link, useNavigate } from 'react-router';
import { $path } from 'safe-routes';
import { apiKeysStore } from '@/store/modules/apiKeys';
import { Loader2 } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { usersStore } from '@/store/modules/users';

export default function ApiKeysNew() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t, i18n } = useTranslation('common');
  const account = useSelector(usersStore.selector.getUserMe);
  const environments = useSelector(environmentsStore.selector.getEnvironments);
  const isCreateApiKeyWaiting = useSelector(
    (state) => state.apiKeys.isCreateApiKeyWaiting,
  );

  console.log(i18n);

  const schema = object({
    name: string({
      required_error: t('zod_errors_name_required'),
    }),
    environmentId: string({
      required_error: t('zod_errors_environment_required'),
    }),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit(async (data) => {
    if (isCreateApiKeyWaiting) {
      return;
    }
    const resultAction = await dispatch(
      apiKeysStore.asyncActions.createApiKey(data),
    );
    if (apiKeysStore.asyncActions.createApiKey.fulfilled.match(resultAction)) {
      navigate(
        $path('/apiKeys/:apiKeyId', {
          apiKeyId: resultAction.payload.apiKey.id,
        }),
      );
    } else {
      toast({
        title: t('routes_apikeys_new_toast_create_failed'),
        variant: 'destructive',
      });
    }
  });

  useEffect(() => {
    setBreadcrumbsState?.([
      { label: t('breadcrumbs_api_keys'), to: $path('/apiKeys') },
      { label: t('breadcrumbs_create_new') },
    ]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    dispatch(environmentsStore.asyncActions.listEnvironments());
  }, [dispatch]);

  return (
    <div>
      <PageHeader label={t('routes_apikeys_new_page_header')} />
      <Form {...form}>
        <form
          className="flex flex-col gap-6 px-4 py-6 md:px-6"
          onSubmit={onSubmit}
        >
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('routes_apikeys_new_name_label')}</FormLabel>
                <FormControl>
                  <Input
                    placeholder={t('routes_apikeys_new_name_placeholder')}
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="environmentId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  {t('routes_apikeys_new_environment_label')}
                </FormLabel>
                <FormControl>
                  <Select onValueChange={field.onChange}>
                    <SelectTrigger>
                      <SelectValue
                        placeholder={t(
                          'routes_apikeys_new_environment_placeholder',
                        )}
                      />
                    </SelectTrigger>
                    <SelectContent>
                      {environments.map((environment) =>
                        (account?.role !== 'admin' &&
                          environment.slug === 'production') ||
                        environment.slug === 'development' ? null : (
                          <SelectItem
                            key={environment.id}
                            value={environment.id}
                          >
                            {environment.name}
                          </SelectItem>
                        ),
                      )}
                    </SelectContent>
                  </Select>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex flex-row justify-start gap-3">
            <Button type="submit" disabled={isCreateApiKeyWaiting}>
              {isCreateApiKeyWaiting && (
                <Loader2 className="size-4 animate-spin" />
              )}
              {t('routes_apikeys_new_create_button')}
            </Button>
            <Button variant="outline" asChild disabled={isCreateApiKeyWaiting}>
              <Link to={$path('/apiKeys')}>
                {t('routes_apikeys_new_cancel_button')}
              </Link>
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
