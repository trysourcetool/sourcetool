import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams, useNavigate } from 'react-router';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { useDispatch } from '@/store';
import { usersStore } from '@/store/modules/users';
import { useToast } from '@/hooks/use-toast';
import { $path } from 'safe-routes';

export type SearchParams = {
  token: string;
};

export default function MagicLinkAuth() {
  const isInitialLoading = useRef(false);
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { toast } = useToast();
  const { t } = useTranslation('common');

  useEffect(() => {
    if (isInitialLoading.current) {
      return;
    }
    isInitialLoading.current = true;

    const authenticate = async () => {
      if (!token) {
        toast({
          title: t('routes_signin_magic_link_toast_invalid_token'),
          description: t('routes_signin_magic_link_toast_try_again'),
          variant: 'destructive',
        });
        navigate($path('/signin'));
        return;
      }

      const resultAction = await dispatch(
        usersStore.asyncActions.authenticateWithMagicLink({
          data: { token },
        }),
      );

      if (usersStore.asyncActions.authenticateWithMagicLink.fulfilled.match(resultAction)) {
        const result = resultAction.payload;
        
        if (result.isNewUser) {
          navigate($path('/signup/followup', { token: result.token }));
        } else {
          const saveAuthResult = await dispatch(
            usersStore.asyncActions.saveAuth({
              authUrl: result.authUrl,
              data: { token: result.token },
            }),
          );

          if (usersStore.asyncActions.saveAuth.fulfilled.match(saveAuthResult)) {
            window.location.replace(saveAuthResult.payload.redirectUrl);
          } else {
            const meResult = await dispatch(usersStore.asyncActions.getUsersMe());
            if (usersStore.asyncActions.getUsersMe.fulfilled.match(meResult)) {
              if (!meResult.payload.user.organization) {
                navigate($path('/organizations/new'));
              }
            }
          }
        }
      } else {
        toast({
          title: t('routes_signin_magic_link_toast_error'),
          description: t('routes_signin_magic_link_toast_try_again'),
          variant: 'destructive',
        });
        navigate($path('/signin'));
      }
    };

    authenticate();
  }, [dispatch, navigate, toast, t]);

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
        <CardHeader className="p-0">
          <CardTitle>{t('routes_signin_magic_link_title')}</CardTitle>
          <CardDescription>
            {t('routes_signin_magic_link_description')}
          </CardDescription>
        </CardHeader>
      </Card>
    </div>
  );
} 