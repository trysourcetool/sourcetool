---
sidebar_position: 16
---

# Columns

`Columns` lets you split the page into *n* vertical blocks and build each block with its own `UIBuilder`.

## Signature

```go
cols := ui.Columns(count int, opts ...columns.Option) []sourcetool.UIBuilder
```

`cols[i]` is a **child builder** whose widgets render in column *i*.

## Option helpers

| Helper | Use | Notes |
|--------|-----|-------|
| `columns.WithWeight(3, 1)` | Set relative widths. Must pass exactly *count* positive ints. | If omitted or invalid the SDK falls back to equal widths.

## Behaviour

* Weight numbers are normalised to `float64` so `WithWeight(2, 1, 1)` results in column weights `0.5, 0.25, 0.25`.
* Passing `count ≤ 0` returns `nil` and renders nothing.
* The parent builder’s cursor advances by **one** after the call, so subsequent widgets appear *below* the column row.

## Examples

### Two equal columns

```go
cols := ui.Columns(2)
cols[0].Markdown("### Left")
cols[1].Markdown("### Right")
```

### 3‑column 2:1:1 layout

```go
cols := ui.Columns(3, columns.WithWeight(2, 1, 1))
cols[0].Markdown("## Main")
cols[1].Markdown("Sidebar A")
cols[2].Markdown("Sidebar B")
```

### Responsive form section

```go
field := ui.Columns(2)
name  := field[0].TextInput("Name")
email := field[1].TextInput("Email")

submit := ui.Button("Save")
if submit {
    save(name, email)
}
```

---

### Related widgets

* [`Form`](./form) – higher‑level container with validation & submit button.