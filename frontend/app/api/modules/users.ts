import * as api from '@/api/instance';

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

export const getMe = async () => {
  const res = await api.get<{
    user: User;
  }>({ path: '/users/me', auth: true });

  return res;
};

export const updateMe = async (params: {
  data: {
    firstName?: string;
    lastName?: string;
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: '/users/me',
    data: params.data,
    auth: true,
  });

  return res;
};

export const sendUpdateMeEmailInstructions = async (params: {
  data: {
    email: string;
    emailConfirmation: string;
  };
}) => {
  const res = await api.post({
    path: '/users/me/email/instructions',
    data: params.data,
    auth: true,
  });

  return res;
};

export const updateMeEmail = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: '/users/me/email',
    data: params.data,
    auth: true,
  });

  return res;
};

export const createUserInvitations = async (params: {
  data: {
    emails: string[];
    role: UserRole;
  };
}) => {
  const res = await api.post<{
    userInvitations: UserInvitation[];
  }>({
    path: '/users/invitations',
    data: params.data,
    auth: true,
  });

  return res;
};

export const resendUserInvitation = async (params: {
  invitationId: string;
}) => {
  const res = await api.post<{
    userInvitation: UserInvitation;
  }>({
    path: `/users/invitations/${params.invitationId}/resend`,
    auth: true,
  });

  return res;
};

export const listUsers = async () => {
  const res = await api.get<{
    userInvitations: UserInvitation[];
    users: User[];
  }>({ path: '/users', auth: true });

  return res;
};

export const updateUser = async (params: {
  userId: string;
  data: {
    role?: UserRole;
    groupIds?: string[];
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: `/users/${params.userId}`,
    data: params.data,
    auth: true,
  });

  return res;
};

export const deleteUser = async (params: {
  userId: string;
}) => {
  const res = await api.del<void>({
    path: `/users/${params.userId}`,
    auth: true,
  });

  return res;
};
