import { useForm } from 'react-hook-form';
import { object, string } from 'zod';
import type { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslation } from 'react-i18next';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Loader2 } from 'lucide-react';
import type { ChangeEvent } from 'react';
import { useDispatch, useSelector } from '@/store';
import { organizationsStore } from '@/store/modules/organizations';
import { useDebouncedCallback } from 'use-debounce';
import { authStore } from '@/store/modules/auth';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import type { ErrorResponse } from '@/api/instance';
import { useToast } from '@/hooks/use-toast';

export default function OrganizationsNew() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { t } = useTranslation('common');

  const isCheckSubdomainAvailabilityWaiting = useSelector(
    (state) => state.organizations.isCheckSubdomainAvailabilityWaiting,
  );

  const schema = object({
    subdomain: string({
      required_error: t('zod_errors_subdomain_required'),
    }).min(3, t('zod_errors_subdomain_min')),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit(async (data) => {
    const resultAction = await dispatch(
      organizationsStore.asyncActions.createOrganization({
        data,
      }),
    );
    if (
      organizationsStore.asyncActions.createOrganization.fulfilled.match(
        resultAction,
      )
    ) {
      const result = await dispatch(authStore.asyncActions.obtainAuthToken());
      if (authStore.asyncActions.obtainAuthToken.fulfilled.match(result)) {
        const res = await dispatch(
          authStore.asyncActions.saveAuth({
            authUrl: result.payload.authUrl,
            data: { token: result.payload.token },
          }),
        );
        if (!authStore.asyncActions.saveAuth.fulfilled.match(res)) {
          throw new Error(
            t('routes_auth_magic_link_toast_save_auth_failed_desc' as any),
          );
        }
        window.location.replace(res.payload.redirectUrl);
      }
    }
  });

  const handleChange = async (event: ChangeEvent<HTMLInputElement>) => {
    if (!event.target.value.length || event.target.value.length < 3) {
      form.setError('subdomain', {
        message: t('zod_errors_subdomain_min'),
      });
      return;
    }
    const result = await dispatch(
      organizationsStore.asyncActions.checkSubdomainAvailability({
        subdomain: event.target.value,
      }),
    );
    if (
      organizationsStore.asyncActions.checkSubdomainAvailability.fulfilled.match(
        result,
      )
    ) {
      form.clearErrors('subdomain');
    } else {
      if (result.error && result.payload) {
        if ((result.payload as ErrorResponse).status === 401) {
          navigate({ to: '/login' });
          toast({
            title: t('routes_organizations_new_toast_auth_error'),
            description: t(
              'routes_organizations_new_toast_auth_error_description',
            ),
            variant: 'destructive',
          });
        }
      }
      form.setError('subdomain', {
        message: t('routes_organizations_new_subdomain_taken', {
          subdomain: event.target.value,
        }),
      });
    }
  };

  const debouncedHandleChange = useDebouncedCallback(handleChange, 300);

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>{t('routes_organizations_new_title')}</CardTitle>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <FormField
              control={form.control}
              name="subdomain"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {t('routes_organizations_new_subdomain_label')}
                  </FormLabel>
                  <FormControl>
                    <div className="flex items-center gap-2">
                      <Input
                        placeholder={t(
                          'routes_organizations_new_subdomain_placeholder',
                        )}
                        {...field}
                        onChange={(e) => {
                          field.onChange(e.target.value);
                          debouncedHandleChange(e);
                        }}
                      />
                      <p className="text-muted-foreground text-sm font-medium">
                        {t('routes_organizations_new_domain_suffix')}
                      </p>
                    </div>
                  </FormControl>
                  <FormMessage />
                  {isCheckSubdomainAvailabilityWaiting && (
                    <Loader2 className="size-4 animate-spin" />
                  )}
                </FormItem>
              )}
            />
            <Button type="submit" className="w-full">
              {t('routes_organizations_new_continue_button')}
            </Button>
          </form>
        </Card>
      </Form>
    </div>
  );
}

export const Route = createFileRoute('/_default/organizations/new/')({
  component: OrganizationsNew,
});
