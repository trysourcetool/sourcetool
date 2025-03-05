---
sidebar_position: 4
---

# Number Input

The Number Input widget provides a specialized input field for numerical values with optional validation and formatting.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `double` | `None` | Current numeric value of the input |
| `label` | `string` | `""` | Label text displayed above the input |
| `placeholder` | `string` | `""` | Placeholder text displayed when the input is empty |
| `default_value` | `double` | `None` | Initial value of the input |
| `required` | `bool` | `False` | Whether the input is required |
| `disabled` | `bool` | `False` | Whether the input is disabled |
| `max_value` | `double` | `None` | Maximum allowed value |
| `min_value` | `double` | `None` | Minimum allowed value |

## Event Handling

The Number Input widget emits events when its value changes:

```go
// Define a number input with change handler
numberInput := widget.NewNumberInput(ctx, widget.NumberInputOptions{
    Label: "Quantity",
    DefaultValue: 1,
    Min: 0,
    Max: 100,
    Step: 1,
    OnChange: func(value float64) {
        // Handle number input value change
        fmt.Printf("Number input value changed: %v\n", value)
    },
})
```

## Examples

### Basic Number Input

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic number input
    numberInput := widget.NewNumberInput(ctx, widget.NumberInputOptions{
        Label: "Age",
        Placeholder: "Enter your age",
        Min: 0,
        Max: 120,
    })
    
    // Add the number input to your UI
    container.Add(numberInput)
}
```

### Decimal Number Input

```go
// Create a number input for decimal values
decimalInput := widget.NewNumberInput(ctx, widget.NumberInputOptions{
    Label: "Price",
    DefaultValue: 9.99,
    Min: 0,
    Step: 0.01,
    Precision: 2,
    Format: "$%.2f", // Display as currency
    OnChange: func(value float64) {
        fmt.Printf("Price updated: $%.2f\n", value)
    },
})
```

### Disabled Number Input

```go
// Create a disabled number input
disabledInput := widget.NewNumberInput(ctx, widget.NumberInputOptions{
    Label: "Score (Read Only)",
    DefaultValue: 85,
    Disabled: true,
    ShowControls: false,
})
```

## Best Practices

1. Always provide a meaningful label and placeholder text to guide users
2. Set appropriate `min` and `max` values to prevent invalid input
3. Choose a suitable `step` value based on the expected precision
4. Use the `format` property to display values in a user-friendly way (e.g., currency, percentages)
5. Consider disabling the controls for read-only values
6. Provide clear validation feedback for user input
7. Use appropriate precision for decimal values to avoid confusion

## Related Components

- [TextInput](./text-input) - For general text input
- `Slider` - For selecting a number from a range using a visual slider
- `CurrencyInput` - Specialized input for currency values with formatting
