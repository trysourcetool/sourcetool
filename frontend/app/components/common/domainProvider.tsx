import { checkSubDomain } from '@/lib/checkSubDomain';
import { useSelector } from '@/store';
import { usersStore } from '@/store/modules/users';
import {
  createContext,
  useEffect,
  useMemo,
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
  const subDomain = useMemo(
    () => (typeof window !== 'undefined' ? checkSubDomain() : null),
    [],
  );

  const account = useSelector(usersStore.selector.getMe);

  const checkComplete = () => {
    setTimeout(() => {
      isChecking.current = false;
    }, 0);
  };

  useEffect(() => {
    if (isChecking.current) {
      return;
    }
    isChecking.current = true;

    if (subDomain && subDomain !== 'auth') {
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
      if (account) {
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
      }
    }
    if (subDomain && subDomain === 'auth') {
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
    }
    checkComplete();
  }, [pathname, account]);

  return (
    <domainContext.Provider value={{}}>{props.children}</domainContext.Provider>
  );
};
