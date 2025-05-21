import { v4 as uuidv4 } from 'uuid';
import { WidgetState } from './widget';

export const WidgetTypeForm = 'form' as const;

export class FormState implements WidgetState {
  constructor(
    public id: string = uuidv4(),
    public value: boolean = false,
    public buttonLabel: string = 'Submit',
    public buttonDisabled: boolean = false,
    public clearOnSubmit: boolean = false,
  ) {
    this.type = WidgetTypeForm;
  }

  getType(): 'form' {
    return WidgetTypeForm;
  }

  public type: 'form' = WidgetTypeForm;
}
