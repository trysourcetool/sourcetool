import type { WidgetState } from '@/store/modules/widgets/slice';
import type { WidgetJson } from '@trysourcetool/proto/widget/v1/widget';
import z from 'zod';

export const createWidgetState = (widget: WidgetJson): WidgetState | null => {
  if (!widget.id) {
    return null;
  }

  if (widget.textInput) {
    return {
      id: widget.id,
      type: 'textInput',
      value: widget.textInput.value ?? undefined,
      error: null,
    };
  }

  if (widget.numberInput) {
    return {
      id: widget.id,
      type: 'numberInput',
      value: Number.isFinite(widget.numberInput.value)
        ? widget.numberInput.value
        : undefined,
      error: null,
    };
  }

  if (widget.selectbox) {
    return {
      id: widget.id,
      type: 'selectbox',
      value: widget.selectbox.value ?? undefined,
      error: null,
    };
  }

  if (widget.checkbox) {
    return {
      id: widget.id,
      type: 'checkbox',
      value: widget.checkbox.value ?? undefined,
      error: null,
    };
  }

  if (widget.checkboxGroup) {
    return {
      id: widget.id,
      type: 'checkboxGroup',
      value: widget.checkboxGroup.value ?? undefined,
      error: null,
    };
  }

  if (widget.radio) {
    return {
      id: widget.id,
      type: 'radio',
      value: widget.radio.value ?? undefined,
      error: null,
    };
  }

  if (widget.multiSelect) {
    return {
      id: widget.id,
      type: 'multiSelect',
      value: widget.multiSelect.value ?? undefined,
      error: null,
    };
  }

  if (widget.dateInput) {
    return {
      id: widget.id,
      type: 'dateInput',
      value: widget.dateInput.value ?? undefined,
      error: null,
    };
  }

  if (widget.dateTimeInput) {
    return {
      id: widget.id,
      type: 'dateTimeInput',
      value: widget.dateTimeInput.value ?? undefined,
      error: null,
    };
  }

  if (widget.timeInput) {
    return {
      id: widget.id,
      type: 'timeInput',
      value: widget.timeInput.value ?? undefined,
      error: null,
    };
  }

  if (widget.textArea) {
    return {
      id: widget.id,
      type: 'textArea',
      value: widget.textArea.value ?? undefined,
      error: null,
    };
  }

  if (widget.table) {
    return {
      id: widget.id,
      type: 'table',
      value: widget.table.value ?? undefined,
      error: null,
    };
  }
  return null;
};
export const validateWidgetValue = <
  T extends WidgetJson,
  U extends WidgetState['type'],
  V extends WidgetState['value'],
>(
  widget: T,
  widgetType: U,
  value: V,
) => {
  // ==============================
  // textInput
  if (widget.textInput && widgetType === 'textInput') {
    let schema = z.string();
    if (widget.textInput.required) {
      schema = schema.min(1, {
        message: 'This field is required',
      });
    }
    if (widget.textInput.minLength) {
      schema = schema.min(widget.textInput.minLength, {
        message: `Min length is ${widget.textInput.minLength}`,
      });
    }
    if (widget.textInput.maxLength) {
      schema = schema.max(widget.textInput.maxLength, {
        message: `Max length is ${widget.textInput.maxLength}`,
      });
    }
    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }
  // ==============================
  // numberInput
  if (widget.numberInput && widgetType === 'numberInput') {
    let schema = z
      .number()
      .optional()
      .superRefine((value, ctx) => {
        if (widget.numberInput?.required && !value) {
          ctx.addIssue({
            code: 'custom',
            message: 'This field is required',
          });
        }

        const minValue = widget.numberInput?.minValue as number;
        if (minValue && Number.isFinite(minValue)) {
          if (value && value < minValue) {
            ctx.addIssue({
              code: 'custom',
              message: `Min is ${minValue}`,
            });
          }
        }

        const maxValue = widget.numberInput?.maxValue as number;
        if (maxValue && Number.isFinite(maxValue)) {
          if (value && value > maxValue) {
            ctx.addIssue({
              code: 'custom',
              message: `Max is ${maxValue}`,
            });
          }
        }
      });

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // dateInput
  if (widget.dateInput && widgetType === 'dateInput') {
    const schema = z.string().superRefine((value, ctx) => {
      if (widget.dateInput?.required && !value) {
        ctx.addIssue({
          code: 'custom',
          message: 'This field is required',
        });
      }

      if (widget.dateInput?.minValue) {
        const minDate = new Date(widget.dateInput.minValue);
        const valueDate = new Date(value);
        if (valueDate < minDate) {
          ctx.addIssue({
            code: 'custom',
            message: `Min is ${widget.dateInput.minValue}`,
          });
        }
      }

      if (widget.dateInput?.maxValue) {
        const maxDate = new Date(widget.dateInput.maxValue);
        const valueDate = new Date(value);
        if (valueDate > maxDate) {
          ctx.addIssue({
            code: 'custom',
            message: `Max is ${widget.dateInput.maxValue}`,
          });
        }
      }
    });

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // dateTimeInput

  if (widget.dateTimeInput && widgetType === 'dateTimeInput') {
    const schema = z.string().superRefine((value, ctx) => {
      if (widget.dateTimeInput?.required && !value) {
        ctx.addIssue({
          code: 'custom',
          message: 'This field is required',
        });
      }

      if (widget.dateTimeInput?.minValue) {
        const minDate = new Date(widget.dateTimeInput.minValue);
        const valueDate = new Date(value);
        if (valueDate < minDate) {
          ctx.addIssue({
            code: 'custom',
            message: `Min is ${widget.dateTimeInput.minValue}`,
          });
        }
      }

      if (widget.dateTimeInput?.maxValue) {
        const maxDate = new Date(widget.dateTimeInput.maxValue);
        const valueDate = new Date(value);
        if (valueDate > maxDate) {
          ctx.addIssue({
            code: 'custom',
            message: `Max is ${widget.dateTimeInput.maxValue}`,
          });
        }
      }
    });

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // timeInput

  if (widget.timeInput && widgetType === 'timeInput') {
    let schema = z.string().superRefine((value, ctx) => {
      if (widget.timeInput?.required && !value) {
        ctx.addIssue({
          code: 'custom',
          message: 'This field is required',
        });
      }
    });

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // textArea
  if (widget.textArea && widgetType === 'textArea') {
    let schema = z.string();
    if (widget.textArea.required) {
      schema = schema.min(1, {
        message: 'This field is required',
      });
    }
    if (widget.textArea.minLength) {
      schema = schema.min(widget.textArea.minLength, {
        message: `Min length is ${widget.textArea.minLength}`,
      });
    }
    if (widget.textArea.maxLength) {
      schema = schema.max(widget.textArea.maxLength, {
        message: `Max length is ${widget.textArea.maxLength}`,
      });
    }
    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // checkbox
  if (widget.checkbox && widgetType === 'checkbox') {
    let schema = z
      .boolean()
      .optional()
      .refine((value) => (widget.checkbox?.required ? value === true : true), {
        message: 'This field is required',
      });

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // checkboxGroup
  if (widget.checkboxGroup && widgetType === 'checkboxGroup') {
    let schema = z
      .array(z.string())
      .refine(
        (value) => (widget.checkboxGroup?.required ? value.length > 0 : true),
        {
          message: 'This field is required',
        },
      );

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // radio
  if (widget.radio && widgetType === 'radio') {
    let schema = z
      .number()
      .optional()
      .refine((value) => (widget.radio?.required ? !value : true), {
        message: 'This field is required',
      });

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // selectbox
  if (widget.selectbox && widgetType === 'selectbox') {
    let schema = z
      .number()
      .optional()
      .refine(
        (value) => (widget.selectbox?.required ? value !== undefined : true),
        {
          message: 'This field is required',
        },
      );

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  // ==============================
  // multiSelect
  if (widget.multiSelect && widgetType === 'multiSelect') {
    let schema = z
      .array(z.number().optional())
      .refine(
        (value) => (widget.multiSelect?.required ? value.length > 0 : true),
        {
          message: 'This field is required',
        },
      );

    return {
      success: schema.safeParse(value).success,
      error: schema.safeParse(value).error?.issues?.[0]?.message || null,
    };
  }

  return {
    success: true,
    error: null,
  };
};
