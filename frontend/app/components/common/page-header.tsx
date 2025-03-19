import { cn } from '@/lib/utils';
import type { FC } from 'react';

export const PageHeader: FC<{
  label: string;
  description?: string;
  border?: boolean;
}> = ({ label, description, border = true }) => {
  return (
    <div className={cn('flex flex-col gap-2', border && 'border-b p-4 md:p-6')}>
      <h1 className="text-3xl font-bold text-foreground">{label}</h1>
      {description && (
        <p className="text-base font-normal text-muted-foreground">
          {description}
        </p>
      )}
    </div>
  );
};
