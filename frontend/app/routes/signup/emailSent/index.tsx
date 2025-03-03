import { CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useSearchParams } from 'react-router';

export type SearchParams = {
  email: string;
};

export default function EmailSent() {
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email');

  return (
    <div className="m-auto flex items-center justify-center">
      <div className="flex max-w-[374px] flex-col gap-6 p-6">
        <CardHeader className="p-0">
          <CardTitle>Email sent!</CardTitle>
          <CardDescription>
            Please check the email we sent to {email} to create a new Sourcetool
            org
          </CardDescription>
        </CardHeader>
      </div>
    </div>
  );
}
