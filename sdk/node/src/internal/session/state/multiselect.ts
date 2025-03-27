import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeMultiSelect: WidgetType = 'multiselect';

export class MultiSelectState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: number[] = [],
    public label: string = '',
    public options: string[] = [],
    public placeholder: string = '',
    public defaultValue: number[] = [],
    public required: boolean = false,
    public disabled: boolean = false,
  ) {
    this.type = WidgetTypeMultiSelect;
  }

  getType(): WidgetType {
    return WidgetTypeMultiSelect;
  }

  public type: WidgetType = WidgetTypeMultiSelect;
}

export interface MultiSelectValue {
  values: string[];
  indexes: number[];
}
