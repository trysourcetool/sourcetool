import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeNumberInput: WidgetType = 'numberInput';

export class NumberInputState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: number | null = null,
    public label: string = '',
    public placeholder: string = '',
    public defaultValue: number | null = null,
    public required: boolean = false,
    public disabled: boolean = false,
    public maxValue: number | null = null,
    public minValue: number | null = null,
  ) {
    this.type = WidgetTypeNumberInput;
  }

  getType(): WidgetType {
    return WidgetTypeNumberInput;
  }

  public type: WidgetType = WidgetTypeNumberInput;
}
