import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeCheckboxGroup = 'checkboxGroup' as const;

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

  getType(): 'checkboxGroup' {
    return WidgetTypeCheckboxGroup;
  }

  public type: 'checkboxGroup' = WidgetTypeCheckboxGroup;
}

export interface CheckboxGroupValue {
  values: string[];
  indexes: number[];
}
