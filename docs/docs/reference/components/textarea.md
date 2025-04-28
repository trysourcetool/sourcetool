---
sidebar_position: 2
---

# Text Area

`TextArea` is a multi‑line text box—perfect for comments, descriptions, or code snippets.

## Signature

```go
content := ui.TextArea(label string, opts ...textarea.Option) string
```

`content` is the current value; the empty string until the user types or a `DefaultValue` is supplied.

## Option helpers

| Helper | Purpose | Default |
|--------|---------|---------|
| `textarea.WithPlaceholder("Write here…")` | Grey hint text. | `""` |
| `textarea.WithDefaultValue("Lorem…")` | Pre‑fill on first render. | none |
| `textarea.WithRequired(true)` | Inside a [`Form`](./form) blocks submit if empty. | `false` |
| `textarea.WithDisabled(true)` | Read‑only field. | `false` |
| `textarea.WithMaxLength(1_000)` | Clamp characters. | unlimited |
| `textarea.WithMinLength(10)` | Enforce minimum length (client‑side). | 0 |
| `textarea.WithMaxLines(20)` | Vertical limit before scroll. | unlimited |
| `textarea.WithMinLines(2)` | Initial visible rows. | `2` |
| `textarea.WithAutoResize(false)` | Disable automatic height growth. | `true` |

## Behaviour notes

* **Session state** – value persists between reruns; use `WithDefaultValue` for first render only.
* **Auto‑resize** – when enabled (default) the control grows until it hits `MaxLines` (if set).
* **Validation** – `Required`, length, and line limits are enforced in the browser; re‑validate on the backend when necessary.

## Examples

### Basic comment box

```go
comment := ui.TextArea("Comment",
    textarea.WithPlaceholder("Leave your thoughts…"),
)
```

### Fixed‑height log viewer

```go
ui.TextArea("Logs",
    textarea.WithDefaultValue(loadLogs()),
    textarea.WithDisabled(true),
    textarea.WithAutoResize(false),
    textarea.WithMaxLines(15),
)
```

### Required feedback with length limits

```go
fb := ui.TextArea("Feedback",
    textarea.WithRequired(true),
    textarea.WithMinLength(20),
    textarea.WithMaxLength(500),
)
```

---

### Related widgets

* [`TextInput`](./text-input) – single‑line input.  
* [`Markdown`](./markdown) – display formatted text.