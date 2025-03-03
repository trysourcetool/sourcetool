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
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { Link, useNavigate } from 'react-router';
import { $path } from 'safe-routes';
import { object, string } from 'zod';
import type { z } from 'zod';

export default function EnvironmentNew() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const isCreateEnvironmentWaiting = useSelector(
    (state) => state.environments.isCreateEnvironmentWaiting,
  );

  const schema = object({
    name: string({
      required_error: t('zod_errors_name_required'),
    }),
    slug: string({
      required_error: t('zod_errors_slug_required'),
    }),
    color: string({
      required_error: t('zod_errors_color_required'),
    }).regex(/^#(?:[0-9a-fA-F]{3}){1,2}$/, t('zod_errors_color_format')),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
    defaultValues: {
      color: '#9333EA',
    },
  });

  const onSubmit = form.handleSubmit(async (data: Schema) => {
    if (isCreateEnvironmentWaiting) {
      return;
    }
    const resultAction = await dispatch(
      environmentsStore.asyncActions.createEnvironment({ data }),
    );

    if (
      environmentsStore.asyncActions.createEnvironment.fulfilled.match(
        resultAction,
      )
    ) {
      navigate(
        $path('/environments/:environmentId', {
          environmentId: resultAction.payload.environment.id,
        }),
      );
    } else {
      console.error(resultAction.error);
      toast({
        title: t('routes_environments_new_toast_create_failed'),
        variant: 'destructive',
      });
    }
  });

  useEffect(() => {
    setBreadcrumbsState?.([
      { label: t('breadcrumbs_environment'), to: $path('/environments') },
      { label: t('breadcrumbs_create_new') },
    ]);
  }, [setBreadcrumbsState, t]);

  return (
    <div>
      <PageHeader label={t('routes_environments_new_page_header')} />
      <Form {...form}>
        <form className="flex flex-col gap-6 p-6" onSubmit={onSubmit}>
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('routes_environments_new_name_label')}</FormLabel>
                <FormControl>
                  <Input
                    placeholder={t('routes_environments_new_name_placeholder')}
                    {...field}
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
                  {t('routes_environments_new_color_label')}
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

          <FormField
            control={form.control}
            name="slug"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('routes_environments_new_slug_label')}</FormLabel>
                <FormControl>
                  <Input
                    placeholder={t('routes_environments_new_slug_placeholder')}
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex flex-col justify-start gap-3 md:flex-row">
            <Button type="submit" disabled={isCreateEnvironmentWaiting}>
              {isCreateEnvironmentWaiting && (
                <Loader2 className="size-4 animate-spin" />
              )}
              {t('routes_environments_new_create_button')}
            </Button>
            <Button
              variant="outline"
              asChild
              disabled={isCreateEnvironmentWaiting}
            >
              <Link to={$path('/environments')}>
                {t('routes_environments_new_cancel_button')}
              </Link>
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
