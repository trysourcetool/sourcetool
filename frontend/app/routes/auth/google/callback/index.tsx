import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import {
  useSearch,
  useNavigate,
  createFileRoute,
} from '@tanstack/react-router';
import { useDispatch } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';
import { zodValidator } from '@tanstack/zod-adapter';
import { object, string } from 'zod';

export default function GoogleAuthenticate() {
  const isProcessing = useRef(false);
  const search = useSearch({
    from: '/_default/auth/google/callback/',
  });
  const code = search.code;
  const state = search.state;
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
        navigate({ to: '/login' });
        return;
      }

      try {
        const authResultAction = await dispatch(
          authStore.asyncActions.authenticateWithGoogle({
            data: { code, state },
          }),
        );

        if (
          !authStore.asyncActions.authenticateWithGoogle.fulfilled.match(
            authResultAction,
          )
        ) {
          const errorPayload = (authResultAction.payload as any)?.error || {};
          const errorMessage =
            errorPayload.message ||
            t('routes_auth_google_toast_auth_failed_desc' as any);
          throw new Error(errorMessage);
        }

        const authResult = authResultAction.payload;

        if (authResult.hasMultipleOrganizations) {
          toast({
            title: t(
              'routes_auth_google_toast_multiple_organizations_title' as any,
            ),
            description: t(
              'routes_auth_google_toast_multiple_organizations_desc' as any,
            ),
          });
          navigate({ to: '/login' });
          return;
        }

        if (authResult.isNewUser) {
          const registerResultAction = await dispatch(
            authStore.asyncActions.registerWithGoogle({
              data: {
                token: authResult.token,
              },
            }),
          );

          if (
            !authStore.asyncActions.registerWithGoogle.fulfilled.match(
              registerResultAction,
            )
          ) {
            const errorPayload =
              (registerResultAction.payload as any)?.error || {};
            const errorMessage =
              errorPayload.message ||
              t('routes_auth_google_toast_reg_failed_desc' as any);
            throw new Error(errorMessage);
          }

          if (!registerResultAction.payload.hasOrganization) {
            navigate({ to: '/organizations/new' });
            return;
          }

          const saveAuthResultAction = await dispatch(
            authStore.asyncActions.saveAuth({
              authUrl: registerResultAction.payload.authUrl,
              data: { token: registerResultAction.payload.token },
            }),
          );

          if (
            !authStore.asyncActions.saveAuth.fulfilled.match(
              saveAuthResultAction,
            )
          ) {
            throw new Error(
              t('routes_auth_google_toast_save_auth_failed_desc' as any),
            );
          }

          window.location.replace(saveAuthResultAction.payload.redirectUrl);
        } else {
          if (!authResult.hasOrganization) {
            navigate({ to: '/organizations/new' });
            return;
          }

          const saveAuthResultAction = await dispatch(
            authStore.asyncActions.saveAuth({
              authUrl: authResult.authUrl,
              data: { token: authResult.token },
            }),
          );

          if (
            !authStore.asyncActions.saveAuth.fulfilled.match(
              saveAuthResultAction,
            )
          ) {
            throw new Error(
              t('routes_auth_google_toast_save_auth_failed_desc' as any),
            );
          }

          window.location.replace(saveAuthResultAction.payload.redirectUrl);
        }
      } catch (error: any) {
        toast({
          title: t('routes_auth_google_toast_error_title' as any),
          description:
            error.message ||
            t('routes_auth_google_toast_try_again_desc' as any),
          variant: 'destructive',
        });
        navigate({ to: '/login' });
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

export const Route = createFileRoute('/_default/auth/google/callback/')({
  component: GoogleAuthenticate,
  validateSearch: zodValidator(
    object({
      code: string(),
      state: string(),
    }),
  ),
});
