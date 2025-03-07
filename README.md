# Sourcetool

This repository contains the Sourcetool application with both backend and frontend components.

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Node.js 20+ (for local development)
- Make

### Setup

1. Copy the sample environment file:
   ```bash
   cp .env.example .env
   ```

2. Generate required keys and add them to your `.env` file:
   ```bash
   make gen-encryption-key
   make gen-jwt-key
   ```

3. Start the application:
   ```bash
   # For Community Edition (CE)
   make up
   
   # For Enterprise Edition (EE)
   make up-ee
   ```

4. Access the application:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - API Documentation: http://localhost:8080/swagger/index.html

## Project Structure

- `/backend` - Go backend service
- `/frontend` - React frontend application
- `/proto` - Protocol Buffers definitions
- `/compose.yaml` - Docker Compose configuration for CE
- `/compose.ee.yaml` - Docker Compose configuration for EE
- `/Makefile` - Consolidated commands for development

## Available Commands

Run `make help` to see all available commands. Here are the most commonly used ones:

### Docker Compose Commands

- `make up` - Start the Community Edition (CE) services
- `make up-ee` - Start the Enterprise Edition (EE) services
- `make down` - Stop the CE services
- `make down-ee` - Stop the EE services
- `make logs` - View logs for CE services
- `make logs-ee` - View logs for EE services

### Development Commands

- `make gen-keys` - Generate both encryption and JWT keys
- `make swagger` - Generate Swagger documentation
- `make swagger-open` - Open Swagger UI in browser
- `make lint` - Run linters on the CE codebase
- `make lint-ee` - Run linters on the EE codebase

### Database Commands

- `make db-migrate` - Run database migrations

### Protocol Buffer Commands

- `make proto-generate` - Generate Go code from proto files
- `make proto-lint` - Lint proto files
- `make proto-format` - Format proto files

### Maintenance Commands

- `make remove-docker-images` - Remove untagged Docker images
- `make remove-docker-builder` - Prune Docker builder cache

## Environment Variables

The application uses a single `.env` file at the root level for all services. See `.env.example` for the required variables.

## Development

### Backend

The backend is a Go application that provides the API for the frontend. See the [backend README](backend/README.md) for more details.

### Frontend

The frontend is a React application built with Vite. See the [frontend README](frontend/README.md) for more details.

### Protocol Buffers

The application uses Protocol Buffers for API definitions. See the [proto README](proto/README.md) for more details.
