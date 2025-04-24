import { useToast } from '@/hooks/use-toast';
import { useDispatch } from '@/store';
import { usersStore } from '@/store/modules/users';
import { Loader2 } from 'lucide-react';
import { useEffect, useRef } from 'react';
import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { zodValidator } from '@tanstack/zod-adapter';
import { object, string } from 'zod';

export default function UsersEmailUpdateConfirm() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const search = useSearch({ from: '/_default/users/email/update/confirm/' });
  const token = search.token;

  const navigate = useNavigate();
  const { toast } = useToast();
  console.log({ token, search });

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
          navigate({ to: '/' });
        }
      } else {
        toast({
          title: 'Invalid token',
          description: 'Please try again',
          variant: 'destructive',
        });
        navigate({ to: '/' });
      }
    })();
  }, [token]);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}

export const Route = createFileRoute('/_default/users/email/update/confirm/')({
  component: UsersEmailUpdateConfirm,
  validateSearch: zodValidator(
    object({
      token: string(),
    }),
  ),
});
