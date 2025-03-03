import { ENVIRONMENTS } from '@/environments';

export function checkSubDomain() {
  const hostname = window.location.hostname;
  const subdomainRegex = new RegExp(
    `^(?:http[s]?:\\/\\/)?([^.]+)\\.${ENVIRONMENTS.DOMAIN}`,
  );

  const matches = hostname.match(subdomainRegex);
  console.log({ matches });

  if (matches && matches[1]) {
    return matches[1];
  }

  return null;
}
