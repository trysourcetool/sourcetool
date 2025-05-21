import { UIBuilder } from './uibuilder';
import { RouterInterface } from './router';
import { ButtonComponentOptions } from './uibuilder/widgets/button';
import { CheckboxComponentOptions } from './uibuilder/widgets/checkbox';
import { TextInputComponentOptions } from './uibuilder/widgets/textinput';
import { NumberInputComponentOptions } from './uibuilder/widgets/numberinput';
import { DateInputComponentOptions } from './uibuilder/widgets/dateinput';
import { DateTimeInputComponentOptions } from './uibuilder/widgets/datetimeinput';
import { TimeInputComponentOptions } from './uibuilder/widgets/timeinput';
import { RadioComponentOptions } from './uibuilder/widgets/radio';
import { SelectboxComponentOptions } from './uibuilder/widgets/selectbox';
import { MultiSelectComponentOptions } from './uibuilder/widgets/multiselect';
import { CheckboxGroupComponentOptions } from './uibuilder/widgets/checkboxgroup';
import { TextAreaComponentOptions } from './uibuilder/widgets/textarea';
import {
  TableComponentOptions,
  TableOnSelect,
  TableRowSelection,
} from './uibuilder/widgets/table';
import { FormComponentOptions } from './uibuilder/widgets/form';
import { ColumnsComponentOptions } from './uibuilder/widgets/columns';

import { Sourcetool, SourcetoolConfig } from './sourcetool';

// Export all components and types
export {
  // Enums
  TableOnSelect,
  TableRowSelection,

  // Sourcetool
  Sourcetool,
};

export type {
  // Component options
  ButtonComponentOptions,
  CheckboxComponentOptions,
  TextInputComponentOptions,
  NumberInputComponentOptions,
  DateInputComponentOptions,
  DateTimeInputComponentOptions,
  TimeInputComponentOptions,
  RadioComponentOptions,
  SelectboxComponentOptions,
  MultiSelectComponentOptions,
  CheckboxGroupComponentOptions,
  TextAreaComponentOptions,
  TableComponentOptions,
  FormComponentOptions,
  ColumnsComponentOptions,

  // Builder type
  UIBuilder,

  // Router
  RouterInterface,

  // Sourcetool
  SourcetoolConfig,
};
