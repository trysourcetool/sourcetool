import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeTextArea = 'textArea' as const;

export class TextAreaState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: string | null = null,
    public label: string = '',
    public placeholder: string = '',
    public defaultValue: string | null = null,
    public required: boolean = false,
    public disabled: boolean = false,
    public maxLength: number | null = null,
    public minLength: number | null = null,
    public maxLines: number | null = null,
    public minLines: number | null = 2,
    public autoResize: boolean = true,
  ) {
    this.type = WidgetTypeTextArea;
  }

  getType(): 'textArea' {
    return WidgetTypeTextArea;
  }

  public type: 'textArea' = WidgetTypeTextArea;
}
