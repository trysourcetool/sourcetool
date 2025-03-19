import { PageHeader } from '@/components/common/page-header';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';

import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useToast } from '@/hooks/use-toast';
import { useDispatch, useSelector } from '@/store';
import { environmentsStore } from '@/store/modules/environments';
import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2 } from 'lucide-react';
import { useEffect, useRef } from 'react';
import { useForm } from 'react-hook-form';
import { Link, useParams } from 'react-router';
import { $path } from 'safe-routes';
import { object, string } from 'zod';
import type { z } from 'zod';

export default function EnvironmentEdit() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const { toast } = useToast();
  const { environmentId } = useParams();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');

  const isCreateEnvironmentWaiting = useSelector(
    (state) => state.environments.isCreateEnvironmentWaiting,
  );
  const environment = useSelector((state) =>
    environmentsStore.selector.getEnvironment(state, environmentId ?? ''),
  );

  const schema = object({
    name: string({
      required_error: t('zod_errors_name_required'),
    }),
    color: string({
      required_error: t('zod_errors_color_required'),
    }).regex(/^#(?:[0-9a-fA-F]{3}){1,2}$/, t('zod_errors_color_format')),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit(async (data: Schema) => {
    if (!environmentId) {
      return;
    }
    const resultAction = await dispatch(
      environmentsStore.asyncActions.updateEnvironment({
        environmentId,
        data,
      }),
    );

    if (
      environmentsStore.asyncActions.updateEnvironment.fulfilled.match(
        resultAction,
      )
    ) {
      toast({
        title: t('routes_environments_edit_toast_updated'),
      });
    } else {
      toast({
        title: t('routes_environments_edit_toast_update_failed'),
        description: (resultAction.error as any)?.message ?? '',
        variant: 'destructive',
      });
    }
  });

  useEffect(() => {
    setBreadcrumbsState?.([
      { label: t('breadcrumbs_environment'), to: $path('/environments') },
      { label: t('breadcrumbs_edit_environment') },
    ]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current && environmentId) {
      isInitialLoading.current = true;
      (async () => {
        dispatch(
          environmentsStore.asyncActions.getEnvironment({ environmentId }),
        );
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch, environmentId]);

  useEffect(() => {
    if (environment) {
      form.reset({
        name: environment.name,
        color: environment.color,
      });
    }
  }, [environment]);

  return (
    <div>
      <PageHeader label={t('routes_environments_edit_page_header')} />
      {environment && (
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
                  <FormLabel>
                    {t('routes_environments_edit_name_label')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t(
                        'routes_environments_edit_name_placeholder',
                      )}
                      {...field}
                      defaultValue={environment.name}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="color"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {t('routes_environments_edit_color_label')}
                  </FormLabel>
                  <FormControl>
                    <div className="flex items-center gap-2">
                      <div
                        className="size-5 rounded-md"
                        style={{ backgroundColor: field.value }}
                      />
                      <Input {...field} />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormItem>
              <FormLabel>{t('routes_environments_edit_slug_label')}</FormLabel>
              <FormControl>
                <Input disabled value={environment.slug} />
              </FormControl>
            </FormItem>

            <div className="flex flex-row justify-start gap-3">
              <Button type="submit" disabled={isCreateEnvironmentWaiting}>
                {isCreateEnvironmentWaiting && (
                  <Loader2 className="size-4 animate-spin" />
                )}
                {t('routes_environments_edit_save_button')}
              </Button>
              <Button
                variant="outline"
                asChild
                disabled={isCreateEnvironmentWaiting}
              >
                <Link to={$path('/environments')}>
                  {t('routes_environments_edit_cancel_button')}
                </Link>
              </Button>
            </div>
          </form>
        </Form>
      )}
    </div>
  );
}
