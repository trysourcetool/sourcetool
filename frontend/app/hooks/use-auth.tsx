import { checkDomain } from '@/lib/checkDomain';
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
  isSourcetoolDomain: boolean;
  isSubDomainMatched: boolean;
  isAuthChecked: 'checking' | 'checked';
  environments: 'production' | 'staging' | 'local' | null;
};

export const authContext = createContext<AuthState>({
  subDomain: null,
  isSourcetoolDomain: false,
  isSubDomainMatched: false,
  isAuthChecked: 'checking',
  environments: null,
});

export const AuthProvider: FC<{ children: ReactNode }> = (props) => {
  const isInitialAuthChecking = useRef(false);
  const dispatch = useDispatch();

  const domain = useMemo(
    () => (typeof window !== 'undefined' ? checkDomain() : null),
    [],
  );

  const isAuthChecked = useSelector((state) => state.users.isAuthChecked);
  const isAuthSucceeded = useSelector((state) => state.users.isAuthSucceeded);
  const isRefreshTokenWaiting = useSelector(
    (state) => state.users.isRefreshTokenWaiting,
  );

  const subDomainMatched = useSelector((state) =>
    usersStore.selector.getSubDomainMatched(state, domain?.subDomain ?? null),
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
            if (
              userOrganization &&
              userOrganization?.subdomain === domain?.subDomain
            ) {
              // TODO: Redirect if outside the authentication screen
            } else if (userOrganization && domain?.subDomain === 'auth') {
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
    domain?.subDomain,
    isAuthSucceeded,
  ]);

  return (
    <authContext.Provider
      value={{
        subDomain: domain?.subDomain ?? null,
        isSourcetoolDomain: !!domain?.isSourcetoolDomain,
        isSubDomainMatched: subDomainMatched.isMatched,
        isAuthChecked: subDomainMatched.status,
        environments: domain?.environments ?? null,
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
  const {
    subDomain,
    isSubDomainMatched,
    isAuthChecked,
    isSourcetoolDomain,
    environments,
  } = useContext(authContext);

  const handleNoAuthRoute = useCallback(() => {
    if (isSourcetoolDomain && !subDomain) {
      console.log(`auth.${window.location.host}`);
      window.location.replace(
        `${window.location.protocol}//auth.${window.location.host}/signin`,
      );
    } else if (isSourcetoolDomain && subDomain && subDomain !== 'auth') {
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

  return {
    subDomain,
    isSubDomainMatched,
    isAuthChecked,
    handleNoAuthRoute,
    isSourcetoolDomain,
    environments,
  };
};
