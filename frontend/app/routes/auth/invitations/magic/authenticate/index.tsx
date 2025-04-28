import { useEffect, useRef } from 'react';
import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { useDispatch } from '@/store';
import { authStore } from '@/store/modules/auth';
import { useToast } from '@/hooks/use-toast';
import { useTranslation } from 'react-i18next';
import { Loader2 } from 'lucide-react';
import { zodValidator } from '@tanstack/zod-adapter';
import { object, string } from 'zod';

export default function InvitationMagicLinkAuthenticate() {
  const dispatch = useDispatch();
  const { toast } = useToast();
  const navigate = useNavigate();
  const search = useSearch({
    from: '/_default/auth/invitations/magic/authenticate/',
  });
  const { t } = useTranslation('common');
  const isInitialLoading = useRef(false);

  const token = search.token;

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
        navigate({ to: '/login' });
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
          navigate({
            to: '/auth/invitations/signup/followup',
            search: { token: result.token },
          });
        } else {
          const saveAuthResult = await dispatch(
            authStore.asyncActions.saveAuth({
              authUrl: result.authUrl,
              data: { token: result.token },
            }),
          );

          if (
            !authStore.asyncActions.saveAuth.fulfilled.match(saveAuthResult)
          ) {
            throw new Error(
              t(
                'routes_auth_invitations_magic_link_toast_save_auth_failed_desc' as any,
              ),
            );
          }

          window.location.replace(saveAuthResult.payload.redirectUrl);
        }
      } else {
        toast({
          title: t('routes_login_toast_failed'),
          description: t('routes_login_toast_failed_description'),
          variant: 'destructive',
        });
        navigate({ to: '/login' });
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

export const Route = createFileRoute(
  '/_default/auth/invitations/magic/authenticate/',
)({
  component: InvitationMagicLinkAuthenticate,
  validateSearch: zodValidator(
    object({
      token: string(),
    }),
  ),
});
