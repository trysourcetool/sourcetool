---
sidebar_position: 4
---

# Number Input

`NumberInput` captures a numeric value (floating‑point). It supports min/max validation and placeholder text.

## Signature

```go
val := ui.NumberInput(label string, opts ...numberinput.Option) *float64
```

`val` points to the current number. If the user has not entered anything yet it defaults to **0.0**.

## Option helpers

| Helper | Description | Default |
|--------|-------------|---------|
| `numberinput.WithPlaceholder("0–100")` | Grey text hint when the field is empty | `""` |
| `numberinput.WithDefaultValue(42.5)` | Initial value for a new session | `0.0` |
| `numberinput.WithRequired(true)` | Inside a [`Form`](./form) the submit button is blocked until a value is entered | `false` |
| `numberinput.WithDisabled(true)` | Renders read‑only | `false` |
| `numberinput.WithMaxValue(100)` | Upper limit; client prevents larger input | *none* |
| `numberinput.WithMinValue(0)` | Lower limit; client prevents smaller input | *none* |

## Behaviour notes

* The field accepts both integers and decimals – it is parsed as `float64`.
* `MinValue` and `MaxValue` are enforced client‑side; back‑end should still validate if data integrity is critical.
* The value persists in the session, so rerunning the page keeps the last input.

## Examples

### Simple number

```go
age := ui.NumberInput("Age")
if age != nil {
    log.Println("age:", *age)
}
```

### With range and placeholder

```go
score := ui.NumberInput("Score",
    numberinput.WithPlaceholder("0 – 100"),
    numberinput.WithMinValue(0),
    numberinput.WithMaxValue(100),
)
```

### Required within a form

```go
formUI, submitted := ui.Form("Save")
price := formUI.NumberInput("Price", numberinput.WithRequired(true))

if submitted {
    fmt.Printf("Price set to %.2f\n", *price)
}
```

---

### Related widgets

* [`TextInput`](./text-input) – unrestricted text field.