---
sidebar_position: 7
---

# Multi Select

`MultiSelect` renders a dropdown that lets users pick **zero or more** items.

## Signature

```go
sel := ui.MultiSelect(label string, opts ...multiselect.Option) *multiselect.Value
```

### Return value

If the user has chosen at least one item `sel` is:

```go
type Value struct {
    Values  []string // selected option labels
    Indexes []int    // matching indexes
}
```

Otherwise it is `nil`.

## Option helpers

| Helper | Explanation | Default |
|--------|-------------|---------|
| `multiselect.WithOptions("A", "B", ...)` | **Required** – list of choices. | – |
| `multiselect.WithPlaceholder("Select …")` | Grey text when nothing chosen. | `""` |
| `multiselect.WithDefaultValue("B", "C")` | Pre‑selected values on first render. | none |
| `multiselect.WithRequired(true)` | Inside a [`Form`](./form) makes at least one choice mandatory. | `false` |
| `multiselect.WithDisabled(true)` | Greys out the control. | `false` |
| `multiselect.WithFormatFunc(func(v string,i int)string)` | Transform labels for display (e.g. uppercase) while keeping underlying values intact. | identity |

## Behaviour notes

* **State type** – session stores an `[]int32` of indexes. Changing the order of `WithOptions` between releases will break old sessions.
* **Nil vs. empty slice** – when the widget is rendered for the first time with nothing selected, `sel` is `nil` not `&Value{}`.
* `WithDefaultValue` matches by *label text*, not index. Make sure the strings exist in the options list.

## Examples

### Basic usage

```go
skills := ui.MultiSelect("Skills",
    multiselect.WithOptions("Go", "JS", "Rust"),
)
if skills != nil {
    log.Println("chosen skills:", skills.Values)
}
```

### Placeholder & required

```go
ui.MultiSelect("Databases",
    multiselect.WithOptions("Postgres", "MySQL", "SQLite"),
    multiselect.WithPlaceholder("Pick at least one"),
    multiselect.WithRequired(true),
)
```

### Default values and custom formatting

```go
format := func(v string, _ int) string { return strings.ToUpper(v) }
ui.MultiSelect("Languages",
    multiselect.WithOptions("English", "Japanese", "Spanish"),
    multiselect.WithDefaultValue("Japanese"),
    multiselect.WithFormatFunc(format),
)
```

---

### Related widgets

* [`Select`](./select) – choose exactly one item.  
* [`CheckboxGroup`](./checkbox-group) – visible checklist alternative.