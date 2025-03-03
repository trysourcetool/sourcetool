import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
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
import { useToast } from '@/hooks/use-toast';
import { useDispatch } from '@/store';
import { usersStore } from '@/store/modules/users';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useNavigate, useSearchParams } from 'react-router';
import { object, string, boolean } from 'zod';
import type { z } from 'zod';

export type SearchParams = {
  isUserExists: string; // 'true' or 'false'
  token: string;
  email: string;
};

export default function InvitationSignupFollowup() {
  const dispatch = useDispatch();
  const [searchParams] = useSearchParams();
  const isUserExists = searchParams.get('isUserExists');
  const email = searchParams.get('email');
  const token = searchParams.get('token');

  const navigate = useNavigate();
  const { toast } = useToast();
  const { t } = useTranslation('common');
  console.log({ token, email, isUserExists });

  const schema = object({
    firstName: string({
      required_error: t('zod_errors_firstName_required'),
    }),
    lastName: string({
      required_error: t('zod_errors_lastName_required'),
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
    agreeToTerms: boolean({
      required_error: t('zod_errors_agreeToTerms_required'),
    }),
  }).superRefine(({ password, passwordConfirmation }, ctx) => {
    if (password !== passwordConfirmation) {
      ctx.addIssue({
        code: 'custom',
        message: t('zod_errors_passwordConfirmation_match'),
        path: ['passwordConfirmation'],
      });
    }
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit(async (data) => {
    if (!token) {
      toast({
        title: t('routes_invitation_signup_followup_toast_invalid_token'),
        description: t('routes_invitation_signup_followup_toast_try_again'),
        variant: 'destructive',
      });
      return;
    }
    const resultAction = await dispatch(
      usersStore.asyncActions.invitationsSignup({
        data: {
          invitationToken: token,
          ...data,
        },
      }),
    );
    if (
      usersStore.asyncActions.invitationsSignup.fulfilled.match(resultAction)
    ) {
      const result = await dispatch(usersStore.asyncActions.getUsersMe());
      if (usersStore.asyncActions.getUsersMe.fulfilled.match(result)) {
        if (result.payload.user.organization) {
          navigate('/');
          toast({
            title: t('routes_invitation_signup_followup_toast_success'),
          });
        }
      }
    } else {
      toast({
        title: t('routes_invitation_signup_followup_toast_error'),
        description: t('routes_invitation_signup_followup_toast_try_again'),
        variant: 'destructive',
      });
    }
  });

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>
              {t('routes_invitation_signup_followup_title')}
            </CardTitle>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <div className="flex items-start gap-3">
              <FormField
                control={form.control}
                name="firstName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {t('routes_invitation_signup_followup_first_name_label')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t(
                          'routes_invitation_signup_followup_first_name_placeholder',
                        )}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="lastName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {t('routes_invitation_signup_followup_last_name_label')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t(
                          'routes_invitation_signup_followup_last_name_placeholder',
                        )}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {t('routes_invitation_signup_followup_password_label')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t(
                        'routes_invitation_signup_followup_password_placeholder',
                      )}
                      {...field}
                      type="password"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="passwordConfirmation"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {t(
                      'routes_invitation_signup_followup_confirm_password_label',
                    )}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t(
                        'routes_invitation_signup_followup_password_placeholder',
                      )}
                      {...field}
                      type="password"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="agreeToTerms"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className="flex items-start gap-2">
                      <Checkbox
                        name="agreeToTerms"
                        checked={field.value}
                        onCheckedChange={() => field.onChange(!field.value)}
                      />
                      <label
                        htmlFor="terms1"
                        className="text-sm font-normal leading-4 text-foreground"
                        dangerouslySetInnerHTML={{
                          __html: t(
                            'routes_invitation_signup_followup_terms_text',
                          ).replace(
                            /<a>/g,
                            '<a href="#" target="_blank" class="underline">',
                          ),
                        }}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="w-full">
              {t('routes_invitation_signup_followup_continue_button')}
            </Button>
          </form>
        </Card>
      </Form>
    </div>
  );
}
