---
sidebar_position: 14
---

# Radio

`Radio` renders a mutually‑exclusive choice set. Users can select **exactly one** option.

## Signature

```go
choice := ui.Radio(label string, opts ...radio.Option) *radio.Value
```

### Return value

```go
type Value struct {
    Value string // selected option label
    Index int    // zero‑based index in the options slice
}
```
If the user has never chosen an option and no default is set, `choice` is `nil`. citeturn14file19

## Option helpers

| Helper | Description | Default |
|--------|-------------|---------|
| `radio.WithOptions("A", "B", ...)` | Provide the list of choices. **Required.** | – |
| `radio.WithDefaultValue("B")` | Pre‑select an option on first render (match by label). | none |
| `radio.WithRequired(true)` | Inside a [`Form`](./form) the submit button is blocked until a choice is made. | `false` |
| `radio.WithDisabled(true)` | Greys out the entire group. | `false` |
| `radio.WithFormatFunc(func(v string,i int)string)` | Transform the *displayed* label (e.g. capitalise) while keeping the underlying value intact. | identity |

## Behaviour notes

* **State type** – backend stores a single `int32` index.
* `WithDefaultValue` must match one of the option strings; otherwise no default is applied.
* Changing the order of `WithOptions` between deployments will alter existing state.
* `FormatFunc` is applied every render; it cannot depend on external mutable state unless the page is rerun.

## Examples

### Basic radio

```go
color := ui.Radio("Primary colour",
    radio.WithOptions("Red", "Green", "Blue"),
)
if color != nil {
    fmt.Println("selected:", color.Value)
}
```

### Default & required

```go
ui.Radio("Plan",
    radio.WithOptions("Free", "Pro", "Enterprise"),
    radio.WithDefaultValue("Pro"),
    radio.WithRequired(true),
)
```

### Custom label formatting

```go
format := func(v string, _ int) string { return strings.ToUpper(v) }
ui.Radio("Priority",
    radio.WithOptions("low", "medium", "high"),
    radio.WithFormatFunc(format),
)
```

### Disabled group

```go
ui.Radio("Archived",
    radio.WithOptions("Yes", "No"),
    radio.WithDisabled(true),
)
```

---

### Related widgets

* [`Select`](./select) – dropdown single‑choice.  
* [`CheckboxGroup`](./checkbox-group) – multi‑choice checklist.