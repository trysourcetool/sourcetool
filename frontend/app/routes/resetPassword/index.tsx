import { useForm } from 'react-hook-form';
import { object, string } from 'zod';
import type { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslation } from 'react-i18next';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Link } from 'react-router';
import { $path } from 'safe-routes';

export default function ResetPassword() {
  const { t } = useTranslation('common');

  const schema = object({
    email: string({
      required_error: t('zod_errors_email_required'),
    }).email(t('zod_errors_email_format')),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit((data) => {
    console.log(data);
  });

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>{t('routes_reset_password_title')}</CardTitle>
            <CardDescription>
              {t('routes_reset_password_description')}
            </CardDescription>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {t('routes_reset_password_email_label')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('routes_reset_password_email_placeholder')}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" className="w-full">
              {t('routes_reset_password_button')}
            </Button>

            <p className="text-center text-sm font-normal text-foreground">
              {t('routes_reset_password_already_using')}{' '}
              <Link className="underline" to={$path('/signin')}>
                {t('routes_reset_password_signin_link')}
              </Link>
            </p>
          </form>
        </Card>
      </Form>
    </div>
  );
}
