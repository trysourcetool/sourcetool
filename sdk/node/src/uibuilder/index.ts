import { v5 as uuidv5 } from 'uuid';
import { button, ButtonOptions } from './widgets/button';
import { checkbox, CheckboxOptions } from './widgets/checkbox';
import { markdown } from './widgets/markdown';
import { textInput, TextInputOptions } from './widgets/textinput';
import { numberInput, NumberInputOptions } from './widgets/numberinput';
import { dateInput, DateInputOptions } from './widgets/dateinput';
import { dateTimeInput, DateTimeInputOptions } from './widgets/datetimeinput';
import { timeInput, TimeInputOptions } from './widgets/timeinput';
import { radio, RadioOptions } from './widgets/radio';
import { selectbox, SelectboxOptions } from './widgets/selectbox';
import { multiSelect, MultiSelectOptions } from './widgets/multiselect';
import { checkboxGroup, CheckboxGroupOptions } from './widgets/checkboxgroup';
import { textArea, TextAreaOptions } from './widgets/textarea';
import { table, TableOptions, TableValue } from './widgets/table';
import { form, FormOptions } from './widgets/form';
import { columns, ColumnsOptions } from './widgets/columns';
import { Page } from '../page';
import { Session } from '../session';
import { Runtime } from '../runtime';
import { WidgetType } from '../session/state/widget';
import { SelectboxValue } from '../session/state/selectbox';
import { MultiSelectValue } from '../session/state/multiselect';
import { RadioValue } from '../session/state/radio';
import { CheckboxGroupValue } from '../session/state/checkboxgroup';

export interface UIBuilder {
  markdown(content: string): void;
  textInput(label: string, options?: TextInputOptions): string;
  numberInput(label: string, options?: NumberInputOptions): number | null;
  dateInput(label: string, options?: DateInputOptions): Date | null;
  dateTimeInput(label: string, options?: DateTimeInputOptions): Date | null;
  timeInput(label: string, options?: TimeInputOptions): Date | null;
  selectbox(label: string, options?: SelectboxOptions): SelectboxValue | null;
  multiSelect(
    label: string,
    options?: MultiSelectOptions,
  ): MultiSelectValue | null;
  radio(label: string, options?: RadioOptions): RadioValue | null;
  checkbox(label: string, options?: CheckboxOptions): boolean;
  checkboxGroup(
    label: string,
    options?: CheckboxGroupOptions,
  ): CheckboxGroupValue | null;
  textArea(label: string, options?: TextAreaOptions): string;
  table(data: any, options?: TableOptions): TableValue | null;
  button(label: string, options?: ButtonOptions): boolean;
  form(label: string, options?: FormOptions): [UIBuilder, boolean];
  columns(count: number, options?: ColumnsOptions): UIBuilder[];
}

export const uiBuilderGeneratePageId = (
  pageId: string,
  widgetType: WidgetType,
  path: number[],
): string => {
  const strPath = path.map((v) => v.toString()).join('_');
  return uuidv5(`${widgetType}-${strPath}`, pageId);
};

export class UIBuilderImpl implements UIBuilder {
  private runtime: Runtime;
  private cursor: Cursor;
  private session: Session;
  private page: Page;

  constructor(runtime: Runtime, session: Session, page: Page, cursor?: Cursor) {
    this.runtime = runtime;
    this.cursor = cursor ?? new Cursor();
    this.session = session;
    this.page = page;
  }

  private getContext(): {
    runtime: Runtime;
    cursor: Cursor;
    session: Session;
    page: Page;
  } {
    return {
      runtime: this.runtime,
      cursor: this.cursor,
      session: this.session,
      page: this.page,
    };
  }
  markdown(content: string): void {
    markdown(this.getContext(), content);
  }

  textInput(label: string, options: TextInputOptions = {}): string {
    return textInput(this.getContext(), label, options);
  }

  numberInput(label: string, options: NumberInputOptions = {}): number | null {
    return numberInput(this.getContext(), label, options);
  }

  dateInput(label: string, options: DateInputOptions = {}): Date | null {
    return dateInput(this.getContext(), label, options);
  }

  dateTimeInput(
    label: string,
    options: DateTimeInputOptions = {},
  ): Date | null {
    return dateTimeInput(this.getContext(), label, options);
  }

  timeInput(label: string, options: TimeInputOptions = {}): Date | null {
    return timeInput(this.getContext(), label, options);
  }

  selectbox(
    label: string,
    options: SelectboxOptions = {},
  ): SelectboxValue | null {
    return selectbox(this.getContext(), label, options);
  }

  multiSelect(
    label: string,
    options: MultiSelectOptions = {},
  ): MultiSelectValue | null {
    return multiSelect(this.getContext(), label, options);
  }

  radio(label: string, options: RadioOptions = {}): RadioValue | null {
    return radio(this.getContext(), label, options);
  }

  checkbox(label: string, options: CheckboxOptions = {}): boolean {
    return checkbox(this.getContext(), label, options);
  }

  checkboxGroup(
    label: string,
    options: CheckboxGroupOptions = {},
  ): CheckboxGroupValue | null {
    return checkboxGroup(this.getContext(), label, options);
  }

  textArea(label: string, options: TextAreaOptions = {}): string {
    return textArea(this.getContext(), label, options);
  }

  table(data: any, options: TableOptions = {}): TableValue | null {
    return table(this.getContext(), data, options);
  }

  button(label: string, options: ButtonOptions = {}): boolean {
    return button(this.getContext(), label, options);
  }

  form(label: string, options: FormOptions = {}): [UIBuilder, boolean] {
    return form(this.getContext(), label, options);
  }

  columns(count: number, options: ColumnsOptions = {}): UIBuilder[] {
    return columns(this.getContext(), count, options);
  }
}

export class Cursor {
  parentPath: number[];
  private index: number;

  constructor() {
    this.parentPath = [];
    this.index = 0;
  }

  getPath(): number[] {
    return [...this.parentPath, this.index];
  }

  next(): void {
    this.index++;
  }
}
