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
import { usersStore } from '@/store/modules/users';
import { useToast } from '@/hooks/use-toast';

export type SearchParams = {
  email: string;
};

export default function EmailSent() {
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email');
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { toast } = useToast();
  const isRequestMagicLinkWaiting = useSelector(
    (state) => state.users.isRequestMagicLinkWaiting,
  );

  const handleResendEmail = async () => {
    if (!email || isRequestMagicLinkWaiting) {
      return;
    }

    const resultAction = await dispatch(
      usersStore.asyncActions.requestMagicLink({
        data: { email },
      }),
    );
    if (usersStore.asyncActions.requestMagicLink.fulfilled.match(resultAction)) {
      toast({
        title: t('routes_signin_email_sent_resend_success'),
        description: t('routes_signin_email_sent_resend_success_description'),
      });
    } else {
      toast({
        title: t('routes_signin_toast_failed'),
        description: t('routes_signin_toast_failed_description'),
        variant: 'destructive',
      });
    }
  };

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Card className="flex w-full max-w-[384px] flex-col gap-6 p-6">
        <CardHeader className="space-y-6 p-0">
          <CardTitle className="text-2xl font-semibold text-foreground">
            {t('routes_signin_email_sent_title')}
          </CardTitle>
          <div className="flex items-center gap-3 rounded-md border border-border p-3">
            <Mail className="h-5 w-5 stroke-2 text-foreground" />
            <CardDescription className="text-sm text-muted-foreground">
              {t('routes_signin_email_sent_description')}{' '}
              <span className="font-medium text-foreground">{email}</span>
            </CardDescription>
          </div>
        </CardHeader>

        <p className="text-center text-xs text-muted-foreground">
          {t('routes_signin_email_sent_resend_text')}{' '}
          <Button
            variant="link"
            className="h-auto p-0 text-xs font-normal text-muted-foreground hover:text-foreground underline"
            onClick={handleResendEmail}
            disabled={isRequestMagicLinkWaiting}
          >
            {t('routes_signin_email_sent_resend_button')}
          </Button>
        </p>
      </Card>

      <Button
        variant="secondary"
        className="fixed bottom-8 flex items-center gap-2"
        onClick={() => navigate(-1)}
      >
        <ArrowLeft className="h-4 w-4" />
        {t('routes_signin_email_sent_back')}
      </Button>
    </div>
  );
} 