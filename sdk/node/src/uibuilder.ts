import { v5 as uuidv5, NIL as uuidNil } from 'uuid';
import { button } from './button';
import { checkbox } from './checkbox';
import { markdown } from './markdown';
import { textInput } from './textinput';
import { numberInput } from './numberinput';
import { dateInput } from './dateinput';
import { dateTimeInput } from './datetimeinput';
import { timeInput } from './timeinput';
import { radio } from './radio';
import { selectbox } from './selectbox';
import { multiSelect } from './multiselect';
import { checkboxGroup } from './checkboxgroup';
import { textArea } from './textarea';
import { table } from './table';
import { form } from './form';
import { columns } from './columns';
import { Page } from './internal/page';
import { Session } from './internal/session';
import { Runtime } from './runtime';

export type UIBuilderType = {
  markdown(content: string): void;
  textInput(label: string, options?: any): string;
  numberInput(label: string, options?: any): number | null;
  dateInput(label: string, options?: any): Date | null;
  dateTimeInput(label: string, options?: any): Date | null;
  timeInput(label: string, options?: any): Date | null;
  selectbox(label: string, options?: any): any;
  multiSelect(label: string, options?: any): any;
  radio(label: string, options?: any): any;
  checkbox(label: string, options?: any): boolean;
  checkboxGroup(label: string, options?: any): any;
  textArea(label: string, options?: any): string;
  table(data: any, options?: any): any;
  button(label: string, options?: any): boolean;
  form(label: string, options?: any): [UIBuilderType, boolean];
  columns(count: number, options?: any): UIBuilderType[];
};

export class UIBuilder implements UIBuilderType {
  runtime: Runtime;
  cursor: Cursor;
  session: Session;
  page: Page;

  constructor(runtime: Runtime, session: Session, page: Page) {
    this.runtime = runtime;
    this.cursor = new Cursor();
    this.session = session;
    this.page = page;
  }

  markdown(content: string): void {
    markdown(this, content);
  }

  textInput(label: string, options: any = {}): string {
    return textInput(this, label, options);
  }

  numberInput(label: string, options: any = {}): number | null {
    return numberInput(this, label, options);
  }

  dateInput(label: string, options: any = {}): Date | null {
    return dateInput(this, label, options);
  }

  dateTimeInput(label: string, options: any = {}): Date | null {
    return dateTimeInput(this, label, options);
  }

  timeInput(label: string, options: any = {}): Date | null {
    return timeInput(this, label, options);
  }

  selectbox(label: string, options: any = {}): any {
    return selectbox(this, label, options);
  }

  multiSelect(label: string, options: any = {}): any {
    return multiSelect(this, label, options);
  }

  radio(label: string, options: any = {}): any {
    return radio(this, label, options);
  }

  checkbox(label: string, options: any = {}): boolean {
    return checkbox(this, label, options);
  }

  checkboxGroup(label: string, options: any = {}): any {
    return checkboxGroup(this, label, options);
  }

  textArea(label: string, options: any = {}): string {
    return textArea(this, label, options);
  }

  table(data: any, options: any = {}): any {
    return table(this, data, options);
  }

  button(label: string, options: any = {}): boolean {
    return button(this, label, options);
  }

  form(label: string, options: any = {}): [UIBuilderType, boolean] {
    return form(this, label, options);
  }

  columns(count: number, options: any = {}): UIBuilderType[] {
    return columns(this, count, options);
  }

  generatePageID(widgetType: any, path: number[]): string {
    if (!this.page) {
      return uuidNil;
    }
    const strPath = path.map((v) => v.toString()).join('_');
    return uuidv5(`${widgetType}-${strPath}`, this.page.id);
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

// Additional methods and logic can be added as needed.
