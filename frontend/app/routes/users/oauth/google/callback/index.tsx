import { useToast } from '@/hooks/use-toast';
import { useDispatch } from '@/store';
import { usersStore } from '@/store/modules/users';
import { Loader2 } from 'lucide-react';
import { useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router';
import { $path } from 'safe-routes';

export type SearchParams = {
  isUserExists: string; // 'true' or 'false'
  firstName: string;
  lastName: string;
  token: string;
};

export default function OauthGoogleCallback() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const [searchParams] = useSearchParams();
  const isUserExists = searchParams.get('isUserExists');
  const firstName = searchParams.get('firstName');
  const lastName = searchParams.get('lastName');
  const token = searchParams.get('token');

  const navigate = useNavigate();
  const { toast } = useToast();
  console.log({ token, firstName, lastName, isUserExists });

  useEffect(() => {
    if (isInitialLoading.current) {
      return;
    }
    if (token) {
      isInitialLoading.current = true;
      (async () => {
        if (isUserExists === 'true') {
          const resultAction = await dispatch(
            usersStore.asyncActions.oauthGoogleSignin({
              data: {
                sessionToken: token,
              },
            }),
          );
          if (
            usersStore.asyncActions.oauthGoogleSignin.fulfilled.match(
              resultAction,
            )
          ) {
            if (resultAction.payload.isOrganizationExists) {
              const result = await dispatch(
                usersStore.asyncActions.saveAuth({
                  authUrl: resultAction.payload.authUrl,
                  data: {
                    token: resultAction.payload.token,
                  },
                }),
              );
              if (usersStore.asyncActions.saveAuth.fulfilled.match(result)) {
                window.location.replace(result.payload.redirectUrl);
                return;
              }
            } else {
              navigate($path('/organizations/new'));
              toast({
                title: 'Signin success',
                description: 'Next, create an organization',
              });
              return;
            }
          }
        } else {
          const resultAction = await dispatch(
            usersStore.asyncActions.oauthGoogleSignup({
              data: {
                sessionToken: token,
                firstName: firstName!,
                lastName: lastName!,
              },
            }),
          );
          if (
            usersStore.asyncActions.oauthGoogleSignup.fulfilled.match(
              resultAction,
            )
          ) {
            navigate($path('/organizations/new'));
            toast({
              title: 'Signup success',
              description: 'Next, create an organization',
            });
          }
        }
      })();
    } else {
      toast({
        title: 'Invalid token',
        description: 'Please try again',
        variant: 'destructive',
      });
      navigate('/signup');
      isInitialLoading.current = false;
    }
  }, []);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
