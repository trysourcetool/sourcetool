import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeTextInput = 'textInput' as const;

export class TextInputState implements WidgetState {
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
  ) {
    this.type = WidgetTypeTextInput;
  }

  getType(): 'textInput' {
    return WidgetTypeTextInput;
  }

  public type: 'textInput' = WidgetTypeTextInput;
}
