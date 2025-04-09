import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeCheckbox = 'checkbox' as const;

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

  getType(): 'checkbox' {
    return WidgetTypeCheckbox;
  }

  public type: 'checkbox' = WidgetTypeCheckbox;
}
