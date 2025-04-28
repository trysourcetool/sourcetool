---
sidebar_position: 15
---

# Form

`Form` groups several widgets, adds a submit button, and returns whether the user just pressed it.

## Signature

```go
formUI, submitted := ui.Form(buttonLabel string, opts ...form.Option)
```

* **`formUI`** – a *child* `UIBuilder` that you use to place the fields belonging to the form.
* **`submitted`** – `true` only on the execution immediately *after* the user presses the submit button. It resets to `false` on the next run.

## Option helpers

| Helper | Purpose | Default |
|--------|---------|---------|
| `form.WithButtonDisabled(true)` | Renders the submit button as disabled. | `false` |
| `form.WithClearOnSubmit(true)` | Clears all contained field values once the form is successfully submitted. | `false` |

## Behaviour notes

* Calling widgets **directly** on the parent builder after a `Form` call will render them *below* the form, not inside it. Always use `formUI` for form fields.
* If `WithClearOnSubmit` is disabled, each field keeps its previous value between submissions (session‑scoped state).
* The submit button label is the **first argument**; there is no separate option for it.

## Examples

### Minimal form

```go
func settingsPage(ui sourcetool.UIBuilder) error {
    formUI, submitted := ui.Form("Save")

    username := formUI.TextInput("Username")
    newsletter := formUI.Checkbox("Subscribe to newsletter")

    if submitted {
        return saveSettings(username, newsletter)
    }
    return nil
}
```

### Form that resets after submit

```go
formUI, submitted := ui.Form("Create", form.WithClearOnSubmit(true))
name  := formUI.TextInput("Name",  textinput.Required(true))
email := formUI.TextInput("Email", textinput.Required(true))

if submitted {
    _ = createUser(name, email)
}
```

### Disabled button until prerequisites met

```go
formUI, _ := ui.Form("Coming soon", form.WithButtonDisabled(true))
formUI.Markdown("Feature in beta – signup closed.")
```

---

### Related widgets

* [`TextInput`](./text-input) – single‑line text field.
* [`Checkbox`](./checkbox) – boolean toggle.
* [`Table`](./table) – displaying submitted data.