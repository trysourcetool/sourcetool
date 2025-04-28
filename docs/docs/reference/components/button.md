---
sidebar_position: 3
---

# Button

`Button` renders a clickable button. When the user presses it the client sends a **RerunPage** message, the page is re‑executed, and the call returns `true` for that single run. On the next run the value automatically resets to `false`.

## Signature

```go
pressed := ui.Button(label string, opts ...button.Option) bool
```

## Return value

| Type | Meaning |
|------|---------|
| `bool` | `true` if the button was clicked since the previous execution of the page; otherwise `false`. |

## Options

| Helper | Description |
|--------|-------------|
| `button.WithDisabled(true)` | Renders the button in a disabled (non‑clickable) state. |

## Examples

### Basic button

```go
if ui.Button("Refresh") {
    // refresh data
}
```

### Disabled button

```go
ui.Button("Cannot Click", button.WithDisabled(true))
```

### Using the returned value

```go
clicked := ui.Button("Create")
if clicked {
    if err := createResource(); err != nil {
        ui.Markdown("⚠️ failed to create resource")
        return err
    }
}
```