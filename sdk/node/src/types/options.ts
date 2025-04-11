// Common options interfaces for UI components

export type ButtonOptions = {
  label: string;
  disabled: boolean;
};

export type CheckboxOptions = {
  label: string;
  defaultValue: boolean;
  required: boolean;
  disabled: boolean;
};

export type CheckboxGroupOptions = {
  label: string;
  options: string[];
  defaultValue: string[] | null;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type ColumnsOptions = {
  columns: number;
  weight?: number[];
};

export type DateInputOptions = {
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

export type DateTimeInputOptions = {
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

export type FormOptions = {
  buttonLabel: string;
  buttonDisabled: boolean;
  clearOnSubmit: boolean;
};

export type MarkdownOptions = {
  body: string;
};

export type MultiSelectOptions = {
  label: string;
  options: string[];
  defaultValue: string[] | null;
  placeholder: string;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type NumberInputOptions = {
  label: string;
  placeholder: string;
  defaultValue: number | null;
  required: boolean;
  disabled: boolean;
  maxValue: number | null;
  minValue: number | null;
};

export type RadioOptions = {
  label: string;
  options: string[];
  defaultValue: string | null;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type SelectboxOptions = {
  label: string;
  options: string[];
  defaultValue: string | null;
  placeholder: string;
  required: boolean;
  disabled: boolean;
  formatFunc?: (value: string, index: number) => string;
};

export type TableOptions = {
  header: string;
  description: string;
  height: number | null;
  columnOrder: string[];
  onSelect: string;
  rowSelection: string;
};

export type TextAreaOptions = {
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

export type TextInputOptions = {
  label: string;
  placeholder: string;
  defaultValue: string | null;
  required: boolean;
  disabled: boolean;
  maxLength: number | null;
  minLength: number | null;
};

export type TimeInputOptions = {
  label: string;
  placeholder: string;
  defaultValue: Date | null;
  required: boolean;
  disabled: boolean;
  location: string; // Timezone location
};
