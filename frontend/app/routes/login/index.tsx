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
import { SocialButtonGoogle } from '@/components/common/social-button-google';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useDispatch, useSelector } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';

export default function Login() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const { t } = useTranslation('common');

  const isRequestMagicLinkWaiting = useSelector(
    (state) => state.auth.isRequestMagicLinkWaiting,
  );
  const isOauthGoogleAuthWaiting = useSelector(
    (state) => state.auth.isRequestGoogleAuthLinkWaiting,
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

  const onSubmit = form.handleSubmit(async (data) => {
    if (isOauthGoogleAuthWaiting) {
      return;
    }
    const resultAction = await dispatch(
      authStore.asyncActions.requestMagicLink({ data }),
    );
    if (authStore.asyncActions.requestMagicLink.fulfilled.match(resultAction)) {
      console.log('navigate', {
        to: '/login/emailSent',
        search: { email: data.email },
      });
      navigate({ to: '/login/emailSent', search: { email: data.email } });
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
      authStore.asyncActions.requestGoogleAuthLink(),
    );
    if (
      authStore.asyncActions.requestGoogleAuthLink.fulfilled.match(resultAction)
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
              {t('routes_login_title')}
            </CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              {t('routes_login_description')}
            </CardDescription>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <SocialButtonGoogle
              onClick={handleGoogleAuth}
              label={t('routes_login_google_button')}
            />

            <div className="relative flex items-center justify-center">
              <span className="text-foreground text-sm font-medium">
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
              disabled={isRequestMagicLinkWaiting}
            >
              {isRequestMagicLinkWaiting && (
                <Loader2 className="mr-2 size-4 animate-spin" />
              )}
              {t('routes_login_login_button')}
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

export const Route = createFileRoute('/_default/login/')({
  component: Login,
});
