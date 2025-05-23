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
import { useDispatch } from '@/store';
import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { useToast } from '@/hooks/use-toast';
import { authStore } from '@/store/modules/auth';
import { usersStore } from '@/store/modules/users';
import { ENVIRONMENTS } from '@/environments';
import { zodValidator } from '@tanstack/zod-adapter';

export default function Followup() {
  const dispatch = useDispatch();
  const search = useSearch({
    from: '/_default/signup/followup/',
  });
  const token = search.token;
  const navigate = useNavigate();
  const { toast } = useToast();
  const { t } = useTranslation('common');

  const schema = object({
    firstName: string({
      required_error: t('zod_errors_firstName_required'),
    }),
    lastName: string({
      required_error: t('zod_errors_lastName_required'),
    }),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit(async (data) => {
    if (!token) {
      toast({
        title: t('routes_signup_followup_toast_invalid_token'),
        description: t('routes_signup_followup_toast_try_again'),
        variant: 'destructive',
      });
      return;
    }
    const resultAction = await dispatch(
      authStore.asyncActions.registerWithMagicLink({
        data: {
          token,
          ...data,
        },
      }),
    );
    if (
      authStore.asyncActions.registerWithMagicLink.fulfilled.match(resultAction)
    ) {
      if (
        ENVIRONMENTS.IS_CLOUD_EDITION &&
        !resultAction.payload.hasOrganization
      ) {
        navigate({ to: '/organizations/new' });
        toast({
          title: t('routes_signup_followup_toast_success'),
          description: t('routes_signup_followup_toast_success_description'),
        });
      } else {
        await dispatch(usersStore.asyncActions.getMe());
        navigate({ to: '/' });
      }
    } else {
      toast({
        title: t('routes_signup_followup_toast_error'),
        description: t('routes_signup_followup_toast_try_again'),
        variant: 'destructive',
      });
    }
  });

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>{t('routes_signup_followup_title')}</CardTitle>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <div className="flex items-start gap-3">
              <FormField
                control={form.control}
                name="firstName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {t('routes_signup_followup_first_name_label')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t(
                          'routes_signup_followup_first_name_placeholder',
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
                      {t('routes_signup_followup_last_name_label')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t(
                          'routes_signup_followup_last_name_placeholder',
                        )}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <Button type="submit" className="w-full cursor-pointer">
              {t('routes_signup_followup_continue_button')}
            </Button>
          </form>
        </Card>
      </Form>
    </div>
  );
}

export const Route = createFileRoute('/_default/signup/followup/')({
  component: Followup,
  validateSearch: zodValidator(
    object({
      token: string(),
    }),
  ),
});
