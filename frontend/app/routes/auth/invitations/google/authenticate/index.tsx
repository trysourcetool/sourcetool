import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams, useNavigate } from 'react-router'; // Use react-router
import { useDispatch } from '@/store';
import { usersStore } from '@/store/modules/users';
import { useToast } from '@/hooks/use-toast';
import { $path } from 'safe-routes';
import { Loader2 } from 'lucide-react';

export default function InvitationGoogleAuthenticate() {
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
          usersStore.asyncActions.authenticateWithInvitationGoogleAuthLink({
            data: { code, state },
          }),
        );

        if (
          !usersStore.asyncActions.authenticateWithInvitationGoogleAuthLink.fulfilled.match(
            authResultAction,
          )
        ) {
          // Attempt to extract error message from backend if available
          const errorPayload = (authResultAction.payload as any)?.error || {};
          const errorMessage = errorPayload.message || t('routes_auth_google_toast_auth_failed_desc' as any);
           toast({
             title: t('routes_auth_google_toast_auth_failed_title' as any),
             description: errorMessage,
             variant: 'destructive',
           });
           navigate($path('/login'));
           return; // Stop execution after handling error
        }

        const authResult = authResultAction.payload;

        if (authResult.isNewUser) {
          // New user via invitation: Register directly
          const registerResultAction = await dispatch(
            usersStore.asyncActions.registerWithInvitationGoogleAuthLink({
              data: {
                token: authResult.token, // Registration token from auth step
                // Use names from Google response, provide fallbacks if necessary
                firstName: authResult.firstName || '-', // Provide a fallback
                lastName: authResult.lastName || '-',  // Provide a fallback
              },
            }),
          );

          if (
            !usersStore.asyncActions.registerWithInvitationGoogleAuthLink.fulfilled.match(
              registerResultAction,
            )
          ) {
            const errorPayload = (registerResultAction.payload as any)?.error || {};
            const errorMessage = errorPayload.message || t('routes_auth_google_toast_reg_failed_desc' as any);
            throw new Error(errorMessage);
          }

          // Registration successful, fetch user data and redirect to the app
          const meResultAction = await dispatch(usersStore.asyncActions.getUsersMe());
          if (!usersStore.asyncActions.getUsersMe.fulfilled.match(meResultAction)) {
            // Handle case where fetching user fails after registration
            console.error("Failed to fetch user details after successful registration.");
            toast({
               title: t('routes_auth_google_toast_reg_success_title' as any),
               description: t('routes_auth_google_toast_getme_failed_desc' as any),
               variant: 'default', // Use default variant for informational message
             });
             // Redirect to login, as we can't confirm the user state
             navigate($path('/login')); 
             return; 
          }
          
           toast({
             title: t('routes_auth_google_toast_reg_success_title' as any),
             description: t('routes_auth_google_toast_reg_success_desc' as any),
           });
           // Since it's an invitation, org should exist. Redirect to main app.
           window.location.replace($path('/'));

        } else {
           // Existing user successfully authenticated with invitation link
           if (!authResult.isOrganizationExists) {
             // This case might indicate an issue, as invitation implies an org.
             // Redirect to create org, or handle as error? Following existing patterns:
             console.warn("Authenticated via invitation but organization doesn't exist?");
             navigate($path('/organizations/new'));
             return;
           }

          // Ensure authUrl exists before proceeding for existing user flow
          if (!authResult.authUrl) {
            console.error("AuthURL missing in response for existing user.");
            throw new Error(t('routes_auth_google_toast_auth_failed_desc' as any)); // Or a more specific error
          }

          const saveAuthResultAction = await dispatch(
            usersStore.asyncActions.saveAuth({
              authUrl: authResult.authUrl, // Now guaranteed to be a string
              data: { token: authResult.token }, // Temporary token from auth step
            }),
          );

          if (
            !usersStore.asyncActions.saveAuth.fulfilled.match(
              saveAuthResultAction,
            )
          ) {
             // saveAuth failed, try to get user state
             const meResult = await dispatch(usersStore.asyncActions.getUsersMe());
             if (usersStore.asyncActions.getUsersMe.fulfilled.match(meResult)) {
               if (!meResult.payload.user.organization) {
                 navigate($path('/organizations/new'));
                 return;
               }
             }
             // If saveAuth failed and getMe also failed or user still has no org, show error
             throw new Error(t('routes_auth_google_toast_save_auth_failed_desc' as any));
          }

          // saveAuth successful, redirect to the application
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
  }, [dispatch, navigate, toast, t, code, state]); // Add code and state to dependency array

  return (
    <div className="m-auto flex items-center justify-center p-10">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
} 