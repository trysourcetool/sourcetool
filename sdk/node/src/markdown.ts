import { v4 as uuidv4 } from 'uuid';
import { UIBuilder } from './uibuilder';
import {
  MarkdownState,
  WidgetTypeMarkdown,
} from './internal/session/state/markdown';
import { MarkdownOptions } from './internal/options';

/**
 * Add markdown content to the UI
 * @param builder The UI builder
 * @param body The markdown content
 */
export function markdown(builder: UIBuilder, body: string): void {
  const runtime = builder.runtime;
  const session = builder.session;
  const page = builder.page;
  const cursor = builder.cursor;

  if (!session || !page || !cursor) {
    return;
  }

  const markdownOpts: MarkdownOptions = {
    body,
  };

  const path = cursor.getPath();
  const widgetID = builder.generatePageID(WidgetTypeMarkdown, path);

  let markdownState = session.state.getMarkdown(widgetID);
  if (!markdownState) {
    markdownState = new MarkdownState(widgetID, markdownOpts.body);
  } else {
    markdownState.body = markdownOpts.body;
  }

  session.state.set(widgetID, markdownState);

  const markdownProto = convertStateToMarkdownProto(
    markdownState as MarkdownState,
  );
  runtime.wsClient.enqueue(uuidv4(), {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: {
      id: widgetID,
      type: 'Markdown',
      markdown: markdownProto,
    },
  });

  cursor.next();
}

/**
 * Convert markdown state to proto
 * @param state Markdown state
 * @returns Markdown proto
 */
function convertStateToMarkdownProto(state: MarkdownState): any {
  return {
    body: state.body,
  };
}

/**
 * Convert markdown proto to state
 * @param id Widget ID
 * @param data Markdown proto
 * @returns Markdown state
 */
export function convertMarkdownProtoToState(
  id: string,
  data: any,
): MarkdownState | null {
  if (!data) {
    return null;
  }

  return new MarkdownState(id, data.body);
}

/**
 * Convert path to int32 array
 * @param path Path
 * @returns Int32 array
 */
export function convertPathToInt32Array(path: number[]): number[] {
  return path.map((v) => v);
}
