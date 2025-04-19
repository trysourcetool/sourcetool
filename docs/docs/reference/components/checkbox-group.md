---
sidebar_position: 13
---

# Checkbox Group

`CheckboxGroup` renders a list of checkboxes that share a single label. Users can tick **any number** of options; the selection is sent back on every **RerunPage** and stored in session state.

## Signature

```go
val := ui.CheckboxGroup(label string, opts ...checkboxgroup.Option) *checkboxgroup.Value
```

### Return value

If the user has selected at least one option, `val` is a pointer to:

```go
type Value struct {
    Values  []string // selected option labels (original, not formatted)
    Indexes []int    // corresponding indexes in the options slice
}
```
Otherwise it is `nil`.

## Option helpers

| Helper | Purpose |
|--------|---------|
| `checkboxgroup.WithOptions("A", "B", "C")` | Provide the list of choices. **Required.** |
| `checkboxgroup.WithDefaultValue("B", "C")` | Pre‑select items on first render. |
| `checkboxgroup.WithRequired(true)` | Prevent form submission unless at least one item is checked. |
| `checkboxgroup.WithDisabled(true)` | Grey out the whole group. |
| `checkboxgroup.WithFormatFunc(func(v string, i int) string { return strings.ToUpper(v) })` | Customize the *rendered* label without altering the underlying value. |

## Behaviour notes

* **State type** – the backend keeps an `[]int32` of indexes, so changing the order of `WithOptions` between deployments will break existing sessions.
* **Formatting** – `WithFormatFunc` is applied only to what the user sees; the `Values` returned remain the original strings.
* **Default vs. current** – after the first run, any user change overrides the default; subsequent reruns keep the user’s selection.

## Examples

### Basic group

```go
choice := ui.CheckboxGroup("Interests",
    checkboxgroup.WithOptions("Tech", "Science", "Art"),
)
if choice != nil {
    log.Printf("selected: %v", choice.Values)
}
```

### Default selection & required

```go
_ = ui.CheckboxGroup("Pick at least one pet",
    checkboxgroup.WithOptions("Dog", "Cat", "Bird"),
    checkboxgroup.WithDefaultValue("Dog"),
    checkboxgroup.WithRequired(true),
)
```

### Disabled group

```go
ui.CheckboxGroup("Archived list",
    checkboxgroup.WithOptions("A", "B"),
    checkboxgroup.WithDisabled(true),
)
```

### Custom label formatting

```go
format := func(v string, i int) string { return fmt.Sprintf("%d. %s", i+1, v) }
ui.CheckboxGroup("Steps",
    checkboxgroup.WithOptions("Build", "Test", "Deploy"),
    checkboxgroup.WithFormatFunc(format),
)
```

---

### Related widgets

* [`Checkbox`](./checkbox) – a single true/false toggle.
* [`Radio`](./radio) – choose exactly one option.
* [`MultiSelect`](./multi-select) – dropdown‑based multi‑choice.

