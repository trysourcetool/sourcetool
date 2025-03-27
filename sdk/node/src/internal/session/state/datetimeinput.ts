import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeDateTimeInput: WidgetType = 'datetimeInput';

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
  ) {}

  getType(): WidgetType {
    return WidgetTypeDateTimeInput;
  }
}
