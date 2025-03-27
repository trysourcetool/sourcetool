import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeCheckbox: WidgetType = 'checkbox';

export class CheckboxState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public label: string = '',
    public value: boolean = false,
    public defaultValue: boolean = false,
    public required: boolean = false,
    public disabled: boolean = false,
  ) {
    this.type = WidgetTypeCheckbox;
  }

  getType(): WidgetType {
    return WidgetTypeCheckbox;
  }

  public type: WidgetType = WidgetTypeCheckbox;
}
