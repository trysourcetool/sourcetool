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
import { Separator } from '@/components/ui/separator';
import { object, string } from 'zod';
import type { z } from 'zod';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { usersStore } from '@/store/modules/users';
import { useDispatch, useSelector } from '@/store';
import { useToast } from '@/hooks/use-toast';
import { createFileRoute } from '@tanstack/react-router';
export default function Settings() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');

  const user = useSelector(usersStore.selector.getUserMe);

  const isUpdateMeWaiting = useSelector(
    (state) => state.users.isUpdateMeWaiting,
  );
  const isSendUpdateMeEmailInstructionsWaiting = useSelector(
    (state) => state.users.isSendUpdateMeEmailInstructionsWaiting,
  );

  const accountSchema = object({
    firstName: string({
      required_error: t('zod_errors_firstName_required'),
    }),
    lastName: string({
      required_error: t('zod_errors_lastName_required'),
    }),
  });

  const emailSchema = object({
    email: string({
      required_error: t('zod_errors_email_required'),
    }).email(t('zod_errors_email_format')),
    emailConfirmation: string({
      required_error: t('zod_errors_email_confirmation_required'),
    }),
  }).superRefine(({ email, emailConfirmation }, ctx) => {
    if (email !== emailConfirmation) {
      ctx.addIssue({
        code: 'custom',
        message: t('zod_errors_email_match'),
        path: ['emailConfirmation'],
      });
    }
  });

  type AccountSchema = z.infer<typeof accountSchema>;
  type EmailSchema = z.infer<typeof emailSchema>;

  const accountForm = useForm<AccountSchema>({
    resolver: zodResolver(accountSchema),
  });
  console.log(accountForm.formState.errors, accountForm.getValues());

  const emailForm = useForm<EmailSchema>({
    resolver: zodResolver(emailSchema),
  });

  const handleSubmitAccount = accountForm.handleSubmit(
    async (data: AccountSchema) => {
      if (isUpdateMeWaiting) {
        return;
      }
      console.log({ data });
      const resultAction = await dispatch(
        usersStore.asyncActions.updateMe({
          data,
        }),
      );

      if (usersStore.asyncActions.updateMe.fulfilled.match(resultAction)) {
        toast({
          title: t('routes_settings_toast_user_updated'),
        });
      } else {
        toast({
          title: t('routes_settings_toast_user_update_failed'),
          variant: 'destructive',
        });
      }
    },
  );

  const handleSubmitEmail = emailForm.handleSubmit(
    async (data: EmailSchema) => {
      if (isSendUpdateMeEmailInstructionsWaiting) {
        return;
      }
      console.log({ data });
      const resultAction = await dispatch(
        usersStore.asyncActions.sendUpdateMeEmailInstructions({
          data,
        }),
      );

      if (
        usersStore.asyncActions.sendUpdateMeEmailInstructions.fulfilled.match(
          resultAction,
        )
      ) {
        toast({
          title: t('routes_settings_toast_email_sent'),
          description: t('routes_settings_toast_email_sent_description'),
        });
      } else {
        toast({
          title: t('routes_settings_toast_email_send_failed'),
          variant: 'destructive',
        });
      }
    },
  );

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_settings') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (user) {
      accountForm.reset({
        firstName: user.firstName,
        lastName: user.lastName,
      });
    }
  }, [user]);

  return (
    <div>
      <PageHeader label={t('routes_settings_page_header')} />
      {user && (
        <div className="flex flex-col gap-4 px-4 py-6 md:gap-6 md:px-6 md:py-6">
          <div className="flex flex-col gap-4 md:flex-row md:gap-6">
            <div className="flex flex-1 flex-col gap-2">
              <p className="text-foreground text-xl font-semibold">
                {t('routes_settings_name_title')}
              </p>
              <p className="text-muted-foreground text-sm font-normal">
                {t('routes_settings_name_description')}
              </p>
            </div>
            <Form {...accountForm}>
              <form
                className="flex flex-1 flex-col gap-6"
                onSubmit={handleSubmitAccount}
              >
                <FormField
                  control={accountForm.control}
                  name="firstName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('routes_settings_first_name_label')}
                      </FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={accountForm.control}
                  name="lastName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('routes_settings_last_name_label')}
                      </FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div>
                  <Button type="submit" className="cursor-pointer">
                    {t('routes_settings_save_button')}
                  </Button>
                </div>
              </form>
            </Form>
          </div>

          <div className="hidden md:block">
            <Separator />
          </div>

          <div className="flex flex-col gap-4 md:flex-row md:gap-6">
            <div className="flex flex-1 flex-col gap-2">
              <p className="text-foreground text-xl font-semibold">
                {t('routes_settings_email_title')}
              </p>
              <p className="text-muted-foreground text-sm font-normal">
                {t('routes_settings_email_description')}
              </p>
            </div>
            <Form {...emailForm}>
              <form
                className="flex flex-1 flex-col gap-6"
                onSubmit={handleSubmitEmail}
              >
                <div className="bg-muted flex flex-col gap-1 rounded-md px-4 py-3">
                  <p className="text-muted-foreground text-sm font-normal">
                    {t('routes_settings_current_email')}
                  </p>
                  <p className="text-foreground text-base font-semibold">
                    {user.email}
                  </p>
                </div>

                <FormField
                  control={emailForm.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('routes_settings_email_label')}</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          placeholder={t('routes_settings_email_placeholder')}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={emailForm.control}
                  name="emailConfirmation"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('routes_settings_email_confirmation_label')}
                      </FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          placeholder={t(
                            'routes_settings_email_confirmation_placeholder',
                          )}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div>
                  <Button type="submit" className="cursor-pointer">
                    {t('routes_settings_send_verification_email_button')}
                  </Button>
                </div>
              </form>
            </Form>
          </div>
        </div>
      )}
    </div>
  );
}

export const Route = createFileRoute('/_auth/settings/')({
  component: Settings,
});
