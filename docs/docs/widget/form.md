---
sidebar_position: 15
---

# Form

The Form widget provides a container for organizing form elements with built-in submission handling.

## Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `value` | `bool` | `False` | Current submission state of the form |
| `button_label` | `string` | `"Submit"` | Text displayed on the form's submit button |
| `button_disabled` | `bool` | `False` | Whether the submit button is disabled |
| `clear_on_submit` | `bool` | `False` | Whether to clear form inputs after submission |

## Event Handling

The Form widget emits events when submitted:

```go
// Define a form with submit handler
form := widget.NewForm(ctx, widget.FormOptions{
    ButtonLabel: "Save Changes",
    OnSubmit: func() {
        // Handle form submission
        fmt.Println("Form submitted!")
        // Process form data, validate inputs, etc.
    },
})
```

## Examples

### Basic Form

```go
package main

import (
    "fmt"
    "github.com/sourcetool/widget"
)

func main() {
    // Create a basic form
    form := widget.NewForm(ctx, widget.FormOptions{
        ButtonLabel: "Submit",
        OnSubmit: func() {
            fmt.Println("Form submitted!")
            // Process form data
        },
    })
    
    // Add form elements
    nameInput := widget.NewTextInput(ctx, widget.TextInputOptions{
        Label: "Name",
        Required: true,
    })
    
    emailInput := widget.NewTextInput(ctx, widget.TextInputOptions{
        Label: "Email",
        Required: true,
    })
    
    messageTextarea := widget.NewTextArea(ctx, widget.TextAreaOptions{
        Label: "Message",
        Required: true,
    })
    
    // Add inputs to the form
    form.Add(nameInput)
    form.Add(emailInput)
    form.Add(messageTextarea)
    
    // Add the form to your UI
    container.Add(form)
}
```

### Form with Validation

```go
// Create a form with validation
registrationForm := widget.NewForm(ctx, widget.FormOptions{
    ButtonLabel: "Register",
    OnSubmit: func() {
        // Validate form data
        if !isValidEmail(emailInput.GetValue()) {
            emailInput.SetError("Please enter a valid email address")
            return
        }
        
        if passwordInput.GetValue() != confirmPasswordInput.GetValue() {
            confirmPasswordInput.SetError("Passwords do not match")
            return
        }
        
        // If validation passes, submit the form
        fmt.Println("Registration form submitted!")
        // Process registration
    },
})

// Add form elements with validation
emailInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Email",
    Required: true,
})

passwordInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Password",
    Required: true,
    Type: "password",
    MinLength: 8,
})

confirmPasswordInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Confirm Password",
    Required: true,
    Type: "password",
})

// Add inputs to the form
registrationForm.Add(emailInput)
registrationForm.Add(passwordInput)
registrationForm.Add(confirmPasswordInput)
```

### Multi-Step Form

```go
// Create a multi-step form
checkoutForm := widget.NewForm(ctx, widget.FormOptions{
    ButtonLabel: "Continue",
    OnSubmit: func() {
        if currentStep < totalSteps {
            // Move to the next step
            currentStep++
            updateFormStep()
        } else {
            // Final submission
            fmt.Println("Checkout complete!")
            // Process order
        }
    },
})

// Step 1: Customer Information
customerInfoStep := widget.NewContainer(ctx)
nameInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Full Name",
    Required: true,
})
emailInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Email",
    Required: true,
})
customerInfoStep.Add(nameInput)
customerInfoStep.Add(emailInput)

// Step 2: Shipping Address
shippingAddressStep := widget.NewContainer(ctx)
addressInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Street Address",
    Required: true,
})
cityInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "City",
    Required: true,
})
shippingAddressStep.Add(addressInput)
shippingAddressStep.Add(cityInput)

// Step 3: Payment Information
paymentInfoStep := widget.NewContainer(ctx)
cardNumberInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Card Number",
    Required: true,
})
expiryDateInput := widget.NewTextInput(ctx, widget.TextInputOptions{
    Label: "Expiry Date",
    Required: true,
})
paymentInfoStep.Add(cardNumberInput)
paymentInfoStep.Add(expiryDateInput)

// Function to update the form based on the current step
func updateFormStep() {
    checkoutForm.Clear()
    
    if currentStep == 1 {
        checkoutForm.Add(customerInfoStep)
    } else if currentStep == 2 {
        checkoutForm.Add(shippingAddressStep)
    } else if currentStep == 3 {
        checkoutForm.Add(paymentInfoStep)
        checkoutForm.SetButtonLabel("Complete Order")
    }
}

// Initialize the form with the first step
currentStep = 1
totalSteps = 3
updateFormStep()
```

## Best Practices

1. Group related form elements together logically
2. Provide clear labels for all form inputs
3. Indicate required fields visually (e.g., with an asterisk)
4. Implement client-side validation for immediate feedback
5. Display validation errors clearly next to the relevant inputs
6. Use appropriate input types for different data (e.g., email, password, number)
7. Consider the tab order for keyboard navigation
8. Provide clear and descriptive button labels
9. Consider disabling the submit button until all required fields are filled
10. Provide feedback after form submission (success or error messages)

## Related Components

- [TextInput](./text-input) - For single-line text input
- [TextArea](./textarea) - For multi-line text input
- [Checkbox](./checkbox) - For boolean input
- [Select](./select) - For selecting from predefined options
- [DateInput](./date-input) - For date selection
