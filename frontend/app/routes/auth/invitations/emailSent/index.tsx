import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Mail, ArrowLeft } from 'lucide-react';
import { useNavigate } from 'react-router';
import { useDispatch, useSelector } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { $path } from 'safe-routes';

export type SearchParams = {
  email: string;
  token: string;
};

export default function InvitationEmailSent() {
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email');
  const token = searchParams.get('token');
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { toast } = useToast();
  const isRequestInvitationMagicLinkWaiting = useSelector(
    (state) => state.auth.isRequestInvitationMagicLinkWaiting,
  );

  const handleResendEmail = async () => {
    if (!email || !token || isRequestInvitationMagicLinkWaiting) {
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
          <CardTitle className="text-2xl font-semibold text-foreground">
            {t('routes_login_email_sent_title')}
          </CardTitle>
          <div className="flex items-center gap-3 rounded-md border border-border p-3">
            <Mail className="h-5 w-5" />
            <CardDescription className="flex-1 text-sm text-muted-foreground">
              {t('routes_login_email_sent_description')}{' '}
              <span className="font-medium text-foreground">{email}</span>
            </CardDescription>
          </div>
        </CardHeader>

        <p className="text-center text-xs font-normal text-muted-foreground">
          {t('routes_login_email_sent_resend_text')}{' '}
          <button
            type="button"
            onClick={handleResendEmail}
            disabled={isRequestInvitationMagicLinkWaiting}
            className="cursor-pointer underline"
          >
            {t('routes_login_email_sent_resend_button')}
          </button>
        </p>
      </Card>

      <Button
        variant="secondary"
        className="fixed bottom-8 flex items-center gap-2"
        onClick={() =>
          navigate($path('/auth/invitations/login', { token, email }))
        }
      >
        <ArrowLeft className="h-4 w-4" />
        {t('routes_login_email_sent_back')}
      </Button>
    </div>
  );
}
