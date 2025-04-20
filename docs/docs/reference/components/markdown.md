---
sidebar_position: 12
---

# Markdown

`Markdown` renders static content written in GitHub‑flavoured Markdown. Use it for headings, lists, links, call‑outs, or any other rich text that doesn’t require user interaction.

## Signature

```go
ui.Markdown(body string)
```

There are **no** option helpers and the call does not return a value.

## Behaviour notes

* Markdown is treated as *read‑only*: it never triggers a page rerun by itself.
* The widget supports the standard CommonMark feature set plus GitHub extensions (tables, task lists, strikethrough, etc.).
* Each call advances the builder cursor, so subsequent widgets render **below** the markdown block.

## Examples

### Basic usage

```go
ui.Markdown(`# Welcome\nThis is **Sourcetool**!`)
```

### Multi‑line content

```go
ui.Markdown(`
## Getting Started

- Install the CLI: ` + "`brew install sourcetool`" + `
- Sign in with your API key
- Run \\`sourcetool init\\` in your project root
`)
```

---

### Related widgets

* [`TextArea`](./textarea) – editable multi‑line text.  
* [`Form`](./form) – collect user input and submit.