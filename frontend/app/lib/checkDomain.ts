import { ENVIRONMENTS } from '@/environments';

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
  const isSourcetoolDomain = ENVIRONMENTS.DOMAIN.match(
    /^((staging|local)\.)?trysourcetool\.com$/,
  );

  returnValue.isSourcetoolDomain = !!isSourcetoolDomain;

  if (isSourcetoolDomain) {
    const subdomainRegex = new RegExp(
      `^(?:http[s]?:\\/\\/)?([^.]+)\\.${ENVIRONMENTS.DOMAIN}`,
    );
    const matches = hostname.match(subdomainRegex);
    if (matches && matches[1]) {
      returnValue.subDomain = matches[1];
    }
    if (isSourcetoolDomain[2] === 'staging') {
      returnValue.environments = 'staging';
    } else if (isSourcetoolDomain[2] === 'local') {
      returnValue.environments = 'local';
    } else {
      returnValue.environments = 'production';
    }
  } else {
    if (ENVIRONMENTS.DOMAIN.includes('localhost')) {
      returnValue.environments = 'local';
    } else {
      returnValue.environments = 'production';
    }
  }

  return returnValue;
}
