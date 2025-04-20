---
sidebar_position: 8
---

# Date Input

`DateInput` shows a calendar picker and returns the selected date.

## Signature

```go
picked := ui.DateInput(label string, opts ...dateinput.Option) *time.Time
```

`picked` is `nil` until the user chooses a date.

## Option helpers

| Helper | Purpose | Default |
|--------|---------|---------|
| `dateinput.WithPlaceholder("YYYY‑MM‑DD")` | Placeholder when no value selected | `""` |
| `dateinput.WithDefaultValue(t)` | Pre‑fill with a date | *nil* |
| `dateinput.WithRequired(true)` | Mark field as required inside a `Form` | `false` |
| `dateinput.WithDisabled(true)` | Render read‑only | `false` |
| `dateinput.WithFormat("MM/DD/YYYY")` | Custom display/layout format | `"YYYY/MM/DD"` |
| `dateinput.WithMaxValue(t)` | Latest selectable date | *nil* |
| `dateinput.WithMinValue(t)` | Earliest selectable date | *nil* |
| `dateinput.WithLocation(loc)` | Time‑zone for parsing/formatting | `time.Local` |

### Format string

The format uses [Moment.js‑style] tokens (e.g. `YYYY`, `MM`, `DD`). It does **not** use Go’s reference date layout.

## Behaviour notes

* State is stored as `time.Time` (date‑only part preserved). Changing the format does not affect stored values.
* Validation (`Required`, `Min/MaxValue`) happens client‑side before page rerun.
* `WithLocation` matters if you later compare the value with time‑zoned dates.

## Examples

### Basic input

```go
birthday := ui.DateInput("Birthday")
if birthday != nil {
    log.Println("Selected:", birthday.Format(time.DateOnly))
}
```

### Custom format & placeholder

```go
ui.DateInput("Start",
    dateinput.WithPlaceholder("MM/DD/YYYY"),
    dateinput.WithFormat("MM/DD/YYYY"),
)
```

### Date range constraint

```go
now := time.Now()
ui.DateInput("Deadline",
    dateinput.WithMinValue(now),
    dateinput.WithMaxValue(now.AddDate(0, 0, 30)),
)
```

### Pre‑selected & disabled

```go
ui.DateInput("Creation date",
    dateinput.WithDefaultValue(time.Now()),
    dateinput.WithDisabled(true),
)
```

---

### Related widgets

* [`DateTimeInput`](./date-time-input) – pick date **and** time.  
* [`TimeInput`](./time-input) – pick time only.