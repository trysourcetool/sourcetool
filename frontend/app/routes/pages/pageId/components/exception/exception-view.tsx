import { useSelector } from '@/store';
import { type FC } from 'react';

export const ExceptionView: FC = () => {
  const exception = useSelector((state) => state.pages.exception);
  return (
    exception && (
      <div className="flex flex-col gap-6">
        <div className="flex flex-col gap-2">
          {exception.title && (
            <h3 className="text-lg font-semibold text-foreground">
              {exception.title}
            </h3>
          )}
          {exception.message && (
            <p className="text-sm font-normal text-muted-foreground">
              {exception.message}
            </p>
          )}
        </div>
        <div className="overflow-hidden rounded-lg border bg-destructive">
          <div className="flex flex-col gap-1 bg-white/90 p-4">
            {exception.stackTrace?.map((stack) => (
              <p key={stack} className="text-sm font-normal text-destructive">
                {stack}
              </p>
            ))}
          </div>
        </div>
      </div>
    )
  );
};
