# Sourcetool Backend API

A Go-based backend service for the Sourcetool project.

## Quick Start

1. Copy the sample environment file:
   ```
   cp .env.example .env
   ```

2. Generate required keys:
   ```
   make gen-encryption-key  # Add to .env
   make gen-jwt-key         # Add to .env
   ```

3. Start the development server:
   ```
   make dc-up
   ```

4. Access:
   - API: http://localhost:8080
   - API Documentation: http://localhost:8080/swagger/index.html

## Technology

- Go 1.22
- PostgreSQL 15
- Redis 7
- Docker

## Useful Commands

- `make gen-swagger` - Update API documentation
- `make run-lint` - Run linters