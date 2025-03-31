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
import { SocialButtonGoogle } from '@/components/common/social-button-google';

export default function InvitationLogin() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { t } = useTranslation('common');

  const token = searchParams.get('token');
  const email = searchParams.get('email');

  const isRequestInvitationMagicLinkWaiting = useSelector(
    (state) => state.users.isRequestInvitationMagicLinkWaiting,
  );
  const isOauthGoogleAuthWaiting = useSelector(
    (state) => state.users.isOauthGoogleAuthWaiting,
  );

  const schema = object({
    email: string({
      required_error: t('zod_errors_email_required'),
    }).email(t('zod_errors_email_format')),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: email || '',
    },
  });

  const onSubmit = form.handleSubmit(async (data) => {
    if (!token) {
      return;
    }

    const resultAction = await dispatch(
      usersStore.asyncActions.requestInvitationMagicLink({
        data: { invitationToken: token },
      }),
    );

    if (usersStore.asyncActions.requestInvitationMagicLink.fulfilled.match(resultAction)) {
      navigate($path('/auth/invitations/emailSent', { email: data.email, token }));
    } else {
      toast({
        title: t('routes_signin_toast_failed'),
        description: t('routes_signin_toast_failed_description'),
        variant: 'destructive',
      });
    }
  });

  const handleGoogleAuth = async () => {
    if (isOauthGoogleAuthWaiting) {
      return;
    }
    const resultAction = await dispatch(
      usersStore.asyncActions.oauthGoogleAuthCodeUrl(),
    );
    if (
      usersStore.asyncActions.oauthGoogleAuthCodeUrl.fulfilled.match(
        resultAction,
      )
    ) {
      const url = resultAction.payload.url;
      window.location.href = url;
    } else {
      toast({
        title: t('routes_signin_toast_url_failed'),
        description: t('routes_signin_toast_url_failed_description'),
        variant: 'destructive',
      });
    }
  };

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-[384px] flex-col gap-6 p-6">
          <CardHeader className="space-y-1.5 p-0">
            <CardTitle className="text-2xl font-semibold text-foreground">
              {t('routes_signin_invitation_title')}
            </CardTitle>
            <CardDescription className="text-sm text-muted-foreground">
              {t('routes_signin_invitation_description')}
            </CardDescription>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <SocialButtonGoogle
              onClick={handleGoogleAuth}
              label={t('routes_signin_google_button')}
            />

            <div className="relative flex items-center justify-center">
              <div className="absolute inset-x-0 top-1/2 h-px bg-border" />
              <span className="relative bg-background px-2 text-sm font-medium text-foreground">
                {t('routes_signin_or')}
              </span>
            </div>

            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input
                      placeholder={t('routes_signin_email_placeholder')}
                      className="h-[42px] border-border text-sm"
                      {...field}
                      disabled={!!email}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button
              type="submit"
              className="h-[42px] w-full bg-primary text-sm font-medium text-primary-foreground hover:bg-primary/90"
              disabled={isRequestInvitationMagicLinkWaiting}
            >
              {isRequestInvitationMagicLinkWaiting && <Loader2 className="size-4 animate-spin mr-2" />}
              {t('routes_signin_invitation_login_button')}
            </Button>

            <p className="text-center text-xs text-muted-foreground">
              {t('routes_signin_terms_text')}
            </p>
          </form>
        </Card>
      </Form>
    </div>
  );
} 