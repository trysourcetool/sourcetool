import { v5 as uuidv5 } from 'uuid';
import { UIBuilder } from './uibuilder';
import { Page } from './page';

export function removeDuplicates(groups: string[]): string[] {
  return [...new Set(groups)];
}

/**
 * Router interface
 */
export interface Router {
  /**
   * Add a page to the router
   * @param relativePath Relative path
   * @param name Page name
   * @param handler Page handler
   */
  page(
    relativePath: string,
    name: string,
    handler: (ui: UIBuilder) => Promise<void>,
  ): void;

  /**
   * Add access groups to the router
   * @param groups Access groups
   * @returns Router
   */
  accessGroups(...groups: string[]): Router;

  /**
   * Create a new router group
   * @param relativePath Relative path
   * @returns Router
   */
  group(relativePath: string): Router;
}

type RouterContext = {
  environment: string;
  pages: Record<string, Page>;
  addPage: (id: string, page: Page) => void;
};

/**
 * Router class
 */
export class RouterImpl implements Router {
  private parent: RouterImpl | null;
  private context: RouterContext | null;
  private basePath: string;
  private namespaceDNS: string;
  private groups: string[];

  /**
   * Constructor
   * @param context Context
   * @param namespaceDNS Namespace DNS
   * @param parent Parent router
   * @param basePath Base path
   * @param groups Access groups
   */
  constructor(
    context: RouterContext | null,
    namespaceDNS: string,
    parent: RouterImpl | null = null,
    basePath: string = '',
    groups: string[] = [],
  ) {
    this.parent = parent;
    this.context = context;
    this.basePath = basePath;
    this.namespaceDNS = namespaceDNS;
    this.groups = groups;
  }

  /**
   * Generate a page ID
   * @param fullPath Full path
   * @returns Page ID
   */
  generatePageId(fullPath: string): string {
    const ns = uuidv5(this.namespaceDNS, uuidv5.DNS);
    return uuidv5(`${fullPath}-${this.context?.environment}`, ns);
  }

  /**
   * Join a relative path with the base path
   * @param relativePath Relative path
   * @returns Full path
   */
  joinPath(relativePath: string): string {
    if (!relativePath.startsWith('/')) {
      relativePath = '/' + relativePath;
    }
    if (this.basePath === '') {
      if (relativePath === '/') {
        return relativePath;
      }
      if (relativePath.endsWith('/')) {
        return relativePath.slice(0, -1);
      }
      return relativePath;
    }
    const basePath = this.basePath.endsWith('/')
      ? this.basePath.slice(0, -1)
      : this.basePath;
    const cleanPath = relativePath.startsWith('/')
      ? relativePath.slice(1)
      : relativePath;
    const result = basePath + '/' + cleanPath;

    if (result === '/') {
      return result;
    }
    if (result.endsWith('/')) {
      return result.slice(0, -1);
    }
    return result;
  }

  /**
   * Remove duplicates from an array
   * @param groups Array of strings
   * @returns Array with duplicates removed
   */
  private removeDuplicates(groups: string[]): string[] {
    return removeDuplicates(groups);
  }

  /**
   * Collect all access groups from the router chain
   * @returns Access groups
   */
  private collectGroups(): string[] {
    const groups: string[] = [];
    let current: RouterImpl | null = this;

    while (current !== null) {
      groups.push(...current.groups);
      current = current.parent;
    }

    return groups;
  }

  /**
   * Add a page to the router
   * @param relativePath Relative path
   * @param name Page name
   * @param handler Page handler
   */
  page(
    relativePath: string,
    name: string,
    handler: (ui: UIBuilder) => Promise<void>,
  ): void {
    // Skip page creation only for top-level root path
    if (relativePath === '/' && this.basePath === '') {
      return;
    }

    let fullPath: string;

    if (relativePath === '') {
      if (this.basePath === '') {
        fullPath = '/';
      } else {
        fullPath = this.basePath.endsWith('/')
          ? this.basePath.slice(0, -1)
          : this.basePath;
      }
    } else {
      fullPath = this.joinPath(relativePath);
    }

    const pageId = this.generatePageId(fullPath);

    if (this.context === null) {
      throw new Error('Sourcetool is not set');
    }

    const page = new Page(
      pageId,
      name,
      fullPath,
      [Object.keys(this.context.pages).length],
      handler,
      this.removeDuplicates(this.collectGroups()),
    );

    this.context.addPage(pageId, page);
  }

  /**
   * Add access groups to the router
   * @param groups Access groups
   * @returns Router
   */
  accessGroups(...groups: string[]): Router {
    if (groups.length > 0) {
      this.groups.push(...groups);
    }
    return this;
  }

  /**
   * Create a new router group
   * @param relativePath Relative path
   * @returns Router
   */
  group(relativePath: string): Router {
    if (this.context === null) {
      throw new Error('Sourcetool is not set');
    }
    return new RouterImpl(
      this.context,
      this.namespaceDNS,
      this,
      this.joinPath(relativePath),
      [],
    );
  }
}
