# Sourcetool Protocol Buffers

> **Note:** This project now uses a consolidated setup with Docker Compose and a root Makefile.
> See the [root README.md](../README.md) for instructions on how to start the entire application.

## Overview

Protocol Buffers definitions for Sourcetool. This directory contains the Protocol Buffers schema definitions that are used to generate code for both the backend (Go) and frontend (TypeScript).

## Directory Structure

- `/proto` - Protocol Buffers schema definitions
- `/go` - Generated Go code
- `/ts` - Generated TypeScript code

## Usage

### Go

Add this module as a dependency:

```bash
go get github.com/trysourcetool/sourcetool/proto
```

Then import the generated packages:

```go
import (
    commonv1 "github.com/trysourcetool/sourcetool/proto/go/common/v1"
    exceptionv1 "github.com/trysourcetool/sourcetool/proto/go/exception/v1"
    pagev1 "github.com/trysourcetool/sourcetool/proto/go/page/v1"
    websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
    widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"
)
```

### TypeScript/JavaScript

Add the package as a dependency:

```bash
npm install @trysourcetool/proto
# or
yarn add @trysourcetool/proto
```

Then import and use the generated code:

```typescript
import { Message } from '@trysourcetool/proto/websocket/v1/message';
import { Page } from '@trysourcetool/proto/page/v1/page';
import { Widget } from '@trysourcetool/proto/widget/v1/widget';
import { Exception } from '@trysourcetool/proto/exception/v1/exception';
import { Common } from '@trysourcetool/proto/common/v1/common';
```