import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeTimeInput: WidgetType = 'timeInput';

export class TimeInputState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: Date | null = null,
    public label: string = '',
    public placeholder: string = '',
    public defaultValue: Date | null = null,
    public required: boolean = false,
    public disabled: boolean = false,
    public location: string = 'local',
  ) {
    this.type = WidgetTypeTimeInput;
  }

  getType(): WidgetType {
    return WidgetTypeTimeInput;
  }

  public type: WidgetType = WidgetTypeTimeInput;
}
