import { useToast } from '@/hooks/use-toast';
import { useDispatch } from '@/store';
import { usersStore } from '@/store/modules/users';
import { Loader2 } from 'lucide-react';
import { useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router';
import { $path } from 'safe-routes';

export type SearchParams = {
  token: string;
};

export default function UsersEmailUpdateConfirm() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');

  const navigate = useNavigate();
  const { toast } = useToast();
  console.log({ token, searchParams });

  useEffect(() => {
    if (isInitialLoading.current) {
      return;
    }
    isInitialLoading.current = true;
    (async () => {
      if (token) {
        const resultAction = await dispatch(
          usersStore.asyncActions.updateMeEmail({
            data: {
              token,
            },
          }),
        );
        if (
          usersStore.asyncActions.updateMeEmail.fulfilled.match(resultAction)
        ) {
          toast({
            title: 'Email updated',
            description: 'Your email has been updated',
          });
          location.href = '/';
        } else {
          toast({
            title: 'Invalid token',
            description: 'Please try again',
            variant: 'destructive',
          });
          navigate($path('/'));
        }
      } else {
        toast({
          title: 'Invalid token',
          description: 'Please try again',
          variant: 'destructive',
        });
        navigate($path('/'));
      }
    })();
  }, [token]);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
