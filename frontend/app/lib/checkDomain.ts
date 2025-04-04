import { ENVIRONMENTS, CLOUD_DOMAIN } from '@/environments';

export function checkDomain() {
  const returnValue: {
    isSourcetoolDomain: boolean;
    subDomain: string | null;
    environments: 'production' | 'staging' | 'local' | null;
  } = {
    isSourcetoolDomain: false,
    subDomain: null,
    environments: null,
  };

  const hostname = window.location.hostname;
  const isSourcetoolDomain = ENVIRONMENTS.IS_CLOUD_EDITION;

  returnValue.isSourcetoolDomain = isSourcetoolDomain;

  if (isSourcetoolDomain) {
    const subdomainRegex = new RegExp(
      `^(?:http[s]?:\\/\\/)?([^.]+)\\.${CLOUD_DOMAIN}`,
    );
    const matches = hostname.match(subdomainRegex);
    if (matches && matches[1]) {
      returnValue.subDomain = matches[1];
    }
    
    if (hostname.includes('staging')) {
      returnValue.environments = 'staging';
    } else if (hostname.includes('local')) {
      returnValue.environments = 'local';
    } else {
      returnValue.environments = 'production';
    }
  } else {
    if (hostname.includes('localhost')) {
      returnValue.environments = 'local';
    } else {
      returnValue.environments = 'production';
    }
  }

  return returnValue;
}
