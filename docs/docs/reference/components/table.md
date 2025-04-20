---
sidebar_position: 11
---

# Table

`Table` renders arbitrary slice/array data in a scrollable grid. It supports column re‑ordering, paging (via `Height`) and row selection.

## Signature

```go
result := ui.Table(data any, opts ...table.Option) table.Value
```

* **`data`** can be a slice of structs, maps, or any value encodable by `encoding/json`.
* **`result`** contains the user’s selection (if any).

### Return type

```go
// package table

type Value struct {
    Selection *Selection // nil if no row selected
}

type Selection struct {
    Row  int   // first selected row (for single‑mode)
    Rows []int // all selected rows (for multiple‑mode)
}
```

## Option helpers

| Helper | Purpose | Default |
|--------|---------|---------|
| `table.WithHeader("Users")` | Title above the table. | empty |
| `table.WithDescription("Active accounts")` | Text below the header. | empty |
| `table.WithHeight(10)` | Visible rows before the grid scrolls. | auto (all rows) |
| `table.WithColumnOrder("ID", "Name", "Email")` | Re‑arrange columns by field name / map key. | natural order |
| `table.WithOnSelect(table.OnSelectRerun)` | Behaviour when a row is clicked: `OnSelectRerun` = rerun page; `OnSelectIgnore` = do nothing. | `OnSelectIgnore` |
| `table.WithRowSelection(table.RowSelectionMultiple)` | Selection mode: `Single` or `Multiple`. | `Single` |

## Behaviour notes

* **Data encoding** – the builder marshals `data` to JSON; unsupported types will panic. Make sure the slice elements are serialisable.
* **Selection persistence** – `Selection` is stored in the session. Changing `RowSelection` between runs resets it.
* **No built‑in sorting** – client handles sorting; it sends the same `data` back, so keep row order deterministic for correct index mapping.

## Examples

### Basic users table

```go
type User struct {
    ID    string
    Name  string
    Email string
}

users := []User{...}
res := ui.Table(users,
    table.WithHeader("Users"),
    table.WithHeight(15),
)
if res.Selection != nil {
    fmt.Println("Clicked row:", res.Selection.Row)
}
```

### Multiple selection & rerun

```go
res := ui.Table(items,
    table.WithRowSelection(table.RowSelectionMultiple),
    table.WithOnSelect(table.OnSelectRerun),
)
if res.Selection != nil {
    fmt.Println("Rows:", res.Selection.Rows)
}
```

---

### Related widgets

* [`Form`](./form) – embed a table inside a larger data‑entry flow.
* [`TextInput`](./text-input) – implement client‑side filters.