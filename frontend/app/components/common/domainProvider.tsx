import { useAuth } from '@/hooks/use-auth';
import { checkDomain } from '@/lib/checkDomain';
import { useSelector } from '@/store';
import { usersStore } from '@/store/modules/users';
import {
  createContext,
  useEffect,
  useRef,
  type FC,
  type ReactNode,
} from 'react';
import { useLocation, useNavigate } from 'react-router';
import { $path } from 'safe-routes';

export const domainContext = createContext({});

export const DomainProvider: FC<{ children: ReactNode }> = (props) => {
  const isChecking = useRef<boolean>(false);
  const pathname = useLocation().pathname;
  const navigate = useNavigate();
  if (typeof window !== 'undefined') {
    checkDomain();
  }

  const account = useSelector(usersStore.selector.getMe);
  const { isAuthChecked, isSourcetoolDomain, subDomain } = useAuth();

  const checkComplete = () => {
    setTimeout(() => {
      isChecking.current = false;
    }, 0);
  };

  const handleAuthorizedRoute = () => {
    if (pathname.startsWith('/users/email/update/confirm')) {
      return;
    }
    if (
      pathname.startsWith('/signup') ||
      pathname.startsWith('/organizations/new')
    ) {
      navigate($path('/signin'));
      checkComplete();
      return;
    }
    if (account?.organization) {
      if (
        pathname.startsWith('/signin') ||
        pathname.startsWith('/users/invitation') ||
        pathname.startsWith('/resetPassword')
      ) {
        navigate($path('/'));
        checkComplete();
        return;
      }
      if (
        account?.role === 'member' &&
        ((pathname.startsWith('/users') &&
          !pathname.startsWith('/users/oauth/google')) ||
          pathname.startsWith('/apiKeys') ||
          pathname.startsWith('/environments'))
      ) {
        navigate($path('/'));
        checkComplete();
        return;
      }
    } else if (account && !account.organization) {
      navigate('/organizations/new');
      checkComplete();
      return;
    }
  };

  const handleUnauthorizedRoute = () => {
    if (
      pathname === '/' ||
      pathname.startsWith('/pages') ||
      (pathname.startsWith('/users') &&
        !pathname.startsWith('/users/oauth/google')) ||
      pathname.startsWith('/apiKeys') ||
      pathname.startsWith('/environments') ||
      pathname.startsWith('/onboarding')
    ) {
      navigate($path('/signin'));
      checkComplete();
      return;
    }
  };

  useEffect(() => {
    if (isChecking.current || isAuthChecked !== 'checked') {
      return;
    }
    isChecking.current = true;

    if (isSourcetoolDomain) {
      if (subDomain && subDomain !== 'auth') {
        handleAuthorizedRoute();
      }
      if (subDomain && subDomain === 'auth') {
        handleUnauthorizedRoute();
      }
    } else {
      if (account) {
        handleAuthorizedRoute();
      } else {
        handleUnauthorizedRoute();
      }
    }
    checkComplete();
  }, [pathname, account, isAuthChecked]);

  return (
    <domainContext.Provider value={{}}>{props.children}</domainContext.Provider>
  );
};
