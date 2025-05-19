import { v4 as uuidv4 } from 'uuid';
import { UIBuilder, uiBuilderGeneratePageID } from '../';
import {
  MarkdownState,
  WidgetTypeMarkdown,
} from '../../session/state/markdown';
import { MarkdownOptions } from '../../types/options';
import { create, fromJson, toJson } from '@bufbuild/protobuf';
import {
  Markdown as MarkdownProto,
  MarkdownSchema,
  WidgetSchema,
} from '../../pb/widget/v1/widget_pb';
import { RenderWidgetSchema } from '../../pb/websocket/v1/message_pb';
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
  const widgetID = uiBuilderGeneratePageID(page.id, WidgetTypeMarkdown, path);

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

  const renderWidget = create(RenderWidgetSchema, {
    sessionId: session.id,
    pageId: page.id,
    path: convertPathToInt32Array(path),
    widget: create(WidgetSchema, {
      id: widgetID,
      type: {
        case: 'markdown',
        value: markdownProto,
      },
    }),
  });

  runtime.wsClient.enqueue(uuidv4(), renderWidget);

  cursor.next();
}

/**
 * Convert markdown state to proto
 * @param state Markdown state
 * @returns Markdown proto
 */
export function convertStateToMarkdownProto(
  state: MarkdownState,
): MarkdownProto {
  return fromJson(MarkdownSchema, {
    body: state.body,
  });
}

/**
 * Convert markdown proto to state
 * @param id Widget ID
 * @param data Markdown proto
 * @returns Markdown state
 */
export function convertMarkdownProtoToState(
  id: string,
  data: MarkdownProto | null,
): MarkdownState | null {
  if (!data) {
    return null;
  }

  const d = toJson(MarkdownSchema, data);

  return new MarkdownState(id, d.body);
}

/**
 * Convert path to int32 array
 * @param path Path
 * @returns Int32 array
 */
export function convertPathToInt32Array(path: number[]): number[] {
  return path.map((v) => v);
}
