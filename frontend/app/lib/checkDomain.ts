import { ENVIRONMENTS } from '@/environments';

export function checkDomain() {
  const returnValue: {
    subDomain: string | null;
    environments: 'production' | 'staging' | 'local';
  } = {
    subDomain: null,
    environments: ENVIRONMENTS.MODE === 'development' ? 'local' : 'production',
  };

  if (ENVIRONMENTS.IS_CLOUD_EDITION) {
    const hostname = window.location.hostname;
    const parts = hostname.split('.');
    if (parts.length > 2) {
      returnValue.subDomain = parts[0];
      if (parts[parts.length - 2] === 'staging') {
        returnValue.environments = 'staging';
      }
    }
  }

  return returnValue;
}
