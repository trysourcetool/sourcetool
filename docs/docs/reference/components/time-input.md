---
sidebar_position: 10
---

# Time Input

`TimeInput` lets users pick a time of day using a clock‑style picker. It returns a `*time.Time` whose **date part is zeroed** (00‑01‑01) and whose location is whatever you pass in `WithLocation` (defaults to `time.Local`).

## Signature

```go
t := ui.TimeInput(label string, opts ...timeinput.Option) *time.Time
```

`t` is `nil` until the user chooses a value.

## Option helpers

| Helper | Purpose | Default |
|--------|---------|---------|
| `timeinput.WithPlaceholder("HH:mm")` | Placeholder when empty. | `""` |
| `timeinput.WithDefaultValue(time.Now())` | Pre‑fill on first render. | *nil* |
| `timeinput.WithRequired(true)` | Inside a [`Form`](./form) blocks submit until a value is chosen. | `false` |
| `timeinput.WithDisabled(true)` | Makes the field read‑only. | `false` |
| `timeinput.WithLocation(time.UTC)` | Time‑zone used for parsing and formatting. | `time.Local` |

## Behaviour notes

* **Storage** – backend saves the value in session; subsequent reruns keep the last input.
* **Format** – the widget always serialises as `HH:MM:SS` (Go’s `time.TimeOnly`). There is no custom format option.
* **Time‑zone** – if you care about absolute moments in time, apply the same `Location` when comparing with other timestamps.

## Examples

### Basic picker

```go
start := ui.TimeInput("Start time")
if start != nil {
    fmt.Println("selected:", start.Format(time.TimeOnly))
}
```

### Pre‑selected, disabled

```go
ui.TimeInput("Closing time",
    timeinput.WithDefaultValue(time.Date(0,0,0,18,0,0,0,time.Local)),
    timeinput.WithDisabled(true),
)
```

### Required inside a form

```go
f, submitted := ui.Form("Save")
when := f.TimeInput("Reminder", timeinput.WithRequired(true))
if submitted {
    scheduleReminder(*when)
}
```

---

### Related widgets

* [`DateInput`](./date-input) – pick a calendar date.  
* [`DateTimeInput`](./date-time-input) – pick date **and** time.