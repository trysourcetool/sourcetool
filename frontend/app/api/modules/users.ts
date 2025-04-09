import * as api from '@/api/instance';
import { ENVIRONMENTS } from '@/environments';

export type UserAuthResponse = {
  authUrl: string;
  isOrganizationExists: boolean;
  token: string;
};

export type UserInvitation = {
  createdAt: string;
  email: string;
  id: string;
};

export type UserRole = 'admin' | 'member' | 'developer';

export type User = {
  createdAt: string;
  email: string;
  firstName: string;
  id: string;
  lastName: string;
  role: UserRole;
  updatedAt: string;
  organization: {
    createdAt: string;
    id: string;
    subdomain: string;
    updatedAt: string;
    webSocketEndpoint: string;
  } | null;
};

export const listUsers = async () => {
  const res = await api.get<{
    userInvitations: UserInvitation[];
    users: User[];
  }>({ path: '/users', auth: true });

  return res;
};

export const getUsersMe = async () => {
  const res = await api.get<{
    user: User;
  }>({ path: '/users/me', auth: true });

  return res;
};

export const usersSignout = async () => {
  const res = await api.post({
    path: '/users/signout',
    auth: true,
  });

  return res;
};

export const usersSignup = async (params: {
  data: {
    firstName: string;
    lastName: string;
    password: string;
    passwordConfirmation: string;
    token: string;
  };
}) => {
  const res = await api.post({
    path: '/users/signup',
    data: params.data,
  });

  return res;
};

export const usersSignupInstructions = async (params: {
  data: {
    email: string;
  };
}) => {
  const res = await api.post<{
    email: string;
  }>({
    path: '/users/signup/instructions',
    data: params.data,
  });

  return res;
};

export const usersRefreshToken = async () => {
  const res = await api.post<{
    expiresAt: string;
  }>({
    path: '/users/refreshToken',
    auth: true,
  });

  return res;
};

export const usersObtainAuthToken = async () => {
  const res = await api.post<{
    authUrl: 'string';
    token: 'string';
  }>({
    path: '/users/obtainAuthToken',
    auth: true,
  });

  return res;
};

export const usersSaveAuth: (params: {
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

export const usersRequestGoogleAuthLink = async () => {
  const res = await api.post<{
    authUrl: string;
  }>({
    path: '/users/auth/google/request',
  });

  return res;
};

export const usersAuthenticateWithGoogle = async (params: {
  data: {
    code: string;
    state: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isOrganizationExists: boolean;
    isNewUser: boolean;
    firstName: string;
    lastName: string;
  }>({
    path: '/users/auth/google/authenticate',
    data: params.data,
  });

  return res;
};

export const usersRegisterWithGoogle = async (params: {
  data: {
    token: string;
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isOrganizationExists: boolean;
  }>({
    path: '/users/auth/google/register',
    data: params.data,
  });

  return res;
};

export const usersInvite = async (params: {
  data: {
    emails: string[];
    role: UserRole;
  };
}) => {
  const res = await api.post<{
    userInvitations: UserInvitation[];
  }>({
    path: '/users/invite',
    data: params.data,
    auth: true,
  });

  return res;
};

export const usersInvitationsResend = async (params: {
  data: {
    invitationId: string;
  };
}) => {
  const res = await api.post<{
    userInvitation: UserInvitation;
  }>({
    path: '/users/invitations/resend',
    data: params.data,
    auth: true,
  });

  return res;
};

export const usersRequestInvitationGoogleAuthLink = async (params: {
  data: {
    invitationToken: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
  }>({
    path: '/users/auth/invitations/google/request',
    data: params.data,
  });

  return res;
};

export const usersAuthenticateWithInvitationGoogleAuthLink = async (params: {
  data: {
    code: string;
    state: string;
  };
}) => {
  const res = await api.post<{
    authUrl?: string;
    token: string;
    isOrganizationExists: boolean;
    isNewUser: boolean;
    firstName?: string;
    lastName?: string;
  }>({
    path: '/users/auth/invitations/google/authenticate',
    data: params.data,
  });

  return res;
};

export const usersRegisterWithInvitationGoogleAuthLink = async (params: {
  data: {
    token: string;
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.post({
    path: '/users/auth/invitations/google/register',
    data: params.data,
  });

  return res;
};

export const updateUser = async (params: {
  data: {
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: '/users',
    data: params.data,
    auth: true,
  });

  return res;
};

export const updateUserEmail = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: '/users/email',
    data: params.data,
    auth: true,
  });

  return res;
};

export const usersSendUpdateEmailInstructions = async (params: {
  data: {
    email: string;
    emailConfirmation: string;
  };
}) => {
  const res = await api.post({
    path: '/users/sendUpdateEmailInstructions',
    data: params.data,
    auth: true,
  });

  return res;
};

export const usersRequestMagicLink = async (params: {
  data: {
    email: string;
  };
}) => {
  const res = await api.post({
    path: '/users/auth/magic/request',
    data: params.data,
  });

  return res;
};

export const usersAuthenticateWithMagicLink = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isNewUser: boolean;
    isOrganizationExists: boolean;
  }>({
    path: '/users/auth/magic/authenticate',
    data: params.data,
  });

  return res;
};

export const usersRegisterWithMagicLink = async (params: {
  data: {
    token: string;
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.post<{
    expiresAt: string;
  }>({
    path: '/users/auth/magic/register',
    data: params.data,
  });

  return res;
};

export const usersRequestInvitationMagicLink = async (params: {
  data: {
    invitationToken: string;
  };
}) => {
  const res = await api.post<{
    email: string;
  }>({
    path: '/users/auth/invitations/magic/request',
    data: params.data,
  });

  return res;
};

export const usersAuthenticateWithInvitationMagicLink = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isNewUser: boolean;
  }>({
    path: '/users/auth/invitations/magic/authenticate',
    data: params.data,
  });

  return res;
};

export const usersRegisterWithInvitationMagicLink = async (params: {
  data: {
    token: string;
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.post<{
    expiresAt: string;
  }>({
    path: '/users/auth/invitations/magic/register',
    data: params.data,
  });

  return res;
};
