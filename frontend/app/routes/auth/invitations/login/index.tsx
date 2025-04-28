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
import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { useDispatch, useSelector } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';
import { SocialButtonGoogle } from '@/components/common/social-button-google';
import { zodValidator } from '@tanstack/zod-adapter';

export default function InvitationLogin() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const search = useSearch({
    from: '/_default/auth/invitations/login/',
  });
  const { t } = useTranslation('common');

  const token = search.token;
  const email = search.email;

  const isRequestInvitationMagicLinkWaiting = useSelector(
    (state) => state.auth.isRequestInvitationMagicLinkWaiting,
  );
  const isOauthGoogleAuthWaiting = useSelector(
    (state) => state.auth.isRequestInvitationGoogleAuthLinkWaiting,
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
      authStore.asyncActions.requestInvitationMagicLink({
        data: { invitationToken: token },
      }),
    );

    if (
      authStore.asyncActions.requestInvitationMagicLink.fulfilled.match(
        resultAction,
      )
    ) {
      navigate({
        to: '/auth/invitations/emailSent',
        search: { email: data.email, token },
      });
    } else {
      toast({
        title: t('routes_login_toast_failed'),
        description: t('routes_login_toast_failed_description'),
        variant: 'destructive',
      });
    }
  });

  const handleGoogleAuth = async () => {
    if (isOauthGoogleAuthWaiting) {
      return;
    }
    const resultAction = await dispatch(
      authStore.asyncActions.requestInvitationGoogleAuthLink({
        data: { invitationToken: token || '' },
      }),
    );
    if (
      authStore.asyncActions.requestInvitationGoogleAuthLink.fulfilled.match(
        resultAction,
      )
    ) {
      const url = resultAction.payload.authUrl;
      window.location.href = url;
    } else {
      toast({
        title: t('routes_login_toast_url_failed'),
        description: t('routes_login_toast_url_failed_description'),
        variant: 'destructive',
      });
    }
  };

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-[384px] flex-col gap-6 p-6">
          <CardHeader className="space-y-1.5 p-0">
            <CardTitle className="text-foreground text-2xl font-semibold">
              {t('routes_login_invitation_title')}
            </CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              {t('routes_login_invitation_description')}
            </CardDescription>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <SocialButtonGoogle
              onClick={handleGoogleAuth}
              label={t('routes_login_google_button')}
            />

            <div className="relative flex items-center justify-center">
              <div className="bg-border absolute inset-x-0 top-1/2 h-px" />
              <span className="bg-background text-foreground relative px-2 text-sm font-medium">
                {t('routes_login_or')}
              </span>
            </div>

            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input
                      placeholder={t('routes_login_email_placeholder')}
                      className="border-border h-[42px] text-sm"
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
              size="default"
              className="cursor-pointer"
              disabled={isRequestInvitationMagicLinkWaiting}
            >
              {isRequestInvitationMagicLinkWaiting && (
                <Loader2 className="mr-2 size-4 animate-spin" />
              )}
              {t('routes_login_invitation_login_button')}
            </Button>

            <p className="text-muted-foreground text-center text-xs">
              {t('routes_login_terms_text')}
            </p>
          </form>
        </Card>
      </Form>
    </div>
  );
}

export const Route = createFileRoute('/_default/auth/invitations/login/')({
  component: InvitationLogin,
  validateSearch: zodValidator(
    object({
      token: string(),
      email: string(),
    }),
  ),
});
