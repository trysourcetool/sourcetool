// Common options interfaces for UI components

export type ButtonInternalOptions = {
  label: string;
  disabled: boolean;
};

export type CheckboxInternalOptions = {
  label: string;
  defaultValue: boolean;
  required: boolean;
  disabled: boolean;
};

export type CheckboxGroupInternalOptions = {
  label: string;
  options: string[];
  defaultValue: string[] | null;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type ColumnsInternalOptions = {
  columns: number;
  weight?: number[];
};

export type DateInputInternalOptions = {
  label: string;
  placeholder: string;
  defaultValue: Date | null;
  required: boolean;
  disabled: boolean;
  format: string;
  maxValue: Date | null;
  minValue: Date | null;
  location: string; // Timezone location
};

export type DateTimeInputInternalOptions = {
  label: string;
  placeholder: string;
  defaultValue: Date | null;
  required: boolean;
  disabled: boolean;
  format: string;
  maxValue: Date | null;
  minValue: Date | null;
  location: string; // Timezone location
};

export type FormInternalOptions = {
  buttonLabel: string;
  buttonDisabled: boolean;
  clearOnSubmit: boolean;
};

export type MarkdownInternalOptions = {
  body: string;
};

export type MultiSelectInternalOptions = {
  label: string;
  options: string[];
  defaultValue: string[] | null;
  placeholder: string;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type NumberInputInternalOptions = {
  label: string;
  placeholder: string;
  defaultValue: number | null;
  required: boolean;
  disabled: boolean;
  maxValue: number | null;
  minValue: number | null;
};

export type RadioInternalOptions = {
  label: string;
  options: string[];
  defaultValue: string | null;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type SelectboxInternalOptions = {
  label: string;
  options: string[];
  defaultValue: string | null;
  placeholder: string;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type TableInternalOptions = {
  header: string;
  description: string;
  height: number | null;
  columnOrder: string[];
  onSelect: string;
  rowSelection: string;
};

export type TextAreaInternalOptions = {
  label: string;
  placeholder: string;
  defaultValue: string | null;
  required: boolean;
  disabled: boolean;
  maxLength: number | null;
  minLength: number | null;
  maxLines: number | null;
  minLines: number | null;
  autoResize: boolean;
};

export type TextInputInternalOptions = {
  label: string;
  placeholder: string;
  defaultValue: string | null;
  required: boolean;
  disabled: boolean;
  maxLength: number | null;
  minLength: number | null;
};

export type TimeInputInternalOptions = {
  label: string;
  placeholder: string;
  defaultValue: Date | null;
  required: boolean;
  disabled: boolean;
  location: string; // Timezone location
};
