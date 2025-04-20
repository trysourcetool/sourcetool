---
sidebar_position: 5
---

# Checkbox

`Checkbox` is a single true/false toggle. It is ideal for opt‑in terms\, feature flags\, or any binary decision.

## Signature

```go
checked := ui.Checkbox(label string, opts ...checkbox.Option) bool
```

The function returns the current checked state on *this* execution of the page.

## Option helpers

| Helper | Effect | Default |
|--------|--------|---------|
| `checkbox.WithDefaultValue(true)` | Sets the initial value for a **new** session. | `false` |
| `checkbox.WithRequired(true)` | Marks the field as required inside a [`Form`](./form); the form will not submit until checked. | `false` |
| `checkbox.WithDisabled(true)` | Renders the control as read‑only. | `false` |

## Behaviour

* The value is kept in session state (`bool`) – it does **not** reset between reruns like a Button.
* Changing the label does **not** affect stored data; only visual text.
* `WithDefaultValue` applies only if no previous value exists in the user’s session.

## Examples

### Simple checkbox

```go
subscribed := ui.Checkbox("Subscribe to newsletter")
if subscribed {
    sendWelcomeEmail(user)
}
```

### Pre‑checked & disabled

```go
ui.Checkbox("Beta feature enabled",
    checkbox.WithDefaultValue(true),
    checkbox.WithDisabled(true),
)
```

### Required checkbox inside a form

```go
form, submitted := ui.Form("Submit")
if submitted {
    // form validated, save data
}
agg := ui.Checkbox("I agree to the terms", checkbox.WithRequired(true))
if submitted && !agg {
    return errors.New("should never happen – form blocks submission")
}
```

---

### Related widgets

* [`CheckboxGroup`](./checkbox-group) – multi‑choice set of checkboxes.  
* [`Radio`](./radio) – choose exactly one option out of many.