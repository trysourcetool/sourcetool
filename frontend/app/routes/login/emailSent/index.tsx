import { useTranslation } from 'react-i18next';
import { createFileRoute, useSearch } from '@tanstack/react-router';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Mail, ArrowLeft } from 'lucide-react';
import { useNavigate } from '@tanstack/react-router';
import { useDispatch, useSelector } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { zodValidator } from '@tanstack/zod-adapter';
import { object, string } from 'zod';

export default function EmailSent() {
  const search = useSearch({ from: '/_default/login/emailSent/' });
  const email = search.email;
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { toast } = useToast();
  const isRequestMagicLinkWaiting = useSelector(
    (state) => state.auth.isRequestMagicLinkWaiting,
  );

  const handleResendEmail = async () => {
    if (!email || isRequestMagicLinkWaiting) {
      return;
    }

    const resultAction = await dispatch(
      authStore.asyncActions.requestMagicLink({
        data: { email },
      }),
    );
    if (authStore.asyncActions.requestMagicLink.fulfilled.match(resultAction)) {
      toast({
        title: t('routes_login_email_sent_resend_success'),
        description: t('routes_login_email_sent_resend_success_description'),
      });
    } else {
      toast({
        title: t('routes_login_toast_failed'),
        description: t('routes_login_toast_failed_description'),
        variant: 'destructive',
      });
    }
  };

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Card className="flex w-full max-w-[384px] flex-col gap-4 p-6">
        <CardHeader className="space-y-6 p-0">
          <CardTitle className="text-foreground text-2xl font-semibold">
            {t('routes_login_email_sent_title')}
          </CardTitle>
          <div className="border-border flex items-center gap-3 rounded-md border p-3">
            <Mail className="h-5 w-5" />
            <CardDescription className="text-muted-foreground flex-1 text-sm">
              {t('routes_login_email_sent_description')}{' '}
              <span className="text-foreground font-medium">{email}</span>
            </CardDescription>
          </div>
        </CardHeader>

        <p className="text-muted-foreground text-center text-xs font-normal">
          {t('routes_login_email_sent_resend_text')}{' '}
          <button
            type="button"
            onClick={handleResendEmail}
            disabled={isRequestMagicLinkWaiting}
            className="cursor-pointer underline"
          >
            {t('routes_login_email_sent_resend_button')}
          </button>
        </p>
      </Card>

      <div className="fixed bottom-8">
        <Button
          variant="secondary"
          className="cursor-pointer"
          onClick={() => navigate({ to: '/login' })}
        >
          <ArrowLeft className="h-4 w-4" />
          {t('routes_login_email_sent_back')}
        </Button>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_default/login/emailSent/')({
  component: EmailSent,
  validateSearch: zodValidator(
    object({
      email: string(),
    }),
  ),
});
