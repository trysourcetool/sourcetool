import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';
import { useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router';
import { $path } from 'safe-routes';

export type SearchParams = {
  token: string;
  email: string;
  isUserExists: string; // 'true' or 'false'
};

export default function InvitationActivate() {
  const isInitialLoading = useRef(false);
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const email = searchParams.get('email');
  const isUserExists = searchParams.get('isUserExists');
  const navigate = useNavigate();
  const { toast } = useToast();
  console.log({ token, email, isUserExists, searchParams });

  useEffect(() => {
    if (isInitialLoading.current) {
      return;
    }
    isInitialLoading.current = true;
    if (token) {
      if (isUserExists === 'false') {
        navigate(
          $path('/users/invitation/signup/followup', {
            token,
            email,
            isUserExists,
          }),
        );
      }
    } else {
      toast({
        title: 'Invalid token',
        description: 'Please try again',
        variant: 'destructive',
      });
      navigate('/signup');
    }
  }, [token, navigate, toast, isUserExists, email]);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
