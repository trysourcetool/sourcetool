import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeTextInput: WidgetType = 'textInput';

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

  getType(): WidgetType {
    return WidgetTypeTextInput;
  }

  public type: WidgetType = WidgetTypeTextInput;
}
