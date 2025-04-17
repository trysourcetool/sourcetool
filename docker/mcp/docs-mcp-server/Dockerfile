# Use Node.js 20 Alpine as base image for smaller size
FROM node:20-alpine

# Set up pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

# Create app directory and set it as the working directory
WORKDIR /app

# Copy workspace configuration
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY mcp/docs-mcp-server/package.json ./mcp/docs-mcp-server/

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy application files
COPY mcp/docs-mcp-server/ ./mcp/docs-mcp-server/

# Build the application
WORKDIR /app/mcp/docs-mcp-server
RUN pnpm build && ls -la

# Remove development dependencies to reduce image size
RUN pnpm prune --prod

# Use non-root user for security
USER node

# Command to run the application
CMD ["node", "dist/index.js"] 