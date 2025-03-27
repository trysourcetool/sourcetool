import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeCheckboxGroup: WidgetType = 'checkboxGroup';

export class CheckboxGroupState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: number[] = [],
    public label: string = '',
    public options: string[] = [],
    public defaultValue: number[] = [],
    public required: boolean = false,
    public disabled: boolean = false,
  ) {
    this.type = WidgetTypeCheckboxGroup;
  }

  getType(): WidgetType {
    return WidgetTypeCheckboxGroup;
  }

  public type: WidgetType = WidgetTypeCheckboxGroup;
}

export interface CheckboxGroupValue {
  values: string[];
  indexes: number[];
}
