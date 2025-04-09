import { expect, test, describe } from 'vitest';
import { createSourcetool, SourcetoolConfig } from './sourcetool';
import { newPageManager, Page } from './internal/page';
import { v4 as uuidv4 } from 'uuid';
const findPageByPath = (pages: Record<string, Page>, path: string): Page => {
  const page = Object.values(pages).find((p) => p.route === path);
  if (!page) {
    throw new Error(`Page not found: ${path}`);
  }
  return page;
};

const pageHandler = async (): Promise<void> => {};

describe('new', () => {
  const config: SourcetoolConfig = {
    apiKey: 'test_api_key',
    endpoint: 'ws://test.trysourcetool.com',
  };

  const sourcetool = createSourcetool(config);

  const tests = [
    { name: 'APIKey', got: sourcetool.apiKey, want: config.apiKey },
    {
      name: 'Endpoint',
      got: sourcetool.endpoint,
      want: 'ws://test.trysourcetool.com/ws',
    },
    {
      name: 'Pages length',
      got: Object.keys(sourcetool.pages).length,
      want: 0,
    },
  ];

  for (const t of tests) {
    test(t.name, () => {
      expect(t.got).toEqual(t.want);
    });
  }
});

describe('page', () => {
  test('Public page', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.page('/public', 'Public Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/public');
    expect(page.route).toEqual('/public');
    expect(page.accessGroups.length).toEqual(0);
  });

  test('Page with direct access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.accessGroups('admin');
    sourcetool.page('/admin', 'Admin Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/admin');
    expect(page.route).toEqual('/admin');
    expect(page.accessGroups[0]).toEqual('admin');
  });

  test('Group with access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const api = sourcetool.group('/api');
    api.accessGroups('api_user');
    api.page('/users', 'Users API', pageHandler);
    api.page('/posts', 'Posts API', pageHandler);

    const usersPage = findPageByPath(sourcetool.pages, '/api/users');
    const postsPage = findPageByPath(sourcetool.pages, '/api/posts');

    expect(usersPage.route).toEqual('/api/users');
    expect(usersPage.accessGroups[0]).toEqual('api_user');
    expect(postsPage.route).toEqual('/api/posts');
    expect(postsPage.accessGroups[0]).toEqual('api_user');
  });

  test('Nested groups with access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const users = sourcetool.group('/users');
    users.accessGroups('admin');
    users.page('/list', 'List users page', pageHandler);
    users
      .accessGroups('customer_support')
      .page('/create', 'Create user page', pageHandler);

    const products = users.group('/products');
    products.accessGroups('product_manager');
    products.page('/list', 'List products page', pageHandler);

    const tests = [
      { path: '/users/list', expectedGroups: ['admin'] },
      { path: '/users/create', expectedGroups: ['admin', 'customer_support'] },
      {
        path: '/users/products/list',
        expectedGroups: ['admin', 'customer_support', 'product_manager'],
      },
    ];

    for (const t of tests) {
      const page = findPageByPath(sourcetool.pages, t.path);
      expect(page.accessGroups.every((g) => t.expectedGroups.includes(g))).toBe(
        true,
      );
    }
  });

  test('Complex group structure', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);

    const admin = sourcetool.group('/admin');
    admin.accessGroups('admin');
    admin.page('/dashboard', 'Admin Dashboard', pageHandler);

    const settings = admin.group('/settings');
    settings.accessGroups('super_admin');
    settings.page('/system', 'System Settings', pageHandler);

    const api = sourcetool.group('/api');
    api.accessGroups('api_user');

    const v1 = api.group('/v1');
    v1.page('/users', 'Users API v1', pageHandler);

    const v2 = api.group('/v2');
    v2.accessGroups('api_v2');
    v2.page('/users', 'Users API v2', pageHandler);

    const tests = [
      { path: '/admin/dashboard', expectedGroups: ['admin'] },
      {
        path: '/admin/settings/system',
        expectedGroups: ['admin', 'super_admin'],
      },
      { path: '/api/v1/users', expectedGroups: ['api_user'] },
      { path: '/api/v2/users', expectedGroups: ['api_user', 'api_v2'] },
    ];

    for (const t of tests) {
      const page = findPageByPath(sourcetool.pages, t.path);
      expect(page.accessGroups.every((g) => t.expectedGroups.includes(g))).toBe(
        true,
      );
    }
  });

  test('Error handling', async () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const errorHandler = async () => {
      throw new Error('test error');
    };
    sourcetool.page('/error', 'Error Page', errorHandler);

    const page = findPageByPath(sourcetool.pages, '/error');

    try {
      await (page.run as any)();
    } catch (error) {
      expect(error).toBeInstanceOf(Error);
      expect((error as Error).message).toEqual('test error');
    }
  });
});

describe('page manager', () => {
  test('Get existing page', () => {
    const pages: Record<string, Page> = {};
    const pageId = uuidv4();
    const testPage = new Page(pageId, 'TestPage');
    pages[pageId] = testPage;

    const pageManager = newPageManager(pages);
    const got = pageManager.getPage(pageId);
    expect(got?.id).toEqual(pageId);
  });

  test('Get non-existent page', () => {
    const pages: Record<string, Page> = {};
    const pageManager = newPageManager(pages);
    const nonExistentId = uuidv4();
    const got = pageManager.getPage(nonExistentId);
    expect(got).toBeUndefined();
  });
});
