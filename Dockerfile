# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY frontend/ ./
# Copy .env file for frontend build
COPY .env /app/frontend/.env
RUN yarn build

# Stage 2: Build backend
FROM golang:1.23 AS backend-builder
WORKDIR /go/src/github.com/trysourcetool/sourcetool
COPY proto/ ./proto/
COPY backend/ ./backend/
WORKDIR /go/src/github.com/trysourcetool/sourcetool/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server ./cmd/server

# Stage 3: Final image
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata postgresql-client curl

# Install migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

WORKDIR /app
COPY --from=backend-builder /go/bin/server /app/
# Copy the entire frontend build directory
COPY --from=frontend-builder /app/frontend/build /app/static-full
# Create a symlink from /app/static to /app/static-full/client
RUN ln -s /app/static-full/client /app/static

# Create debug output directory
RUN mkdir -p /debug_output

# Verify index.html exists in the correct location
RUN echo "Checking for index.html in static-full" && \
    find /app/static-full -name "index.html" && \
    echo "Checking symlink" && \
    ls -la /app/static

COPY backend/migrations /app/migrations
COPY entrypoint.sh /app/
COPY .env /app/

# Make entrypoint executable
RUN chmod +x /app/entrypoint.sh

# Expose port
EXPOSE 8080

# Set environment variables
ENV STATIC_FILES_DIR=/app/static

# Use our custom entrypoint script
ENTRYPOINT ["/app/entrypoint.sh"]
