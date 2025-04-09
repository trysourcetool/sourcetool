import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
// Use react-router provided by Remix
import { useSearchParams, useNavigate } from 'react-router'; 
import { useDispatch } from '@/store';
import { usersStore } from '@/store/modules/users';
import { useToast } from '@/hooks/use-toast';
import { $path } from 'safe-routes';
import { Loader2 } from 'lucide-react';

export default function GoogleAuthenticate() {
  const isProcessing = useRef(false);
  const [searchParams] = useSearchParams();
  const code = searchParams.get('code');
  const state = searchParams.get('state');
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { toast } = useToast();
  const { t } = useTranslation('common');

  useEffect(() => {
    if (isProcessing.current) {
      return;
    }
    isProcessing.current = true;

    const authenticateAndProceed = async () => {
      if (!code || !state) {
        toast({
          title: t('routes_auth_google_toast_missing_params_title' as any),
          description: t('routes_auth_google_toast_missing_params_desc' as any),
          variant: 'destructive',
        });
        navigate($path('/login'));
        return;
      }

      try {
        const authResultAction = await dispatch(
          usersStore.asyncActions.authenticateWithGoogle({
            data: { code, state },
          }),
        );

        if (!usersStore.asyncActions.authenticateWithGoogle.fulfilled.match(authResultAction)) {
          const errorPayload = (authResultAction.payload as any)?.error || {};
          const errorMessage = errorPayload.message || t('routes_auth_google_toast_auth_failed_desc' as any);
          throw new Error(errorMessage);
        }

        const authResult = authResultAction.payload;

        if (authResult.isNewUser) {
          const registerResultAction = await dispatch(
            usersStore.asyncActions.registerWithGoogle({
              data: {
                token: authResult.token,
                firstName: authResult.firstName || '-',
                lastName: authResult.lastName || '-',
              },
            }),
          );

          if (!usersStore.asyncActions.registerWithGoogle.fulfilled.match(registerResultAction)) {
            const errorPayload = (registerResultAction.payload as any)?.error || {};
            const errorMessage = errorPayload.message || t('routes_auth_google_toast_reg_failed_desc' as any);
            throw new Error(errorMessage);
          }

          if (!registerResultAction.payload.isOrganizationExists) {
            navigate($path('/organizations/new'));
            return;
          }

          const saveAuthResultAction = await dispatch(
            usersStore.asyncActions.saveAuth({
              authUrl: registerResultAction.payload.authUrl,
              data: { token: registerResultAction.payload.token },
            }),
          );

          if (!usersStore.asyncActions.saveAuth.fulfilled.match(saveAuthResultAction)) {
            const meResult = await dispatch(usersStore.asyncActions.getUsersMe());
            if (usersStore.asyncActions.getUsersMe.fulfilled.match(meResult)) {
              if (!meResult.payload.user.organization) {
                navigate($path('/organizations/new'));
                return;
              }
            }
            throw new Error(t('routes_auth_google_toast_save_auth_failed_desc' as any));
          }

          window.location.replace(saveAuthResultAction.payload.redirectUrl);
        } else {
          if (!authResult.isOrganizationExists) {
            navigate($path('/organizations/new'));
            return;
          }

          const saveAuthResultAction = await dispatch(
            usersStore.asyncActions.saveAuth({
              authUrl: authResult.authUrl,
              data: { token: authResult.token },
            }),
          );

          if (!usersStore.asyncActions.saveAuth.fulfilled.match(saveAuthResultAction)) {
            const meResult = await dispatch(usersStore.asyncActions.getUsersMe());
            if (usersStore.asyncActions.getUsersMe.fulfilled.match(meResult)) {
              if (!meResult.payload.user.organization) {
                navigate($path('/organizations/new'));
                return;
              }
            }
            throw new Error(t('routes_auth_google_toast_save_auth_failed_desc' as any));
          }

          window.location.replace(saveAuthResultAction.payload.redirectUrl);
        }
      } catch (error: any) {
        toast({
          title: t('routes_auth_google_toast_error_title' as any),
          description: error.message || t('routes_auth_google_toast_try_again_desc' as any),
          variant: 'destructive',
        });
        navigate($path('/login'));
      }
    };

    authenticateAndProceed();
  }, [dispatch, navigate, toast, t, code, state]);

  return (
    // Fixed escaped quotes
    <div className="m-auto flex items-center justify-center p-10">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
} 