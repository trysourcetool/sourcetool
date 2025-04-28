import { Button } from '@/components/ui/button';
import { CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useTranslation } from 'react-i18next';
import { ArrowRight } from 'lucide-react';
import { createFileRoute, Link } from '@tanstack/react-router';

export default function OnboardingComplete() {
  const { t } = useTranslation('common');
  return (
    <div className="m-auto flex w-full justify-center">
      <div className="flex max-w-[374px] flex-col gap-6 p-4 md:p-6">
        <CardHeader className="p-0">
          <CardTitle>{t('routes_onboarding_complete_title')}</CardTitle>
          <CardDescription>
            {t('routes_onboarding_complete_description')}
          </CardDescription>
        </CardHeader>
        <div className="flex">
          <Button asChild>
            <Link to={'/'}>
              {t('routes_onboarding_complete_button')}
              <ArrowRight />
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_default/onboarding/complete/')({
  component: OnboardingComplete,
});
