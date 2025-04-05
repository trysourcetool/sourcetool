import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from '../uibuilder';

/**
 * Page class
 */
export class Page {
  id: string;
  name: string;
  route: string;
  path: number[];
  handler: (ui: UIBuilder) => Promise<void>;
  accessGroups: string[];

  /**
   * Constructor
   * @param id Page ID
   * @param name Page name
   * @param route Page route
   * @param path Page path
   * @param handler Page handler
   * @param accessGroups Access groups
   */
  constructor(
    id: string = uuidv4(),
    name: string = '',
    route: string = '',
    path: number[] = [],
    handler: (ui: UIBuilder) => Promise<void> = async () => {},
    accessGroups: string[] = [],
  ) {
    this.id = id;
    this.name = name;
    this.route = route;
    this.path = path;
    this.handler = handler;
    this.accessGroups = accessGroups;
  }

  /**
   * Run the page handler
   * @param ui UI builder
   * @returns Promise
   */
  async run(ui: UIBuilder): Promise<void> {
    await this.handler(ui);
  }

  /**
   * Check if the user has access to the page
   * @param userGroups User groups
   * @returns Whether the user has access
   */
  hasAccess(userGroups: string[]): boolean {
    if (this.accessGroups.length === 0) {
      return true;
    }

    for (const userGroup of userGroups) {
      for (const requiredGroup of this.accessGroups) {
        if (userGroup === requiredGroup) {
          return true;
        }
      }
    }
    return false;
  }
}

/**
 * Page manager class
 */
export class PageManager {
  private pages: Map<string, Page>;

  /**
   * Constructor
   * @param pages Pages
   */
  constructor(pages: Record<string, Page> = {}) {
    this.pages = new Map();

    // Convert record to map
    for (const [id, page] of Object.entries(pages)) {
      this.pages.set(id, page);
    }
  }

  /**
   * Get a page by ID
   * @param id Page ID
   * @returns Page
   */
  getPage(id: string): Page | undefined {
    return this.pages.get(id);
  }

  /**
   * Add a page
   * @param id Page ID
   * @param page Page
   */
  addPage(id: string, page: Page): void {
    this.pages.set(id, page);
  }

  /**
   * Get all pages
   * @returns Pages
   */
  getPages(): Map<string, Page> {
    return this.pages;
  }
}

/**
 * Create a new page manager
 * @param pages Pages
 * @returns Page manager
 */
export function newPageManager(pages: Record<string, Page> = {}): PageManager {
  return new PageManager(pages);
}
