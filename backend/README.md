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

- Go 1.22
- PostgreSQL 15
- Redis 7
- Docker

## Directory Structure

- `/cmd` - Application entry points
- `/config` - Configuration handling
- `/devtools` - Development tools
- `/docs` - API documentation (Swagger)
- `/ee` - Enterprise features (if applicable)
- `/fixtures` - Test fixtures and data
- `/internal` - Internal application logic and packages
  - `/app` - Application layer services and use cases
  - `/domain` - Core domain models and business logic
  - `/infra` - Infrastructure components (database, external services)
    - `/db` - Database related code and repositories
      - `/postgres` - PostgreSQL specific implementations
    - `/email` - Email service implementations
      - `/smtp` - SMTP based email service
    - `/pubsub` - Publish/Subscribe implementations
      - `/redis` - Redis based pub/sub service
    - `/ws` - WebSocket related implementations
      - `/manager` - WebSocket connection management
  - `/pb` - Protocol Buffer definitions and generated code
  - `/permission` - Permission checking logic
  - `/transport` - API transport layer (HTTP and WebSocket handlers)
    - `/http` - HTTP API handlers and routing
      - `/v1` - Version 1 of the HTTP API
        - `/handlers` - HTTP request handlers
        - `/requests` - Request models and validation
        - `/responses` - Response models
        - `/mapper` - Data mapping utilities
        - `/middleware` - HTTP middleware components
    - `/ws` - WebSocket handlers and routing
      - `/handlers` - WebSocket message handlers
      - `/message` - WebSocket message utilities
      - `/middleware` - WebSocket middleware components
    - `router.go` - Main router configuration
    - `static.go` - Static file serving configuration
  - `/jwt` - JWT handling utilities
  - `/ctxutil` - Context utility functions
- `/logger` - Logging utilities
- `/migrations` - Database migrations
- `/pkg` - Public library code usable by external applications
