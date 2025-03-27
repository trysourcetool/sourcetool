import { ButtonState } from './button';
import { CheckboxState } from './checkbox';
import { CheckboxGroupState } from './checkboxgroup';
import { ColumnsState } from './columns';
import { DateInputState } from './dateinput';
import { DateTimeInputState } from './datetimeinput';
import { FormState } from './form';
import { MarkdownState } from './markdown';
import { MultiSelectState } from './multiselect';
import { NumberInputState } from './numberinput';
import { RadioState } from './radio';
import { SelectboxState } from './selectbox';
import { TableState } from './table';
import { TextAreaState } from './textarea';
import { TextInputState } from './textinput';
import { TimeInputState } from './timeinput';

export type State = {
  button: ButtonState;
  textInput: TextInputState;
  numberInput: NumberInputState;
  dateInput: DateInputState;
  datetimeInput: DateTimeInputState;
  timeInput: TimeInputState;
  selectbox: SelectboxState;
  multiselect: MultiSelectState;
  radio: RadioState;
  checkbox: CheckboxState;
  checkboxGroup: CheckboxGroupState;
  textArea: TextAreaState;
  table: TableState;
  form: FormState;
  columns: ColumnsState;
  markdown: MarkdownState;
};

export type WidgetType =
  | 'button'
  | 'textInput'
  | 'numberInput'
  | 'dateInput'
  | 'datetimeInput'
  | 'timeInput'
  | 'selectbox'
  | 'multiselect'
  | 'radio'
  | 'checkbox'
  | 'checkboxGroup'
  | 'textArea'
  | 'table'
  | 'form'
  | 'columns'
  | 'columnItem'
  | 'markdown';

export interface WidgetState {
  getType(): WidgetType;
}

export function String(w: WidgetType) {
  return w;
}
