import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeButton = 'button' as const;

export class ButtonState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: boolean = false,
    public label: string = '',
    public disabled: boolean = false,
  ) {
    this.type = WidgetTypeButton;
  }

  getType(): 'button' {
    return WidgetTypeButton;
  }

  public type: 'button' = WidgetTypeButton;
}
