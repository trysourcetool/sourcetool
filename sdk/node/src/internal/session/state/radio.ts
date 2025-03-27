import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeRadio: WidgetType = 'radio';

export class RadioState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: number | null = null,
    public label: string = '',
    public options: string[] = [],
    public defaultValue: number | null = null,
    public required: boolean = false,
    public disabled: boolean = false,
  ) {}

  getType(): WidgetType {
    return WidgetTypeRadio;
  }
}

export interface RadioValue {
  value: string;
  index: number;
}
