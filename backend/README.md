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
- `/logger` - Logging utilities
- `/migrations` - Database migrations
- `/pkg` - Public library code usable by external applications
