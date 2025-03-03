import { GoogleIcon } from '../icon/social-media-icon/google';
import { Button } from '../ui/button';

export function SocialButtonGoogle({
  onClick,
  label,
}: {
  onClick: () => void;
  label: string;
}) {
  return (
    <Button
      variant="outline"
      className="w-full"
      onClick={onClick}
      type="button"
    >
      <GoogleIcon className="size-4" />
      {label}
    </Button>
  );
}
