import { v5 as uuidv5 } from 'uuid';
import { button, ButtonComponentOptions } from './widgets/button';
import { checkbox, CheckboxComponentOptions } from './widgets/checkbox';
import { markdown } from './widgets/markdown';
import { textInput, TextInputComponentOptions } from './widgets/textinput';
import {
  numberInput,
  NumberInputComponentOptions,
} from './widgets/numberinput';
import { dateInput, DateInputComponentOptions } from './widgets/dateinput';
import {
  dateTimeInput,
  DateTimeInputComponentOptions,
} from './widgets/datetimeinput';
import { timeInput, TimeInputComponentOptions } from './widgets/timeinput';
import { radio, RadioComponentOptions } from './widgets/radio';
import { selectbox, SelectboxComponentOptions } from './widgets/selectbox';
import {
  multiSelect,
  MultiSelectComponentOptions,
} from './widgets/multiselect';
import {
  checkboxGroup,
  CheckboxGroupComponentOptions,
} from './widgets/checkboxgroup';
import { textArea, TextAreaComponentOptions } from './widgets/textarea';
import { table, TableComponentOptions, TableValue } from './widgets/table';
import { form, FormComponentOptions } from './widgets/form';
import { columns, ColumnsComponentOptions } from './widgets/columns';
import { Page } from '../page';
import { Session } from '../session';
import { Runtime } from '../runtime';
import { WidgetType } from '../session/state/widget';
import { SelectboxValue } from '../session/state/selectbox';
import { MultiSelectValue } from '../session/state/multiselect';
import { RadioValue } from '../session/state/radio';
import { CheckboxGroupValue } from '../session/state/checkboxgroup';

export type UIBuilderType = {
  markdown(content: string): void;
  textInput(label: string, options?: TextInputComponentOptions): string;
  numberInput(
    label: string,
    options?: NumberInputComponentOptions,
  ): number | null;
  dateInput(label: string, options?: DateInputComponentOptions): Date | null;
  dateTimeInput(
    label: string,
    options?: DateTimeInputComponentOptions,
  ): Date | null;
  timeInput(label: string, options?: TimeInputComponentOptions): Date | null;
  selectbox(
    label: string,
    options?: SelectboxComponentOptions,
  ): SelectboxValue | null;
  multiSelect(
    label: string,
    options?: MultiSelectComponentOptions,
  ): MultiSelectValue | null;
  radio(label: string, options?: RadioComponentOptions): RadioValue | null;
  checkbox(label: string, options?: CheckboxComponentOptions): boolean;
  checkboxGroup(
    label: string,
    options?: CheckboxGroupComponentOptions,
  ): CheckboxGroupValue | null;
  textArea(label: string, options?: TextAreaComponentOptions): string;
  table(data: any, options?: TableComponentOptions): TableValue | null;
  button(label: string, options?: ButtonComponentOptions): boolean;
  form(label: string, options?: FormComponentOptions): [UIBuilderType, boolean];
  columns(count: number, options?: ColumnsComponentOptions): UIBuilderType[];
};

export const uiBuilderGeneratePageID = (
  pageId: string,
  widgetType: WidgetType,
  path: number[],
): string => {
  const strPath = path.map((v) => v.toString()).join('_');
  return uuidv5(`${widgetType}-${strPath}`, pageId);
};

export class UIBuilder implements UIBuilderType {
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

  textInput(label: string, options: TextInputComponentOptions = {}): string {
    return textInput(this.getContext(), label, options);
  }

  numberInput(
    label: string,
    options: NumberInputComponentOptions = {},
  ): number | null {
    return numberInput(this.getContext(), label, options);
  }

  dateInput(
    label: string,
    options: DateInputComponentOptions = {},
  ): Date | null {
    return dateInput(this.getContext(), label, options);
  }

  dateTimeInput(
    label: string,
    options: DateTimeInputComponentOptions = {},
  ): Date | null {
    return dateTimeInput(this.getContext(), label, options);
  }

  timeInput(
    label: string,
    options: TimeInputComponentOptions = {},
  ): Date | null {
    return timeInput(this.getContext(), label, options);
  }

  selectbox(
    label: string,
    options: SelectboxComponentOptions = {},
  ): SelectboxValue | null {
    return selectbox(this.getContext(), label, options);
  }

  multiSelect(
    label: string,
    options: MultiSelectComponentOptions = {},
  ): MultiSelectValue | null {
    return multiSelect(this.getContext(), label, options);
  }

  radio(label: string, options: RadioComponentOptions = {}): RadioValue | null {
    return radio(this.getContext(), label, options);
  }

  checkbox(label: string, options: CheckboxComponentOptions = {}): boolean {
    return checkbox(this.getContext(), label, options);
  }

  checkboxGroup(
    label: string,
    options: CheckboxGroupComponentOptions = {},
  ): CheckboxGroupValue | null {
    return checkboxGroup(this.getContext(), label, options);
  }

  textArea(label: string, options: TextAreaComponentOptions = {}): string {
    return textArea(this.getContext(), label, options);
  }

  table(data: any, options: TableComponentOptions = {}): TableValue | null {
    return table(this.getContext(), data, options);
  }

  button(label: string, options: ButtonComponentOptions = {}): boolean {
    return button(this.getContext(), label, options);
  }

  form(
    label: string,
    options: FormComponentOptions = {},
  ): [UIBuilderType, boolean] {
    return form(this.getContext(), label, options);
  }

  columns(
    count: number,
    options: ColumnsComponentOptions = {},
  ): UIBuilderType[] {
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
