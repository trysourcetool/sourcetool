# Sourcetool Frontend

This repository contains the frontend application for Sourcetool, a web-based platform for managing organizations, users, groups, pages, environments, and API keys. The application is built with React and uses modern web technologies to provide a responsive and interactive user experience.

## Prerequisites

- Node.js >= 20.0.0
- pnpm package manager

## Package Installation

To install the project dependencies, run:

```bash
pnpm install
```

## Development

To start the development server, run:

```bash
pnpm dev
```

This will start the development server at `auth.local.trysourcetool.com:5173` or `acme.local.trysourcetool.com:5173`.

Before running the development server, make sure to set up the required environment variables by copying the `.env.example` file to `.env` and updating the values as needed:

```bash
cp .env.example .env
```

## Production Build

To build the application for production, run:

```bash
pnpm build
```

To preview the production build locally, run:

```bash
pnpm preview
```

## Directory Structure

The project follows a modular directory structure:

- `/app`: Main application code
  - `/api`: API client and modules for different resources
  - `/components`: Reusable UI components
    - `/common`: Common components used throughout the application
    - `/icon`: Icon components
    - `/layout`: Layout components
    - `/ui`: UI components (buttons, forms, etc.)
  - `/constants`: Application constants
  - `/environments`: Environment configuration
  - `/hooks`: Custom React hooks
  - `/lib`: Utility functions and libraries
  - `/routes`: Application routes and pages
  - `/store`: Redux store configuration and modules
- `/public`: Static assets and localization files

## Technologies Used

- **React**: UI library
- **React Router v7**: Routing
- **Redux Toolkit**: State management
- **TailwindCSS**: Styling
- **Shadcn UI**: UI component primitives
- **React Hook Form**: Form handling
- **zod**: Schema validation
- **i18next**: Internationalization
- **TypeScript**: Type safety
- **Vite**: Build tool
- **WebSockets**: Real-time communication

## Internationalization

The application supports multiple languages using i18next. Currently, English is the primary language, with Japanese support in development. Localization files are stored in the `/public/locales` directory.

## Type Checking

To run type checking, use:

```bash
pnpm typecheck
```

## Linting

To run the linter, use:

```bash
pnpm lint
```
