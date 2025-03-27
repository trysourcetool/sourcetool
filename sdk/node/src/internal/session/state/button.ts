import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeButton: WidgetType = 'button';

export class ButtonState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: boolean = false,
    public label: string = '',
    public disabled: boolean = false,
  ) {
    this.type = WidgetTypeButton;
  }

  getType(): WidgetType {
    return WidgetTypeButton;
  }

  public type: WidgetType = WidgetTypeButton;
}
