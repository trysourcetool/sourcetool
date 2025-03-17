import { ENVIRONMENTS } from '@/environments';
import { checkSubDomain } from '@/lib/checkSubDomain';
import dayjs from 'dayjs';
import { usersRefreshToken } from './modules/users';
import { checkDomain } from '@/lib/checkDomain';

type SuccessResponse = {
  code: 0;
  message: string;
};

export type ErrorResponse = {
  detail: string;
  id: string;
  meta: {
    [key: string]: string;
  };
  status: number;
  title: string;
};

class Api {
  expiresAt: string | null;
  parseExpiresAt: dayjs.Dayjs | null;
  isRefreshWaiting: boolean;

  constructor() {
    this.expiresAt = null;
    this.parseExpiresAt = null;
    this.isRefreshWaiting = false;
  }

  setExpiresAt(expiresAt: string) {
    console.log({ expiresAt }, dayjs.unix(Number(expiresAt)));
    this.expiresAt = expiresAt;
    this.parseExpiresAt = dayjs.unix(Number(expiresAt));
  }

  async checkRefreshWaiting() {
    return new Promise((resolve) => {
      const timer = setInterval(() => {
        if (!this.isRefreshWaiting) {
          clearInterval(timer);
          resolve(true);
        }
      }, 100);
    });
  }

  async checkTokenExpiresAt() {
    if (this.isRefreshWaiting) {
      await this.checkRefreshWaiting();
    }
    if (
      this.expiresAt &&
      this.parseExpiresAt &&
      this.parseExpiresAt.diff(dayjs(new Date()), 'minute') < 1
    ) {
      console.log(
        this.expiresAt,
        this.parseExpiresAt,
        this.parseExpiresAt.diff(dayjs(new Date()), 'minute') < 1,
      );
      this.isRefreshWaiting = true;
      const xsrfToken = document.cookie
        .split('; ')
        .find((row) => row.startsWith('xsrf_token='));
      if (!xsrfToken) {
        return null;
      }
      try {
        const res = await usersRefreshToken();
        this.expiresAt = res.expiresAt;
        this.parseExpiresAt = dayjs(res.expiresAt);
        this.isRefreshWaiting = false;

        return null;
      } catch (error: any) {
        console.log({ error });
        return new Error();
      } finally {
        this.isRefreshWaiting = false;
      }
    }
    return null;
  }

  getParams(auth?: boolean) {
    const domain = checkDomain();
    console.log({ domain });
    const url = `${window.location.protocol}//${domain.isSourcetoolDomain && domain.subDomain ? `${domain.subDomain}.` : ''}${ENVIRONMENTS.API_BASE_URL}/api/v1`;

    const headers = new Headers();

    console.log('document.cookie', document.cookie);

    const xsrfToken = document.cookie
      .split('; ')
      .find((row) => row.startsWith('xsrf_token='));

    if (auth && xsrfToken) {
      headers.set('X-XSRF-TOKEN', xsrfToken.split('=')[1]);
    }

    return {
      url,
      headers,
    };
  }
}

export const api = new Api();

export const get: <T>(params: {
  path: string;
  params?: Record<string, any>;
  auth?: boolean;
}) => Promise<T> = async ({ params, path, auth }) => {
  const { url, headers } = api.getParams(auth);
  const urlParams = new URLSearchParams();
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        if (typeof value !== 'string') {
          value = JSON.stringify(value);
        }
        urlParams.set(key, value);
      }
    });
  }

  console.log({ url, headers });

  const res = await fetch(
    `${url}${path}${params ? '?' + urlParams.toString() : ''}`,
    {
      method: 'GET',
      credentials:
        ENVIRONMENTS.MODE === 'development' ? 'include' : 'same-origin',
      headers,
      mode: 'cors',
    },
  );

  const json = await res.json();

  if (ENVIRONMENTS.MODE === 'development') {
    console.log('===============================');
    console.log({ path });
    console.log({ ...json });
    console.log('===============================');
  }

  if (res.ok) {
    return json;
  }

  if (!res.ok) {
    throw json as ErrorResponse;
  }

  throw new Error('Unknown error');
};

export const post: <T = SuccessResponse>(params: {
  path: string;
  data?: Record<string, any>;
  auth?: boolean;
}) => Promise<T> = async ({ data, path, auth }) => {
  const { url, headers } = api.getParams(auth);
  const res = await fetch(`${url}${path}`, {
    method: 'POST',
    headers,
    credentials:
      ENVIRONMENTS.MODE === 'development' ? 'include' : 'same-origin',
    mode: 'cors',
    body: JSON.stringify(data),
  });

  const json = await res.json();

  if (ENVIRONMENTS.MODE === 'development') {
    console.log('===============================');
    console.log({ path });
    console.log({ ...json });
    console.log('===============================');
  }

  if (res.ok) {
    return json;
  }

  if (!res.ok) {
    throw json;
  }

  throw new Error('Unknown error');
};

export const put: <T = SuccessResponse>(params: {
  path: string;
  data?: object;
  auth?: boolean;
}) => Promise<T> = async ({ data, path, auth }) => {
  const { url, headers } = api.getParams(auth);

  const res = await fetch(`${url}${path}`, {
    method: 'PUT',
    headers,
    body: JSON.stringify(data),
    credentials:
      ENVIRONMENTS.MODE === 'development' ? 'include' : 'same-origin',
    mode: 'cors',
  });
  const json = await res.json();
  if (ENVIRONMENTS.MODE === 'development') {
    console.log('===============================');
    console.log({ path });
    console.log({ ...json });
    console.log('===============================');
  }

  if (res.ok) {
    return json;
  }

  if (!res.ok) {
    throw json;
  }

  throw new Error('Unknown error');
};

export const patch: <T = SuccessResponse>(params: {
  path: string;
  data?: Record<string, any>;
  auth?: boolean;
}) => Promise<T> = async ({ data, path, auth }) => {
  const { url, headers } = api.getParams(auth);

  console.log({ url, headers });

  const res = await fetch(`${url}${path}`, {
    method: 'PATCH',
    headers,
    body: JSON.stringify(data),
    mode: 'cors',
    credentials:
      ENVIRONMENTS.MODE === 'development' ? 'include' : 'same-origin',
  });
  const json = await res.json();
  if (ENVIRONMENTS.MODE === 'development') {
    console.log('===============================');
    console.log({ path });
    console.log({ ...json });
    console.log('===============================');
  }

  if (res.ok) {
    return json;
  }

  if (!res.ok) {
    throw json;
  }

  throw new Error('Unknown error');
};

export const del: <T = SuccessResponse>(params: {
  path: string;
  data?: Record<string, any>;
  auth?: boolean;
}) => Promise<T> = async ({ data, path, auth }) => {
  const { url, headers } = api.getParams(auth);
  const res = await fetch(`${url}${path}`, {
    method: 'DELETE',
    headers,
    body: JSON.stringify(data),
    credentials:
      ENVIRONMENTS.MODE === 'development' ? 'include' : 'same-origin',
    mode: 'cors',
  });
  const json = await res.json();
  if (ENVIRONMENTS.MODE === 'development') {
    console.log('===============================');
    console.log({ path });
    console.log({ ...json });
    console.log('===============================');
  }

  if (res.ok) {
    return json;
  }

  if (!res.ok) {
    throw json;
  }

  throw new Error('Unknown error');
};
