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
import { SocialButtonGoogle } from '@/components/common/social-button-google';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Link, useNavigate } from 'react-router';
import { $path } from 'safe-routes';
import { useDispatch, useSelector } from '@/store';
import { usersStore } from '@/store/modules/users';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';
import { useAuth } from '@/hooks/use-auth';

export default function Signin() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const { t } = useTranslation('common');

  const { isSourcetoolDomain, subDomain } = useAuth();

  const isRequestMagicLinkWaiting = useSelector((state) => state.users.isRequestMagicLinkWaiting);
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
  });

  const onSubmit = form.handleSubmit(async (data) => {
    if (isOauthGoogleAuthWaiting) {
      return;
    }
    const resultAction = await dispatch(
      usersStore.asyncActions.requestMagicLink({ data }),
    );
    if (usersStore.asyncActions.requestMagicLink.fulfilled.match(resultAction)) {
      navigate($path('/signin/emailSent', { email: data.email }));
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
        <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>{t('routes_signin_title')}</CardTitle>
            <CardDescription>{t('routes_signin_description')}</CardDescription>
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
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="w-full" disabled={isRequestMagicLinkWaiting}>
              {isRequestMagicLinkWaiting && <Loader2 className="size-4 animate-spin" />}
              {t('routes_signin_login_button')}
            </Button>

            <SocialButtonGoogle
              onClick={handleGoogleAuth}
              label={t('routes_signin_google_button')}
            />

            {((isSourcetoolDomain && subDomain === 'auth') ||
              !isSourcetoolDomain) && (
              <p className="text-center text-sm font-normal text-foreground">
                {t('routes_signin_no_account')}{' '}
                <Link className="underline" to={$path('/signup')}>
                  {t('routes_signin_signup_link')}
                </Link>
              </p>
            )}
          </form>
        </Card>
      </Form>
    </div>
  );
}
