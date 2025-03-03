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

export default function Settings() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');

  const user = useSelector(usersStore.selector.getMe);

  const isUpdateUserWaiting = useSelector(
    (state) => state.users.isUpdateUserWaiting,
  );
  const isUpdateUserPasswordWaiting = useSelector(
    (state) => state.users.isUpdateUserPasswordWaiting,
  );
  const isUsersSendUpdateEmailInstructionsWaiting = useSelector(
    (state) => state.users.isUsersSendUpdateEmailInstructionsWaiting,
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

  const passwordSchema = object({
    currentPassword: string({
      required_error: t('zod_errors_current_password_required'),
    }),
    password: string({
      required_error: t('zod_errors_password_required'),
    })
      .min(8, t('zod_errors_password_min'))
      .regex(
        /^(?=.*[A-Za-z])(?=.*\d)[a-zA-Z0-9!?_+*'"`#$%&\-^\\@;:,./=~|[\](){}<>]{8,}$/,
        t('zod_errors_password_format'),
      ),
    passwordConfirmation: string({
      required_error: t('zod_errors_passwordConfirmation_required'),
    }),
  }).superRefine(({ password, passwordConfirmation, currentPassword }, ctx) => {
    if (currentPassword === password) {
      ctx.addIssue({
        code: 'custom',
        message: t('zod_errors_password_same_as_current'),
        path: ['password'],
      });
      ctx.addIssue({
        code: 'custom',
        message: t('zod_errors_password_same_as_current'),
        path: ['passwordConfirmation'],
      });
    }
    if (password !== passwordConfirmation) {
      ctx.addIssue({
        code: 'custom',
        message: t('zod_errors_passwordConfirmation_match'),
        path: ['passwordConfirmation'],
      });
    }
  });

  type AccountSchema = z.infer<typeof accountSchema>;
  type EmailSchema = z.infer<typeof emailSchema>;
  type PasswordSchema = z.infer<typeof passwordSchema>;

  const accountForm = useForm<AccountSchema>({
    resolver: zodResolver(accountSchema),
  });
  console.log(accountForm.formState.errors, accountForm.getValues());

  const emailForm = useForm<EmailSchema>({
    resolver: zodResolver(emailSchema),
  });

  const passwordForm = useForm<PasswordSchema>({
    resolver: zodResolver(passwordSchema),
  });

  const handleSubmitAccount = accountForm.handleSubmit(
    async (data: AccountSchema) => {
      if (isUpdateUserWaiting) {
        return;
      }
      console.log({ data });
      const resultAction = await dispatch(
        usersStore.asyncActions.updateUser({
          data,
        }),
      );

      if (usersStore.asyncActions.updateUser.fulfilled.match(resultAction)) {
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
      if (isUsersSendUpdateEmailInstructionsWaiting) {
        return;
      }
      console.log({ data });
      const resultAction = await dispatch(
        usersStore.asyncActions.usersSendUpdateEmailInstructions({
          data,
        }),
      );

      if (
        usersStore.asyncActions.usersSendUpdateEmailInstructions.fulfilled.match(
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

  const handleSubmitPassword = passwordForm.handleSubmit(
    async (data: PasswordSchema) => {
      if (isUpdateUserPasswordWaiting) {
        return;
      }
      console.log({ data });
      const resultAction = await dispatch(
        usersStore.asyncActions.updateUserPassword({
          data,
        }),
      );

      if (
        usersStore.asyncActions.updateUserPassword.fulfilled.match(resultAction)
      ) {
        toast({
          title: t('routes_settings_toast_password_updated'),
        });
      } else {
        toast({
          title: t('routes_settings_toast_password_update_failed'),
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
        <div className="flex flex-col gap-6 p-6">
          <div className="flex flex-col gap-6 md:flex-row">
            <div className="flex flex-1 flex-col gap-2">
              <p className="text-xl font-semibold text-foreground">
                {t('routes_settings_name_title')}
              </p>
              <p className="text-sm font-normal text-muted-foreground">
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

          <Separator />

          <div className="flex flex-col gap-6 md:flex-row">
            <div className="flex flex-1 flex-col gap-2">
              <p className="text-xl font-semibold text-foreground">
                {t('routes_settings_email_title')}
              </p>
              <p className="text-sm font-normal text-muted-foreground">
                {t('routes_settings_email_description')}
              </p>
            </div>
            <Form {...emailForm}>
              <form
                className="flex flex-1 flex-col gap-6"
                onSubmit={handleSubmitEmail}
              >
                <div className="flex flex-col gap-1 rounded-md bg-muted px-4 py-3">
                  <p className="text-sm font-normal text-muted-foreground">
                    {t('routes_settings_current_email')}
                  </p>
                  <p className="text-base font-semibold text-foreground">
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

          <Separator />

          <div className="flex flex-col gap-6 md:flex-row">
            <div className="flex flex-1 flex-col gap-2">
              <p className="text-xl font-semibold text-foreground">
                {t('routes_settings_password_title')}
              </p>
              <p className="text-sm font-normal text-muted-foreground">
                {t('routes_settings_password_description')}
              </p>
            </div>
            <Form {...passwordForm}>
              <form
                className="flex flex-1 flex-col gap-6"
                onSubmit={handleSubmitPassword}
              >
                <FormField
                  control={passwordForm.control}
                  name="currentPassword"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('routes_settings_current_password_label')}
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder={t(
                            'routes_settings_current_password_placeholder',
                          )}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={passwordForm.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('routes_settings_new_password_label')}
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder={t(
                            'routes_settings_new_password_placeholder',
                          )}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={passwordForm.control}
                  name="passwordConfirmation"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('routes_settings_confirm_password_label')}
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder={t(
                            'routes_settings_confirm_password_placeholder',
                          )}
                          {...field}
                        />
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
        </div>
      )}
    </div>
  );
}
