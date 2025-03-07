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

- `/apikey` - API key management
- `/authz` - Authorization logic
- `/cmd` - Application entry points
- `/config` - Configuration handling
- `/devtools` - Development tools
- `/docs` - API documentation (Swagger)
- `/dto` - Data transfer objects
- `/environment` - Environment management
- `/group` - Group management
- `/hostinstance` - Host instance management
- `/infra` - Infrastructure components
- `/migrations` - Database migrations
- `/organization` - Organization management
- `/page` - Page management
- `/server` - HTTP server implementation
- `/session` - Session management
- `/user` - User management
- `/ws` - WebSocket implementation
