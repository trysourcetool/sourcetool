import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

export type SearchParams = {
  email: string;
};

export default function EmailSent() {
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email');
  const { t } = useTranslation('common');

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
        <CardHeader className="p-0">
          <CardTitle>{t('routes_signin_email_sent_title')}</CardTitle>
          <CardDescription>
            {t('routes_signin_email_sent_description', { email })}
          </CardDescription>
        </CardHeader>
      </Card>
    </div>
  );
} 