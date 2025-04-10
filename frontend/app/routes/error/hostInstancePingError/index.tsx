import { CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';
import { ArrowLeft } from 'lucide-react';

export default function HostInstancePingError() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="m-auto flex items-center justify-center">
      <div className="flex max-w-[374px] flex-col gap-6 p-6 text-center">
        <CardHeader className="p-0">
          <CardTitle>{t('routes_error_host_instance_ping_title')}</CardTitle>
          <CardDescription className="whitespace-pre-wrap">
            {t('routes_error_host_instance_ping_description')}
          </CardDescription>
        </CardHeader>
        <div className="flex justify-center gap-4">
          <Button variant="outline" onClick={() => navigate('/')}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            {t('routes_error_host_instance_ping_back_button')}
          </Button>
          <Button asChild>
            <a href="mailto:a.yoshida@trysourcetool.com">
              {t('routes_error_host_instance_ping_support_button')}
            </a>
          </Button>
        </div>
      </div>
    </div>
  );
}
