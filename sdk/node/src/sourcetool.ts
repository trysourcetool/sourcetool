import { UIBuilder } from './uibuilder';
import { Page, PageManager, newPageManager } from './internal/page';
import { RouterInterface, newRouter } from './router';
import { startRuntime } from './runtime';

/**
 * Sourcetool configuration
 */
export interface SourcetoolConfig {
  /**
   * API key
   */
  apiKey: string;

  /**
   * Endpoint URL
   */
  endpoint: string;
}

/**
 * Sourcetool class
 */
export class Sourcetool implements RouterInterface {
  /**
   * API key
   */
  apiKey: string;

  /**
   * Endpoint URL
   */
  endpoint: string;

  /**
   * Router
   */
  router: RouterInterface;

  /**
   * Pages
   */
  pages: Record<string, Page>;

  /**
   * Page manager
   */
  pageManager: PageManager;

  /**
   * Runtime
   */
  runtime: any;

  /**
   * Constructor
   * @param config Configuration
   */
  constructor(config: SourcetoolConfig) {
    this.apiKey = config.apiKey;

    // Format endpoint URL
    const hostParts = config.endpoint.split('://');
    if (hostParts.length !== 2) {
      throw new Error('Invalid endpoint URL');
    }

    this.endpoint = `${config.endpoint}/ws`;

    // Extract namespace DNS
    const namespaceDNS = hostParts[1].split(':')[0];

    // Initialize pages
    this.pages = {};

    // Initialize router
    this.router = newRouter(this, namespaceDNS);

    // Initialize page manager
    this.pageManager = newPageManager();
  }

  /**
   * Add a page
   * @param id Page ID
   * @param page Page
   */
  addPage(id: string, page: Page): void {
    this.pages[id] = page;
  }

  /**
   * Start the server
   * @returns Promise
   */
  async listen(): Promise<void> {
    // Validate pages
    this.validatePages();

    // Initialize logger
    // await initLogger();

    // Start runtime
    const runtime = await startRuntime(this.apiKey, this.endpoint, this.pages);
    this.runtime = runtime;

    // Wait for runtime to finish
    await runtime.wsClient.wait();
  }

  /**
   * Close the server
   * @returns Promise
   */
  async close(): Promise<void> {
    if (this.runtime) {
      await this.runtime.wsClient.close();
    }
  }

  /**
   * Validate pages
   */
  private validatePages(): void {
    const pagesByRoute: Record<string, string> = {};

    // Find duplicate routes
    for (const [id, page] of Object.entries(this.pages)) {
      pagesByRoute[page.route] = id;
    }

    // Create new pages object with only unique routes
    const newPages: Record<string, Page> = {};
    for (const [, id] of Object.entries(pagesByRoute)) {
      newPages[id] = this.pages[id];
    }

    this.pages = newPages;
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
    this.router.page(relativePath, name, handler);
  }

  /**
   * Add access groups to the router
   * @param groups Access groups
   * @returns Router
   */
  accessGroups(...groups: string[]): RouterInterface {
    return this.router.accessGroups(...groups);
  }

  /**
   * Create a new router group
   * @param relativePath Relative path
   * @returns Router
   */
  group(relativePath: string): RouterInterface {
    return this.router.group(relativePath);
  }
}

/**
 * Create a new Sourcetool instance
 * @param config Configuration
 * @returns Sourcetool instance
 */
export function createSourcetool(config: SourcetoolConfig): Sourcetool {
  return new Sourcetool(config);
}
