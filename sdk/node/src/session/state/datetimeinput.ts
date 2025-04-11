import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeDateTimeInput = 'datetimeInput' as const;

export class DateTimeInputState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: Date | null = null,
    public label: string = '',
    public placeholder: string = '',
    public defaultValue: Date | null = null,
    public required: boolean = false,
    public disabled: boolean = false,
    public format: string = 'YYYY/MM/DD HH:MM:SS',
    public maxValue: Date | null = null,
    public minValue: Date | null = null,
    public location: string = 'local',
  ) {
    this.type = WidgetTypeDateTimeInput;
  }

  getType(): 'datetimeInput' {
    return WidgetTypeDateTimeInput;
  }

  public type: 'datetimeInput' = WidgetTypeDateTimeInput;
}
