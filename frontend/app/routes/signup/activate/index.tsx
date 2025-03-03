import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';
import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router';
import { $path } from 'safe-routes';

export type SearchParams = {
  token: string;
};

export default function SignupActivate() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const navigate = useNavigate();
  const { toast } = useToast();
  console.log({ token, searchParams });

  useEffect(() => {
    console.log('token', token);
    if (token) {
      navigate($path('/signup/followup', { token }));
    } else {
      toast({
        title: 'Invalid token',
        description: 'Please try again',
        variant: 'destructive',
      });
      navigate('/signup');
    }
  }, [token, navigate, toast]);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}
