import {
  index,
  layout,
  physical,
  rootRoute,
  route,
} from '@tanstack/virtual-file-routes';

export default rootRoute('root.tsx', [
  layout('default', 'layout-default.tsx', [
    physical('/login', 'login'),
    physical('/auth', 'auth'),
    physical('/signup', 'signup'),
    route(
      '/users/email/update/confirm/',
      'users/email/update/confirm/index.tsx',
    ),
    physical('/organizations/new', 'organizations/new'),
    physical('/onboarding', 'onboarding'),
    physical('/error', 'error'),
  ]),
  layout('auth', 'layout-auth-external.tsx', [
    index('pages/index.tsx'),
    route('/users', 'users/route.tsx', [
      route('/', 'users/index.tsx'),
      route('/$userId', 'users/$userId/index.tsx'),
    ]),
    physical('/apiKeys', 'apiKeys'),
    physical('/groups', 'groups'),
    physical('/environments', 'environments'),
    physical('/settings', 'settings'),
  ]),
  layout('preview', 'layout-auth-preview.tsx', [
    route('/pages/$', 'pages/pageId/index.tsx'),
  ]),
]);
