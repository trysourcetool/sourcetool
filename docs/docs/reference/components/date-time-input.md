---
sidebar_position: 9
---

# Date Time Input

`DateTimeInput` combines a calendar and clock picker. It returns the full timestamp the user chose.

## Signature

```go
dt := ui.DateTimeInput(label string, opts ...datetimeinput.Option) *time.Time
```

`dt` is `nil` until the user selects a value.

## Option helpers

| Helper | Purpose | Default |
|--------|---------|---------|
| `datetimeinput.WithPlaceholder("YYYY‑MM‑DD HH:mm")` | Placeholder when empty | `""` |
| `datetimeinput.WithDefaultValue(time.Now())` | Pre‑fill on first render | *nil* |
| `datetimeinput.WithRequired(true)` | Require a value inside a [`Form`](./form) | `false` |
| `datetimeinput.WithDisabled(true)` | Read‑only field | `false` |
| `datetimeinput.WithFormat("MM/DD/YYYY HH:mm")` | Display / parse format (Moment‑style tokens) | `"YYYY/MM/DD HH:MM:SS"` |
| `datetimeinput.WithMaxValue(t)` | Latest selectable timestamp | *nil* |
| `datetimeinput.WithMinValue(t)` | Earliest selectable timestamp | *nil* |
| `datetimeinput.WithLocation(loc)` | Time‑zone for parsing & formatting | `time.Local` |

## Behaviour

* **State type** – stored as `time.Time`. Changing `WithFormat` only affects rendering, not stored value.
* **Validation** – `Required`, `Min/MaxValue` enforce constraints client‑side before rerun.
* **Time‑zone** – all conversions use the provided `Location`; this matters when comparing dates across zones.

## Examples

### Basic input

```go
ev := ui.DateTimeInput("Event time")
if ev != nil {
    fmt.Println("selected:", ev.Format(time.RFC3339))
}
```

### Custom format + placeholder

```go
ui.DateTimeInput("Start",
    datetimeinput.WithPlaceholder("MM/DD/YYYY HH:mm"),
    datetimeinput.WithFormat("MM/DD/YYYY HH:mm"),
)
```

### Range‑restricted picker

```go
now := time.Now()
ui.DateTimeInput("Deadline",
    datetimeinput.WithMinValue(now),
    datetimeinput.WithMaxValue(now.Add(24*time.Hour)),
)
```

### Pre‑selected, disabled

```go
ui.DateTimeInput("Created",
    datetimeinput.WithDefaultValue(time.Now()),
    datetimeinput.WithDisabled(true),
)
```

---

### Related widgets

* [`DateInput`](./date-input) – date only.  
* [`TimeInput`](./time-input) – time only.