import { v4 as uuidv4 } from 'uuid';
import { WidgetState, WidgetType } from './widget';

export const WidgetTypeForm: WidgetType = 'form';

export class FormState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: boolean = false,
    public buttonLabel: string = 'Submit',
    public buttonDisabled: boolean = false,
    public clearOnSubmit: boolean = false,
  ) {}

  getType(): WidgetType {
    return WidgetTypeForm;
  }
}
