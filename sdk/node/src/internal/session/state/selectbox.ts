import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeSelectbox: WidgetType = 'selectbox';

export class SelectboxState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: number | null = null,
    public label: string = '',
    public options: string[] = [],
    public placeholder: string = '',
    public defaultValue: number | null = null,
    public required: boolean = false,
    public disabled: boolean = false,
  ) {}

  getType(): WidgetType {
    return WidgetTypeSelectbox;
  }
}

export interface SelectboxValue {
  value: string;
  index: number;
}
