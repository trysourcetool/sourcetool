import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeSelectbox = 'selectbox' as const;

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
  ) {
    this.type = WidgetTypeSelectbox;
  }

  getType(): 'selectbox' {
    return WidgetTypeSelectbox;
  }

  public type: 'selectbox' = WidgetTypeSelectbox;
}

export interface SelectboxValue {
  value: string;
  index: number;
}
