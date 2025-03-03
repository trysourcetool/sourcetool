import { CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';

export default function HostInstancePingError() {
  const { t } = useTranslation();

  return (
    <div className="m-auto flex items-center justify-center">
      <div className="flex max-w-[374px] flex-col gap-6 p-6 text-center">
        <CardHeader className="p-0">
          <CardTitle>{t('routes_error_host_instance_ping_title')}</CardTitle>
          <CardDescription className="whitespace-pre-wrap">
            {t('routes_error_host_instance_ping_description')}
          </CardDescription>
        </CardHeader>
        <div className="flex justify-center">
          <Button>{t('routes_error_host_instance_ping_support_button')}</Button>
        </div>
      </div>
    </div>
  );
}
