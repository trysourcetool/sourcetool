import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams, useNavigate } from 'react-router';
import { useDispatch } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { $path } from 'safe-routes';
import { Loader2 } from 'lucide-react';

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
          title: t('routes_login_magic_link_toast_invalid_token'),
          description: t('routes_login_magic_link_toast_try_again'),
          variant: 'destructive',
        });
        navigate($path('/login'));
        return;
      }

      const resultAction = await dispatch(
        authStore.asyncActions.authenticateWithMagicLink({
          data: { token },
        }),
      );

      if (
        authStore.asyncActions.authenticateWithMagicLink.fulfilled.match(
          resultAction,
        )
      ) {
        const result = resultAction.payload;

        if (result.isNewUser) {
          navigate($path('/signup/followup', { token: result.token }));
        } else {
          if (!result.hasOrganization) {
            navigate($path('/organizations/new'));
            return;
          }

          const saveAuthResult = await dispatch(
            authStore.asyncActions.saveAuth({
              authUrl: result.authUrl,
              data: { token: result.token },
            }),
          );

          if (!authStore.asyncActions.saveAuth.fulfilled.match(saveAuthResult)) {
            throw new Error(t('routes_auth_magic_link_toast_save_auth_failed_desc' as any));
          }

          window.location.replace(saveAuthResult.payload.redirectUrl);
        }
      } else {
        toast({
          title: t('routes_login_magic_link_toast_error'),
          description: t('routes_login_magic_link_toast_try_again'),
          variant: 'destructive',
        });
        navigate($path('/login'));
      }
    };

    authenticate();
  }, [dispatch, navigate, toast, t]);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
