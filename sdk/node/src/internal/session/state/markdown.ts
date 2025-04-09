import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeMarkdown = 'markdown' as const;

export class MarkdownState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public body: string = '',
  ) {
    this.type = WidgetTypeMarkdown;
  }

  getType(): 'markdown' {
    return WidgetTypeMarkdown;
  }

  public type: 'markdown' = WidgetTypeMarkdown;
}
