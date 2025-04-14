import { expect, test, describe } from 'vitest';
import { removeDuplicates, Router } from '../router';
import { createSourcetool, SourcetoolConfig } from '../sourcetool';
import { Page } from '../page';

const findPageByPath = (
  pages: Record<string, Page>,
  path: string,
): Page | null => {
  const page = Object.values(pages).find((p) => p.route === path);
  if (!page) {
    return null;
  }
  return page;
};

describe('join path', () => {
  const tests = [
    {
      name: 'Empty base path',
      basePath: '',
      path: '/users',
      want: '/users',
    },
    {
      name: 'Base path with trailing slash',
      basePath: '/admin/',
      path: 'users',
      want: '/admin/users',
    },
    {
      name: 'Path without leading slash',
      basePath: '/admin',
      path: 'users',
      want: '/admin/users',
    },
    {
      name: 'Both with slashes',
      basePath: '/admin/',
      path: '/users/',
      want: '/admin/users',
    },
    {
      name: 'Nested paths',
      basePath: '/api/v1',
      path: 'users/list',
      want: '/api/v1/users/list',
    },
    {
      name: 'Root path',
      basePath: '',
      path: '/',
      want: '/',
    },
    {
      name: 'Root path with base path',
      basePath: '/admin',
      path: '/',
      want: '/admin',
    },
  ];

  for (const t of tests) {
    test(t.name, () => {
      const router = new Router(
        null,
        'test.trysourcetool.com',
        null,
        t.basePath,
      );
      const result = router.joinPath(t.path);
      expect(result).toBe(t.want);
    });
  }
});

describe('remove duplicates', () => {
  const tests = [
    {
      name: 'No duplicates',
      groups: ['admin', 'user', 'guest'],
      want: ['admin', 'user', 'guest'],
    },
    {
      name: 'With duplicates',
      groups: ['admin', 'user', 'admin', 'guest', 'user'],
      want: ['admin', 'user', 'guest'],
    },
    {
      name: 'Empty slice',
      groups: [],
      want: [],
    },
    {
      name: 'Single element',
      groups: ['admin'],
      want: ['admin'],
    },
    {
      name: 'All duplicates',
      groups: ['admin', 'admin', 'admin'],
      want: ['admin'],
    },
  ];

  for (const t of tests) {
    test(t.name, () => {
      const result = removeDuplicates(t.groups);
      expect(result).toEqual(t.want);
    });
  }
});

describe('generate page id', () => {
  const router = new Router(null, 'test.trysourcetool.com', null);

  const tests = [
    {
      name: 'Simple path',
      path: '/users',
      wantSame: true,
    },
    {
      name: 'Nested path',
      path: '/admin/users/list',
      wantSame: true,
    },
    {
      name: 'Root path',
      path: '/',
      wantSame: true,
    },
  ];

  for (const t of tests) {
    test(t.name, () => {
      const id1 = router.generatePageID(t.path);
      const id2 = router.generatePageID(t.path);
      expect(id1).toBe(id2);

      const differentPath = t.path + '/different';
      const id3 = router.generatePageID(differentPath);
      expect(id3).not.toBe(id1);
    });
  }
});

describe('router access groups', () => {
  const pageHandler = async () => {};

  describe('Group creation before and after AccessGroups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const admin = sourcetool.group('/admin');

    admin.accessGroups('admin');
    admin.accessGroups('super_admin');
    admin.page('/settings', 'Settings', pageHandler);

    const userAdminGroup = admin.group('/');
    userAdminGroup
      .accessGroups('user_admin')
      .page('/users', 'Users', pageHandler);

    const systemAdminGroup = admin.group('/');
    systemAdminGroup
      .accessGroups('system_admin')
      .page('/system', 'System', pageHandler);

    const tests = [
      {
        path: '/admin/settings',
        expectedGroups: ['admin', 'super_admin'],
      },
      {
        path: '/admin/users',
        expectedGroups: ['admin', 'super_admin', 'user_admin'],
      },
      {
        path: '/admin/system',
        expectedGroups: ['admin', 'super_admin', 'system_admin'],
      },
    ];

    for (const t of tests) {
      test(t.path, () => {
        const result = findPageByPath(sourcetool.pages, t.path);
        if (!result) {
          throw new Error(`Page not found: ${t.path}`);
        }
        expect(result.accessGroups.length).toEqual(t.expectedGroups.length);
        expect(
          result.accessGroups.every((group) =>
            t.expectedGroups.includes(group),
          ),
        ).toBe(true);
      });
    }
  });

  describe('Multiple AccessGroups calls', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const admin = sourcetool.group('/admin');
    admin.accessGroups('admin');
    admin.accessGroups('super_admin');
    admin.page('/settings', 'Settings', pageHandler);

    const userAdminGroup = admin.group('/');
    userAdminGroup
      .accessGroups('user_admin')
      .page('/users', 'Users', pageHandler);

    const systemAdminGroup = admin.group('/');
    systemAdminGroup
      .accessGroups('system_admin')
      .page('/system', 'System', pageHandler);

    const tests = [
      {
        path: '/admin/settings',
        expectedGroups: ['admin', 'super_admin'],
      },
      {
        path: '/admin/users',
        expectedGroups: ['admin', 'super_admin', 'user_admin'],
      },
      {
        path: '/admin/system',
        expectedGroups: ['admin', 'super_admin', 'system_admin'],
      },
    ];

    for (const t of tests) {
      test(t.path, () => {
        const result = findPageByPath(sourcetool.pages, t.path);
        if (!result) {
          throw new Error(`Page not found: ${t.path}`);
        }
        expect(result.accessGroups.length).toEqual(t.expectedGroups.length);
        expect(
          result.accessGroups.every((group) =>
            t.expectedGroups.includes(group),
          ),
        ).toBe(true);
      });
    }
  });

  describe('Sibling groups inheritance', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.accessGroups('global');

    const users = sourcetool.group('/users');
    users.accessGroups('user_admin');
    users.page('/list', 'Users', pageHandler);

    const products = sourcetool.group('/products');
    products.accessGroups('product_admin');
    products.page('/list', 'Products', pageHandler);

    const tests = [
      {
        path: '/users/list',
        expectedGroups: ['global', 'user_admin'],
      },
      {
        path: '/products/list',
        expectedGroups: ['global', 'product_admin'],
      },
    ];

    for (const t of tests) {
      test(t.path, () => {
        const result = findPageByPath(sourcetool.pages, t.path);
        if (!result) {
          throw new Error(`Page not found: ${t.path}`);
        }
        expect(result.accessGroups.length).toEqual(t.expectedGroups.length);
        expect(
          result.accessGroups.every((group) =>
            t.expectedGroups.includes(group),
          ),
        ).toBe(true);
      });
    }
  });

  test('Deep nested groups inheritance', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.accessGroups('global');

    const api = sourcetool.group('/api');
    api.accessGroups('api_user');

    const v1 = api.group('/v1');
    v1.accessGroups('v1_user');

    const users = v1.group('/users');
    users.accessGroups('user_admin');

    const settings = users.group('/settings');
    settings.accessGroups('settings_admin');
    settings.page('/profile', 'Profile Settings', pageHandler);

    const page = findPageByPath(
      sourcetool.pages,
      '/api/v1/users/settings/profile',
    );

    const expectedGroups = [
      'global',
      'api_user',
      'v1_user',
      'user_admin',
      'settings_admin',
    ];

    if (!page) {
      throw new Error('Page not found');
    }

    expect(page.accessGroups.length).toEqual(expectedGroups.length);
    expect(
      page.accessGroups.every((group) => expectedGroups.includes(group)),
    ).toBe(true);
  });

  describe('Mixed group and page specific access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const admin = sourcetool.group('/admin');
    admin.accessGroups('admin');
    admin.page('/dashboard', 'Dashboard', pageHandler);

    const settings = admin.group('/settings');
    settings.accessGroups('settings_admin');
    settings.page('/general', 'General Settings', pageHandler);

    const users = settings.group('/users');
    users.accessGroups('user_manager');
    users.page('/profiles', 'User Profiles', pageHandler);

    const tests = [
      {
        path: '/admin/dashboard',
        expectedGroups: ['admin'],
      },
      {
        path: '/admin/settings/general',
        expectedGroups: ['admin', 'settings_admin'],
      },
      {
        path: '/admin/settings/users/profiles',
        expectedGroups: ['admin', 'settings_admin', 'user_manager'],
      },
    ];

    for (const t of tests) {
      test(t.path, () => {
        const result = findPageByPath(sourcetool.pages, t.path);
        if (!result) {
          throw new Error(`Page not found: ${t.path}`);
        }
        expect(result.accessGroups.length).toEqual(t.expectedGroups.length);
        expect(
          result.accessGroups.every((group) =>
            t.expectedGroups.includes(group),
          ),
        ).toBe(true);
      });
    }
  });
});

describe('router groups', () => {
  const pageHandler = async () => {};

  test('group access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const admin = sourcetool.group('/admin');
    const settings = admin.group('/settings');
    settings.page('/users', 'User Settings', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/admin/settings/users');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/admin/settings/users');
  });

  test('Multiple nested groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const api = sourcetool.group('/api');
    const v1 = api.group('/v1');
    const users = v1.group('/users');

    users.page('/list', 'Users List', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/api/v1/users/list');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/api/v1/users/list');
  });
});

describe('router page', () => {
  const pageHandler = async () => {};

  test('Basic page', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.page('/test', 'Test Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/test');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/test');
  });

  test('Skip top-level root path', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);

    sourcetool.page('/', 'Root Page', pageHandler);
    console.log(sourcetool.pages);
    const page = findPageByPath(sourcetool.pages, '/');
    expect(page).toBe(null);

    sourcetool.page('/other', 'Other Page', pageHandler);
    const otherPage = findPageByPath(sourcetool.pages, '/other');
    if (!otherPage) {
      throw new Error('Page not found');
    }
    expect(otherPage.route).toBe('/other');
  });

  test('Allow nested root path', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const users = sourcetool.group('/users');
    users.page('/', 'Users List', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/users');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/users');

    const admin = sourcetool.group('/admin');
    const settings = admin.group('/settings');
    settings.page('/', 'Settings Home', pageHandler);

    const settingsPage = findPageByPath(sourcetool.pages, '/admin/settings');
    if (!settingsPage) {
      throw new Error('Page not found');
    }
    expect(settingsPage.route).toBe('/admin/settings');
  });

  test('Page with access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.accessGroups('admin');
    sourcetool.page('/admin', 'Admin Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/admin');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/admin');
    expect(page.accessGroups.length).toEqual(1);
    expect(page.accessGroups[0]).toBe('admin');
  });

  test('Page with error handler', async () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.page('/error', 'Error Page', async () => {
      throw new Error('Test error');
    });

    const page = findPageByPath(sourcetool.pages, '/error');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/error');

    try {
      await (page.run as any)();
    } catch (error) {
      expect(error).toBeInstanceOf(Error);
      expect((error as Error).message).toBe('Test error');
    }
  });

  test('Page with empty route', async () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.page('', 'Root Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/');
  });

  test('Page with duplicate route', async () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.page('/duplicate', 'First Page', pageHandler);
    sourcetool.page('/duplicate', 'Second Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/duplicate');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/duplicate');
    expect(page.name).toBe('Second Page');
  });
});

describe('router group', () => {
  const pageHandler = async () => {};
  test('Basic group', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const group = sourcetool.group('/test');

    group.page('/page', 'Test Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/test/page');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/test/page');
  });

  test('Group with access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const group = sourcetool.group('/admin');
    group.accessGroups('admin');
    group.page('/dashboard', 'Admin Dashboard', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/admin/dashboard');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/admin/dashboard');
    expect(page.accessGroups.length).toEqual(1);
    expect(page.accessGroups[0]).toBe('admin');
  });

  test('Nested groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const parent = sourcetool.group('/parent');
    const child = parent.group('/child');
    child.page('/page', 'Test Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/parent/child/page');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/parent/child/page');
  });

  test('Group with empty path', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    const group = sourcetool.group('');
    group.page('/page', 'Test Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/page');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/page');
  });
});

describe('router access groups', () => {
  const pageHandler = async () => {};
  test('Set access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.accessGroups('admin', 'user');
    sourcetool.page('/test', 'Test Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/test');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/test');
    expect(page.accessGroups.length).toEqual(2);
    expect(page.accessGroups[0]).toBe('admin');
    expect(page.accessGroups[1]).toBe('user');
  });

  test('Clear access groups', () => {
    const config: SourcetoolConfig = {
      apiKey: 'test_api_key',
      endpoint: 'ws://test.trysourcetool.com',
    };

    const sourcetool = createSourcetool(config);
    sourcetool.accessGroups('admin');
    sourcetool.accessGroups();

    sourcetool.page('/test', 'Test Page', pageHandler);

    const page = findPageByPath(sourcetool.pages, '/test');
    if (!page) {
      throw new Error('Page not found');
    }
    expect(page.route).toBe('/test');
    expect(page.accessGroups.length).toEqual(1);
    expect(page.accessGroups[0]).toBe('admin');
  });
});
