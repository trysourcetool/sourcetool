import { v5 as uuidv5 } from 'uuid';
import { UIBuilder } from './uibuilder';
import { Page } from './page';
import { Sourcetool } from './sourcetool';

export function removeDuplicates(groups: string[]): string[] {
  return [...new Set(groups)];
}

/**
 * Router interface
 */
export interface RouterInterface {
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
  accessGroups(...groups: string[]): RouterInterface;

  /**
   * Create a new router group
   * @param relativePath Relative path
   * @returns Router
   */
  group(relativePath: string): RouterInterface;
}

/**
 * Router class
 */
export class Router implements RouterInterface {
  private parent: Router | null;
  private sourcetool: Sourcetool | null;
  private basePath: string;
  private namespaceDNS: string;
  private groups: string[];

  /**
   * Constructor
   * @param sourcetool Sourcetool instance
   * @param namespaceDNS Namespace DNS
   * @param parent Parent router
   * @param basePath Base path
   * @param groups Access groups
   */
  constructor(
    sourcetool: Sourcetool | null = null,
    namespaceDNS: string,
    parent: Router | null = null,
    basePath: string = '',
    groups: string[] = [],
  ) {
    this.parent = parent;
    this.sourcetool = sourcetool;
    this.basePath = basePath;
    this.namespaceDNS = namespaceDNS;
    this.groups = groups;
  }

  /**
   * Generate a page ID
   * @param fullPath Full path
   * @returns Page ID
   */
  generatePageID(fullPath: string): string {
    const ns = uuidv5(this.namespaceDNS, uuidv5.DNS);
    return uuidv5(`${fullPath}-${this.sourcetool?.environment}`, ns);
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
    let current: Router | null = this;

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

    const pageID = this.generatePageID(fullPath);

    if (this.sourcetool === null) {
      throw new Error('Sourcetool is not set');
    }

    const page = new Page(
      pageID,
      name,
      fullPath,
      [Object.keys(this.sourcetool.pages).length],
      handler,
      this.removeDuplicates(this.collectGroups()),
    );

    this.sourcetool.addPage(pageID, page);
  }

  /**
   * Add access groups to the router
   * @param groups Access groups
   * @returns Router
   */
  accessGroups(...groups: string[]): RouterInterface {
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
  group(relativePath: string): RouterInterface {
    return new Router(
      this.sourcetool,
      this.namespaceDNS,
      this,
      this.joinPath(relativePath),
      [],
    );
  }
}
