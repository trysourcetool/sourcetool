---
sidebar_position: 1
---

# Quickstart

## Why Sourcetool?

Sourcetool gives you a **React‑quality UI** while writing **pure Go**. Widgets are declarative, live‑update over WebSockets, and are fully type‑safe. That means you stay in one language and one repo, yet ship internal tools that feel indistinguishable from bespoke front‑ends.

## Installation

```bash
go get github.com/trysourcetool/sourcetool-go
```

### Prerequisites

1. Create a project and copy the API key from the **Sourcetool Dashboard**.
2. Go 1.21 or newer.

## Hello, users page

```go
package main

import (
    "log"
    "time"

    "github.com/trysourcetool/sourcetool-go"
    "github.com/trysourcetool/sourcetool-go/textinput"
    "github.com/trysourcetool/sourcetool-go/table"
)

type User struct {
    ID    string
    Name  string
    Email string
    Joined time.Time
}

func listUsers(filter string) ([]User, error) {
    return []User{{ID: "1", Name: "Jane", Email: "jane@example.com", Joined: time.Now()}}, nil
}

func usersPage(ui sourcetool.UIBuilder) error {
    ui.Markdown("## Users")

    name := ui.TextInput("Name", textinput.WithPlaceholder("filter by name"))

    users, err := listUsers(name)
    if err != nil {
        return err
    }

    ui.Table(users,
        table.WithHeader("Users list"),
        table.WithOnSelect(table.OnSelectRerun),
    )
    return nil
}

func main() {
    st := sourcetool.New(&sourcetool.Config{ // New now takes a *Config
        APIKey:   "YOUR_API_KEY",
        Endpoint: "wss://api.trysourcetool.com", // bare host is fine; SDK appends /ws
    })

    st.Page("/users", "Users", usersPage)

    if err := st.Listen(); err != nil {
        log.Fatal(err)
    }
}
```

### What happened?

* `sourcetool.New(&sourcetool.Config{…})` creates a host process and connects to the cloud endpoint.
* The `usersPage` function runs **every time** the client opens `/users`, or when a widget (e.g. table row) asks for a rerun.
* Widgets such as `TextInput` or `Table` persist their state in the session between reruns.

## Building forms

```go
func createUserPage(ui sourcetool.UIBuilder) error {
    formUI, submitted := ui.Form("Create", form.WithClearOnSubmit(true))

    name  := formUI.TextInput("Name",  textinput.WithRequired(true))
    email := formUI.TextInput("Email", textinput.WithRequired(true))

    role := formUI.Selectbox("Role",
        selectbox.WithOptions("Admin", "User", "Guest"),
        selectbox.WithRequired(true),
    )

    if submitted {
        _ = createUser(name, email, role.Value)
    }
    return nil
}
```

## Component catalogue

| Group | Components |
|-------|------------|
| Input | **TextInput**, TextArea, NumberInput, DateInput, DateTimeInput, TimeInput |
| Selection | Selectbox, MultiSelect, Radio, Checkbox, CheckboxGroup |
| Layout | Columns, Form, Table |
| Display | Markdown |
| Action | Button |

Each widget offers `With…` option helpers for behaviour—see the individual reference pages.

## Next steps

1. Dive into the [concepts](../concepts/pages) to understand sessions, reruns, and access control.  
2. Browse the complete [widget reference](../reference/components).  
3. Explore examples in the main repo for advanced patterns such as auth middleware and streaming logs.

Need help? Join the community on Discord or open a GitHub discussion.