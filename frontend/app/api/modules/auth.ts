import * as api from '@/api/instance';
import { ENVIRONMENTS } from '@/environments';

export const logout = async () => {
  const res = await api.post({
    path: '/auth/logout',
    auth: true,
  });
  return res;
};

export const refreshToken = async () => {
  const res = await api.post<{
    expiresAt: string;
  }>({
    path: '/auth/refresh',
    auth: true,
  });

  return res;
};

export const obtainAuthToken = async () => {
  const res = await api.post<{
    authUrl: 'string';
    token: 'string';
  }>({
    path: '/auth/token/obtain',
    auth: true,
  });

  return res;
};

export const saveAuth: (params: {
  authUrl: string;
  data: {
    token: string;
  };
}) => Promise<{
  expiresAt: string;
  redirectUrl: string;
}> = async (params) => {
  const res = await fetch(`${params.authUrl}`, {
    method: 'POST',
    credentials: 'include',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(params.data),
  });

  const json = await res.json();

  if (ENVIRONMENTS.MODE === 'development') {
    console.log('===============================');
    console.log({ authUrl: params.authUrl });
    console.log({ ...json });
    console.log('===============================');
  }

  if (res.ok) {
    return json;
  }

  if (json.errors) {
    throw json.errors;
  }

  throw new Error('Unknown error');
};

export const requestGoogleAuthLink = async () => {
  const res = await api.post<{
    authUrl: string;
  }>({
    path: '/auth/google/request',
  });

  return res;
};

export const authenticateWithGoogle = async (params: {
  data: {
    code: string;
    state: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    hasOrganization: boolean;
    hasMultipleOrganizations: boolean;
    isNewUser: boolean;
  }>({
    path: '/auth/google/authenticate',
    data: params.data,
  });

  return res;
};

export const registerWithGoogle = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    hasOrganization: boolean;
  }>({
    path: '/auth/google/register',
    data: params.data,
  });

  return res;
};

export const requestInvitationGoogleAuthLink = async (params: {
  data: {
    invitationToken: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
  }>({
    path: '/auth/invitations/google/request',
    data: params.data,
  });

  return res;
};

export const authenticateWithInvitationGoogleAuthLink = async (params: {
  data: {
    code: string;
    state: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isNewUser: boolean;
  }>({
    path: '/auth/invitations/google/authenticate',
    data: params.data,
  });

  return res;
};

export const registerWithInvitationGoogleAuthLink = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.post<{
    expiresAt: string;
  }>({
    path: '/auth/invitations/google/register',
    data: params.data,
  });

  return res;
};

export const requestMagicLink = async (params: {
  data: {
    email: string;
  };
}) => {
  const res = await api.post<{
    email: string;
    isNew: boolean;
  }>({
    path: '/auth/magic/request',
    data: params.data,
  });

  return res;
};

export const authenticateWithMagicLink = async (params: {
  data: {
    token: string;
    firstName?: string;
    lastName?: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    hasOrganization: boolean;
    isNewUser: boolean;
  }>({
    path: '/auth/magic/authenticate',
    data: params.data,
  });

  return res;
};

export const registerWithMagicLink = async (params: {
  data: {
    token: string;
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.post<{
    hasOrganization: boolean;
    expiresAt: string;
  }>({
    path: '/auth/magic/register',
    data: params.data,
  });

  return res;
};

export const requestInvitationMagicLink = async (params: {
  data: {
    invitationToken: string;
  };
}) => {
  const res = await api.post<{
    email: string;
  }>({
    path: '/auth/invitations/magic/request',
    data: params.data,
  });

  return res;
};

export const authenticateWithInvitationMagicLink = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isNewUser: boolean;
  }>({
    path: '/auth/invitations/magic/authenticate',
    data: params.data,
  });

  return res;
};

export const registerWithInvitationMagicLink = async (params: {
  data: {
    token: string;
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.post<{
    expiresAt: string;
  }>({
    path: '/auth/invitations/magic/register',
    data: params.data,
  });

  return res;
};
