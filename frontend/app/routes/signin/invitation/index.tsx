import { object, string } from 'zod';
import type { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
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
import { useNavigate, useSearchParams } from 'react-router';
import { $path } from 'safe-routes';
import { useDispatch, useSelector } from '@/store';
import { usersStore } from '@/store/modules/users';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';
import { useEffect } from 'react';

export default function InvitationSignin() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { t } = useTranslation('common');

  const invitationToken = searchParams.get('token');

  const isRequestInvitationMagicLinkWaiting = useSelector(
    (state) => state.users.isRequestInvitationMagicLinkWaiting,
  );
  const isAuthenticateWithInvitationMagicLinkWaiting = useSelector(
    (state) => state.users.isAuthenticateWithInvitationMagicLinkWaiting,
  );

  const schema = object({
    email: string({
      required_error: t('zod_errors_email_required'),
    }).email(t('zod_errors_email_format')),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  useEffect(() => {
    if (invitationToken) {
      const requestInvitationMagicLink = async () => {
        const resultAction = await dispatch(
          usersStore.asyncActions.requestInvitationMagicLink({
            data: { invitationToken },
          }),
        );
        if (
          usersStore.asyncActions.requestInvitationMagicLink.fulfilled.match(
            resultAction,
          )
        ) {
          form.setValue('email', resultAction.payload.email);
        } else {
          toast({
            title: t('routes_signin_toast_failed'),
            description: t('routes_signin_toast_failed_description'),
            variant: 'destructive',
          });
        }
      };
      requestInvitationMagicLink();
    }
  }, [dispatch, invitationToken, form, toast, t]);

  const onSubmit = form.handleSubmit(async (data) => {
    if (!invitationToken) {
      return;
    }

    const resultAction = await dispatch(
      usersStore.asyncActions.authenticateWithInvitationMagicLink({
        data: { token: invitationToken },
      }),
    );

    if (
      usersStore.asyncActions.authenticateWithInvitationMagicLink.fulfilled.match(
        resultAction,
      )
    ) {
      const { isNewUser } = resultAction.payload;
      if (isNewUser) {
        navigate($path('/signup/invitation', { token: invitationToken }));
      } else {
        navigate($path('/signin/emailSent', { email: data.email }));
      }
    } else {
      toast({
        title: t('routes_signin_toast_failed'),
        description: t('routes_signin_toast_failed_description'),
        variant: 'destructive',
      });
    }
  });

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>{t('routes_signin_invitation_title')}</CardTitle>
            <CardDescription>
              {t('routes_signin_invitation_description')}
            </CardDescription>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('routes_signin_email_label')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('routes_signin_email_placeholder')}
                      {...field}
                      disabled={isRequestInvitationMagicLinkWaiting}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button
              type="submit"
              className="w-full"
              disabled={
                isRequestInvitationMagicLinkWaiting ||
                isAuthenticateWithInvitationMagicLinkWaiting
              }
            >
              {(isRequestInvitationMagicLinkWaiting ||
                isAuthenticateWithInvitationMagicLinkWaiting) && (
                <Loader2 className="size-4 animate-spin" />
              )}
              {t('routes_signin_invitation_login_button')}
            </Button>
          </form>
        </Card>
      </Form>
    </div>
  );
}
