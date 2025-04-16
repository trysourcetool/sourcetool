import { useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router';
import { useDispatch } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { $path } from 'safe-routes';
import { useTranslation } from 'react-i18next';
import { Loader2 } from 'lucide-react';

export type SearchParams = {
  token: string;
};

export default function InvitationMagicLinkAuthenticate() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { t } = useTranslation('common');
  const isInitialLoading = useRef(false);

  const token = searchParams.get('token');

  useEffect(() => {
    if (isInitialLoading.current) {
      return;
    }
    isInitialLoading.current = true;

    const authenticate = async () => {
      if (!token) {
        toast({
          title: t('routes_login_toast_failed'),
          description: t('routes_login_toast_failed_description'),
          variant: 'destructive',
        });
        navigate($path('/login'));
        return;
      }

      const resultAction = await dispatch(
        authStore.asyncActions.authenticateWithInvitationMagicLink({
          data: { token },
        }),
      );

      if (
        authStore.asyncActions.authenticateWithInvitationMagicLink.fulfilled.match(
          resultAction,
        )
      ) {
        const result = resultAction.payload;

        if (result.isNewUser) {
          navigate(
            $path('/auth/invitations/signup/followup', { token: result.token }),
          );
        } else {
          const saveAuthResult = await dispatch(
            authStore.asyncActions.saveAuth({
              authUrl: result.authUrl,
              data: { token: result.token },
            }),
          );

          if (!authStore.asyncActions.saveAuth.fulfilled.match(saveAuthResult)) {
            throw new Error(t('routes_auth_invitations_magic_link_toast_save_auth_failed_desc' as any));
          }

          window.location.replace(saveAuthResult.payload.redirectUrl);
        }
      } else {
        toast({
          title: t('routes_login_toast_failed'),
          description: t('routes_login_toast_failed_description'),
          variant: 'destructive',
        });
        navigate($path('/login'));
      }
    };

    authenticate();
  }, [dispatch, token, navigate, toast, t]);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
