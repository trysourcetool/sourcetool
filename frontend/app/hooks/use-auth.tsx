import { checkSubDomain } from '@/lib/checkSubDomain';
import { useDispatch, useSelector } from '@/store';
import { usersStore } from '@/store/modules/users';
import { Loader2 } from 'lucide-react';
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  type FC,
  type ReactNode,
} from 'react';

type AuthState = {
  subDomain: string | null;
  subDomainMatched: {
    isMatched: boolean;
    status: 'checking' | 'checked';
  };
};

export const authContext = createContext<AuthState>({
  subDomain: null,
  subDomainMatched: {
    isMatched: false,
    status: 'checking',
  },
});

export const AuthProvider: FC<{ children: ReactNode }> = (props) => {
  const isInitialAuthChecking = useRef(false);
  const dispatch = useDispatch();

  const subDomain = useMemo(
    () => (typeof window !== 'undefined' ? checkSubDomain() : null),
    [],
  );

  const isAuthChecked = useSelector((state) => state.users.isAuthChecked);
  const isAuthSucceeded = useSelector((state) => state.users.isAuthSucceeded);
  const isRefreshTokenWaiting = useSelector(
    (state) => state.users.isRefreshTokenWaiting,
  );

  const subDomainMatched = useSelector((state) =>
    usersStore.selector.getSubDomainMatched(state, subDomain),
  );

  useEffect(() => {
    if (
      typeof window !== 'undefined' &&
      !isInitialAuthChecking.current &&
      !isAuthChecked &&
      !isAuthSucceeded &&
      !isRefreshTokenWaiting
    ) {
      isInitialAuthChecking.current = true;
      (async () => {
        const resultAction = await dispatch(
          usersStore.asyncActions.refreshToken(),
        );
        if (
          usersStore.asyncActions.refreshToken.fulfilled.match(resultAction)
        ) {
          const resultUser = await dispatch(
            usersStore.asyncActions.getUsersMe(),
          );
          if (usersStore.asyncActions.getUsersMe.fulfilled.match(resultUser)) {
            const user = resultUser.payload.user;
            const userOrganization = user.organization;
            if (userOrganization && userOrganization?.subdomain === subDomain) {
              // TODO: redirect if outside authentication screen
            } else if (userOrganization && subDomain === 'auth') {
              window.location.replace(
                `${window.location.protocol}//${window.location.host.replace(
                  'auth',
                  userOrganization.subdomain,
                )}`,
              );
            }
          }
        }

        isInitialAuthChecking.current = false;
      })();
    }
  }, [
    dispatch,
    isAuthChecked,
    isRefreshTokenWaiting,
    subDomain,
    isAuthSucceeded,
  ]);

  return (
    <authContext.Provider
      value={{
        subDomain,
        subDomainMatched,
      }}
    >
      {isAuthChecked ? (
        props.children
      ) : (
        <div className="flex h-screen items-center justify-center">
          <Loader2 className="size-8 animate-spin" />
        </div>
      )}
    </authContext.Provider>
  );
};

export const useAuth = () => {
  const { subDomain, subDomainMatched } = useContext(authContext);

  const handleNoAuthRoute = useCallback(() => {
    if (!subDomain) {
      console.log(`auth.${window.location.host}`);
      window.location.replace(
        `${window.location.protocol}//auth.${window.location.host}/signin`,
      );
    } else if (subDomain && subDomain !== 'auth') {
      console.log(
        `${window.location.protocol}//${window.location.host}/signin`,
      );
      window.location.replace(
        `${window.location.protocol}//${window.location.host.replace(
          subDomain,
          'auth',
        )}/signin`,
      );
    }
  }, [subDomain]);

  console.log({
    subDomain,
    subDomainMatched,
  });

  return { subDomain, subDomainMatched, handleNoAuthRoute };
};
