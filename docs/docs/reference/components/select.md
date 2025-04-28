---
sidebar_position: 6
---

# Select

`Selectbox` shows a dropdown where the user picks **exactly one** item.

## Signature

```go
sel := ui.Selectbox(label string, opts ...selectbox.Option) *selectbox.Value
```

`sel` is `nil` until the user makes a choice.

### Return struct

```go
type Value struct {
    Value string // chosen option label
    Index int    // index inside options slice
}
```

## Option helpers

| Helper | Purpose | Default |
|--------|---------|---------|
| `selectbox.WithOptions("A", "B")` | **Required** – list of dropdown options. | – |
| `selectbox.WithPlaceholder("Choose …")` | Grey hint when nothing selected. | `""` |
| `selectbox.WithDefaultValue("B")` | Pre‑select on first render. | *nil* |
| `selectbox.WithRequired(true)` | Inside a [`Form`](./form) blocks submit until a value is chosen. | `false` |
| `selectbox.WithDisabled(true)` | Renders read‑only dropdown. | `false` |
| `selectbox.WithFormatFunc(func(v string,i int)string)` | Modify display text (e.g. upper‑case) without changing stored value. | identity |

## Behaviour notes

* Backend stores the selected index (`int32`). Changing the order of `WithOptions` will break old sessions.
* `WithDefaultValue` matches by label string; ensure it exists in the options slice.
* When the widget is disabled the stored value remains unchanged.

## Examples

### Basic select

```go
country := ui.Selectbox("Country",
    selectbox.WithOptions("USA", "Canada", "Japan"),
)
if country != nil {
    log.Println("selected:", country.Value)
}
```

### Placeholder and required

```go
ui.Selectbox("Language",
    selectbox.WithOptions("Go", "Rust", "Python"),
    selectbox.WithPlaceholder("Pick one"),
    selectbox.WithRequired(true),
)
```

### Default value & custom formatting

```go
fmtOpt := func(v string, _ int) string { return strings.ToUpper(v) }
ui.Selectbox("Plan",
    selectbox.WithOptions("Free", "Pro", "Enterprise"),
    selectbox.WithDefaultValue("Pro"),
    selectbox.WithFormatFunc(fmtOpt),
)
```

---

### Related widgets

* [`MultiSelect`](./multi-select) – choose many items.  
* [`Radio`](./radio) – visible list of mutually exclusive choices.