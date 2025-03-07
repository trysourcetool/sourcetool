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

export const usersSignin = async (params: {
  data: {
    email: string;
    password: string;
  };
}) => {
  const res = await api.post<UserAuthResponse>({
    path: '/users/signin',
    data: params.data,
  });

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

export const usersOauthGoogleUserInfo = async (params: {
  data: {
    code: string;
    state: string;
  };
}) => {
  const res = await api.post<{
    firstName: string;
    isUserExists: true;
    lastName: string;
    sessionToken: string;
  }>({
    path: '/users/oauth/google/userInfo',
    data: params.data,
  });

  return res;
};

export const usersOauthGoogleSignin = async (params: {
  data: {
    sessionToken: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isOrganizationExists: boolean;
  }>({
    path: '/users/oauth/google/signin',
    data: params.data,
  });

  return res;
};

export const usersOauthGoogleSignup = async (params: {
  data: {
    firstName: string;
    lastName: string;
    sessionToken: string;
  };
}) => {
  const res = await api.post({
    path: '/users/oauth/google/signup',
    data: params.data,
  });

  return res;
};

export const usersOauthGoogleAuthCodeUrl = async () => {
  const res = await api.post<{
    url: string;
  }>({
    path: '/users/oauth/google/authCodeUrl',
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

export const usersInvitationsSignup = async (params: {
  data: {
    firstName: string;
    invitationToken: string;
    lastName: string;
    password: string;
    passwordConfirmation: string;
  };
}) => {
  const res = await api.post({
    path: '/users/invitations/signup',
    data: params.data,
  });

  return res;
};

export const usersInvitationsSignin = async (params: {
  data: {
    invitationToken: string;
    password: string;
  };
}) => {
  const res = await api.post({
    path: '/users/invitations/signin',
    data: params.data,
  });

  return res;
};

export const usersInvitationsOauthGoogleUserInfo = async (params: {
  data: {
    code: string;
    state: string;
  };
}) => {
  const res = await api.post<{
    firstName: string;
    isUserExists: true;
    lastName: string;
    sessionToken: string;
  }>({
    path: '/users/invitations/oauth/google/userInfo',
    data: params.data,
  });

  return res;
};

export const usersInvitationsOauthGoogleSignup = async (params: {
  data: {
    firstName: string;
    invitationToken: string;
    lastName: string;
    sessionToken: string;
  };
}) => {
  const res = await api.post({
    path: '/users/invitations/oauth/google/signup',
    data: params.data,
  });

  return res;
};

export const usersInvitationsOauthGoogleSignin = async (params: {
  data: {
    invitationToken: string;
    sessionToken: string;
  };
}) => {
  const res = await api.post({
    path: '/users/invitations/oauth/google/signin',
    data: params.data,
  });

  return res;
};

export const usersInvitationsOauthGoogleAuthCodeUrl = async (params: {
  data: {
    invitationToken: string;
  };
}) => {
  const res = await api.post<{
    url: string;
  }>({
    path: '/users/invitations/oauth/google/authCodeUrl',
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
  }>({ path: '/users', data: params.data, auth: true });

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

export const updateUserPassword = async (params: {
  data: {
    currentPassword: string;
    password: string;
    passwordConfirmation: string;
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: '/users/password',
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
