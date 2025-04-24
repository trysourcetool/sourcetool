import { PageHeader } from '@/components/common/page-header';
import { object, string } from 'zod';
import type { z } from 'zod';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from '@/store';
import { useEffect, useRef } from 'react';
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

import { Button } from '@/components/ui/button';
import { createFileRoute, Link, useParams } from '@tanstack/react-router';
import { Copy, Loader2 } from 'lucide-react';
import { apiKeysStore } from '@/store/modules/apiKeys';
import { useToast } from '@/hooks/use-toast';
import dayjs from 'dayjs';

export default function ApiKeysEdit() {
  const isInitialLoading = useRef(false);
  const { toast } = useToast();
  const dispatch = useDispatch();
  const { apiKeyId } = useParams({
    from: '/_auth/apiKeys/$apiKeyId/',
  });
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');

  const apiKey = useSelector((state) =>
    apiKeysStore.selector.getApiKey(state, apiKeyId ?? ''),
  );
  const isUpdateApiKeyWaiting = useSelector(
    (state) => state.apiKeys.isUpdateApiKeyWaiting,
  );

  const schema = object({
    name: string({
      required_error: t('zod_errors_name_required'),
    })
      .trim()
      .min(1, t('zod_errors_name_required')),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit(async (data) => {
    console.log({ data });
    if (isUpdateApiKeyWaiting || !apiKeyId) {
      return;
    }
    const resultAction = await dispatch(
      apiKeysStore.asyncActions.updateApiKey({ apiKeyId, data }),
    );
    if (apiKeysStore.asyncActions.updateApiKey.fulfilled.match(resultAction)) {
      toast({
        title: t('routes_apikeys_edit_toast_updated'),
      });
    } else {
      toast({
        title: t('routes_apikeys_edit_toast_update_failed'),
        description: (resultAction.error as any)?.message ?? '',
        variant: 'destructive',
      });
    }
  });

  const onCopy = async () => {
    if (apiKey?.key) {
      try {
        await navigator.clipboard.writeText(apiKey.key);
        toast({
          title: t('routes_apikeys_edit_toast_copied'),
        });
      } catch (error) {
        toast({
          title: t('routes_apikeys_edit_toast_copy_failed'),
          description: (error as any)?.message ?? '',
          variant: 'destructive',
        });
      }
    }
  };

  useEffect(() => {
    setBreadcrumbsState?.([
      { label: t('breadcrumbs_api_keys'), to: '/apiKeys' },
      { label: t('breadcrumbs_edit_api_key') },
    ]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      if (apiKeyId) {
        (async () => {
          await dispatch(apiKeysStore.asyncActions.getApiKey({ apiKeyId }));
          isInitialLoading.current = false;
        })();
      }
    }
  }, [dispatch, apiKeyId]);

  useEffect(() => {
    if (apiKey) {
      form.reset({
        name: apiKey.name,
      });
    }
  }, [apiKey]);

  return (
    <div>
      <PageHeader label={t('routes_apikeys_edit_page_header')} />
      {apiKey && (
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
                  <FormLabel>{t('routes_apikeys_edit_name_label')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('routes_apikeys_edit_name_placeholder')}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormItem>
              <FormLabel>
                {t('routes_apikeys_edit_environment_label')}
              </FormLabel>
              <Input
                type="text"
                value={apiKey?.environment?.name ?? ''}
                disabled
              />
            </FormItem>

            <div className="flex justify-start gap-3">
              <Button type="submit">
                {isUpdateApiKeyWaiting && (
                  <Loader2 className="size-4 animate-spin" />
                )}
                {t('routes_apikeys_edit_save_button')}
              </Button>
              <Button variant="outline" asChild>
                <Link to={'/apiKeys'}>
                  {t('routes_apikeys_edit_cancel_button')}
                </Link>
              </Button>
            </div>

            <div className="bg-muted flex flex-col gap-4 px-6 py-4">
              <div className="flex flex-col gap-1">
                <p className="text-foreground text-lg font-bold">
                  {t('routes_apikeys_edit_key_title')}
                </p>
                <p className="text-muted-foreground text-sm font-normal">
                  {t('routes_apikeys_edit_created_on', {
                    date: dayjs
                      .unix(Number(apiKey?.createdAt))
                      .format('MMM D, YYYY'),
                  })}
                </p>
              </div>
              <div className="flex items-center gap-3">
                <div className="bg-input text-foreground flex min-w-0 flex-1 truncate rounded-md px-3 py-2.5 text-sm font-normal">
                  <p className="line-clamp-1 whitespace-break-spaces break-all">
                    {apiKey?.key ?? ''}
                  </p>
                </div>
                <Button
                  variant="outline"
                  size="icon"
                  type="button"
                  onClick={onCopy}
                >
                  <Copy />
                </Button>
              </div>
            </div>
          </form>
        </Form>
      )}
    </div>
  );
}

export const Route = createFileRoute('/_auth/apiKeys/$apiKeyId/')({
  component: ApiKeysEdit,
});
