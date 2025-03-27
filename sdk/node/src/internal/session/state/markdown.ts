import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeMarkdown: WidgetType = 'markdown';

export class MarkdownState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public body: string = '',
  ) {
    this.type = WidgetTypeMarkdown;
  }

  getType(): WidgetType {
    return WidgetTypeMarkdown;
  }

  public type: WidgetType = WidgetTypeMarkdown;
}
