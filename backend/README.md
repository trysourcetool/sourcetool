# Sourcetool Backend API

A Go-based backend service for the Sourcetool project.

> **Note:** This project now uses a consolidated setup with Docker Compose and a root Makefile.
> See the [root README.md](../README.md) for instructions on how to start the entire application.

## Overview

The backend provides the API for the Sourcetool application, handling:
- User authentication and authorization
- Organization management
- Group management
- Page management
- Environment management
- API key management
- WebSocket connections

## Technology

- Go 1.24
- PostgreSQL 15
- Redis 7
- Docker

## Directory Structure

```
backend/
├── cmd/                # entry points
│   ├── server/         # main.go – wires everything + graceful shutdown
│   └── internal/       # helpers only needed by binaries (redis/smtp/upgrader/fixtures)
├── devtools/           # CLI & scripts (migrations, etc.)
├── internal/           # all import‑able code lives here (go’s internal visibility)
│   ├── config/         # env parsing, URL helpers
│   ├── core/           # *domain* models – NO I/O
│   ├── postgres/       # SQL repos (`*_query.go`) + db logger
│   ├── pubsub/         # Redis adapter (fan‑out)
│   ├── mail/           # SMTP adapter
│   ├── ws/             # WebSocket manager / ping loop
│   ├── server/         # HTTP & WS handlers, DTOs, middleware, CORS (CE/EE split via build‑tags)
│   ├── jwt/            # generic JWT helpers & claim structs
│   ├── google/         # OAuth client
│   ├── permission/     # RBAC checker
│   ├── errdefs/        # error taxonomy (maps → HTTP status)
│   ├── logger/         # zap wrapper
│   └── pb/             # generated Protobuf (widget, page, websocket, …)
├── migrations/         # SQL schema
└── go.mod
```