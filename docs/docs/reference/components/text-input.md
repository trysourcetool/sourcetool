---
sidebar_position: 1
---

# Text Input

`TextInput` renders a single‑line edit box and returns the current text.

## Signature

```go
val := ui.TextInput(label string, opts ...textinput.Option) string
```

`val` is an empty string until the user types or a `DefaultValue` is supplied.

## Option helpers

| Helper | Use | Default |
|--------|-----|---------|
| `textinput.WithPlaceholder("Your name")` | Hint text when empty | `""` |
| `textinput.WithDefaultValue("Alice")` | Pre‑fill on first render | none |
| `textinput.WithRequired(true)` | Inside a [`Form`](./form) blocks submit if blank | `false` |
| `textinput.WithDisabled(true)` | Read‑only field | `false` |
| `textinput.WithMaxLength(64)` | Upper character limit | unlimited |
| `textinput.WithMinLength(3)` | Lower character limit (client‑side) | 0 |

## Behaviour notes

* The value is stored in session state, so it persists across page reruns.
* **Validation** – `MinLength`, `MaxLength`, and `Required` are enforced in the browser; always re‑validate on the backend if critical.
* Changing the label does not clear or change the stored value.

## Examples

### Basic input with placeholder

```go
name := ui.TextInput("Name", textinput.WithPlaceholder("Enter your name"))
```

### Required field inside a form

```go
f, submitted := ui.Form("Save")
email := f.TextInput("Email", textinput.WithRequired(true))
if submitted {
    log.Printf("email: %s", email)
}
```

### Length‑constrained password

```go
pwd := ui.TextInput("Password",
    textinput.WithPlaceholder("8–64 chars"),
    textinput.WithMinLength(8),
    textinput.WithMaxLength(64),
)
```

### Disabled, pre‑filled key

```go
ui.TextInput("API Key",
    textinput.WithDefaultValue(os.Getenv("API_KEY")),
    textinput.WithDisabled(true),
)
```

---

### Related widgets

* [`TextArea`](./textarea) – multi‑line text.  
* [`NumberInput`](./number-input) – numeric entry.