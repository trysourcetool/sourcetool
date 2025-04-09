import {
  type RouteConfig,
  layout,
  index,
  route,
} from '@react-router/dev/routes';

export default [
  layout('./routes/layout-default.tsx', [
    route('/login', './routes/login/index.tsx'),
    route('/login/emailSent', './routes/login/emailSent/index.tsx'),
    route(
      '/auth/invitations/login',
      './routes/auth/invitations/login/index.tsx',
    ),
    route(
      '/auth/invitations/emailSent',
      './routes/auth/invitations/emailSent/index.tsx',
    ),
    route(
      '/auth/invitations/magic/authenticate',
      './routes/auth/invitations/magic/authenticate/index.tsx',
    ),
    route(
      '/auth/invitations/signup/followup',
      './routes/auth/invitations/signup/followup/index.tsx',
    ),
    route(
      '/auth/magic/authenticate',
      './routes/auth/magic/authenticate/index.tsx',
    ),
    route(
      '/auth/google/callback',
      './routes/auth/google/callback/index.tsx',
    ),
    route('/signup/followup', './routes/signup/followup/index.tsx'),
    route(
      '/users/email/update/confirm',
      './routes/users/email/update/confirm/index.tsx',
    ),
    route('/organizations/new', './routes/organizations/new/index.tsx'),
    route('/onboarding', './routes/onboarding/index.tsx'),
    route('/onboarding/complete', './routes/onboarding/complete/index.tsx'),
    route(
      '/error/hostInstancePingError',
      './routes/error/hostInstancePingError/index.tsx',
    ),
  ]),
  layout('./routes/layout-auth-external.tsx', [
    index('./routes/pages/index.tsx'),
    route('/users', './routes/users/index.tsx', [
      route('/users/invite', './routes/users/invite/index.tsx'),
    ]),
    route('/users/:userId', './routes/users/userId/index.tsx'),
    route('/apiKeys', './routes/apiKeys/index.tsx'),
    route('/apiKeys/new', './routes/apiKeys/new/index.tsx'),
    route('/apiKeys/:apiKeyId', './routes/apiKeys/apiKeyId/index.tsx'),
    route('/groups', './routes/groups/index.tsx'),
    route('/groups/new', './routes/groups/new/index.tsx'),
    route('/groups/:groupId', './routes/groups/groupId/index.tsx'),
    route('/environments', './routes/environments/index.tsx'),
    route('/environments/new', './routes/environments/new/index.tsx'),
    route(
      '/environments/:environmentId',
      './routes/environments/environmentId/index.tsx',
    ),
    route('/settings', './routes/settings/index.tsx'),
  ]),
  layout('./routes/layout-auth-preview.tsx', [
    route('/pages/*', './routes/pages/pageId/index.tsx'),
  ]),
] satisfies RouteConfig;
