FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

# Set environment variables for the build
ENV VITE_IS_CLOUD_EDITION=true

COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY frontend/ ./
RUN yarn build

FROM golang:1.23 AS backend-builder
WORKDIR /go/src/github.com/trysourcetool/sourcetool
COPY backend/ ./backend/
WORKDIR /go/src/github.com/trysourcetool/sourcetool/backend/ee
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server ./cmd/server
# Build migration tool
WORKDIR /go/src/github.com/trysourcetool/sourcetool/backend/devtools/cmd/db
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/migrate

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /go/bin/server /app/
COPY --from=backend-builder /go/bin/migrate /app/
COPY backend/migrations /app/migrations
COPY --from=frontend-builder /app/frontend/build /app/static-full
RUN ln -s /app/static-full/client /app/static
ENV STATIC_FILES_DIR=/app/static
EXPOSE 8080

# Default command to start the server
CMD ["/app/server"]
