---
sidebar_position: 3
---

# Deployment Guide

This guide will help you deploy Sourcetool in your environment using Docker. The deployment process is straightforward and requires minimal setup.

## Prerequisites

- Docker installed on your system
- Access to a server or cloud environment where you want to deploy
- Basic understanding of environment variables and Docker
- Environment that supports WebSocket connections

## Infrastructure Requirements

Sourcetool requires the following infrastructure:

- PostgreSQL database (version 15)
- Redis (version 7)
- WebSocket-capable environment (for real-time features)

You can use your preferred method to host these services:
- Managed database services (e.g., Amazon RDS, Google Cloud SQL)
- Managed Redis services (e.g., Amazon ElastiCache, Google Cloud Memorystore)
- Self-hosted instances

Note: If you're using a reverse proxy or load balancer, make sure it's configured to support WebSocket connections.

## Docker Image

The application is available as a Docker image that contains both the frontend and backend components:

```
ghcr.io/trysourcetool/sourcetool:latest
```

The image exposes the following:
- Port: 8080
- Health check endpoint: `/api/health`
- WebSocket endpoint: `/ws`

You can also use specific version tags instead of `latest` for better stability:
```
ghcr.io/trysourcetool/sourcetool:v1.0.0
```

## Deployment Options

### 1. Container Services

You can deploy the application on various container platforms:

- **Google Cloud Run**
  - Supports WebSocket connections
  - Set minimum instances to 1 for better performance
  - Configure memory and CPU based on your needs (recommended: 1CPU, 2GB memory)
  - Use Cloud SQL and Memorystore for managed database services
  - Enable session affinity for WebSocket connections
  - Configure custom domain and SSL/TLS certificates

- **AWS ECS/Fargate**
  - Configure Application Load Balancer with WebSocket support
  - Use RDS and ElastiCache for managed services
  - Set up Auto Scaling for the ECS service
  - Configure target groups with appropriate health checks
  - Use AWS Certificate Manager for SSL/TLS

- **Azure Container Apps**
  - Enable WebSocket support in configuration
  - Use Azure Database for PostgreSQL and Azure Cache for Redis
  - Configure scaling rules and minimum replica count
  - Set up custom domains and managed certificates

For all container services, ensure:
- Memory and CPU settings are appropriate (recommended: 1CPU, 2GB memory)
- Health check endpoints are properly configured
- Environment variables are securely stored (using secret management services)
- Network policies allow communication between services

### 2. Self-hosted

For self-hosted environments, you can run the container directly on your server. Follow the deployment steps below for detailed instructions.

## Environment Variables

The application requires several environment variables to be set. Below are the essential variables you need to configure for production:

```sh
# Environment
ENV=prod
BASE_URL=https://your-domain.com

# Security (make sure to use strong, unique values)
ENCRYPTION_KEY=<your-secure-encryption-key> # you can generate this using `make gen-encryption-key`
JWT_KEY=<your-secure-jwt-key> # you can generate this using `make gen-jwt-key`

# Database configuration
# Use your production database connection details
POSTGRES_USER=<your-db-user>
POSTGRES_PASSWORD=<your-secure-password>
POSTGRES_DB=sourcetool
POSTGRES_HOST=<your-db-host>
POSTGRES_PORT=5432

# Redis configuration
# Use your production Redis connection details
REDIS_HOST=<your-redis-host>
REDIS_PASSWORD=<your-secure-redis-password> # if not using a password, leave this empty
REDIS_PORT=6379

# Google OAuth configuration
# For Google OAuth, make sure to configure {BASE_URL}/auth/google/callback as the callback URL in your Google OAuth settings screen
GOOGLE_OAUTH_CLIENT_ID=<your-google-oauth-client-id>
GOOGLE_OAUTH_CLIENT_SECRET=<your-google-oauth-client-secret>


# SMTP configuration
SMTP_HOST=<your-smtp-host>
SMTP_PORT=<your-smtp-port>
SMTP_USERNAME=<your-smtp-username>
SMTP_PASSWORD=<your-smtp-password>
SMTP_FROM_EMAIL=<your-smtp-from-email>
```

Note: The example `.env` file in the repository is configured for local development. Make sure to adjust these values for your production environment.

## Deployment Steps

### 1. Running the Container

For self-hosted environments, you can run the container directly:

```bash
docker run -d \
  --name sourcetool \
  --env-file /path/to/your/production.env \
  -p 8080:8080 \
  ghcr.io/trysourcetool/sourcetool:latest
```

### 2. Database and Redis Connection

Ensure that your application can connect to your production PostgreSQL and Redis instances:

1. Configure the correct connection details in your environment variables
2. Make sure the network allows connections from your application to the database and Redis
3. Use appropriate security groups and firewall rules
4. Consider using SSL/TLS for database connections in production

### 3. Network Configuration

1. **WebSocket Support**:
   - Configure your load balancer or reverse proxy to support WebSocket connections
   - Ensure timeout settings are appropriate for long-lived WebSocket connections
   - Configure proper headers and protocol upgrades for WebSocket support

### 4. Running Migrations

The application will automatically run necessary database migrations on startup. However, if you need to run migrations manually, the Docker image includes a dedicated migration tool:

```bash
# Run migrations using the migration tool
docker run --rm \
  --env-file /path/to/your/production.env \
  ghcr.io/trysourcetool/sourcetool:latest \
  /app/migrate
```

This can be useful in scenarios such as:
- Running migrations before deploying a new version
- Verifying database schema changes
- Troubleshooting database issues

### 5. Verifying the Deployment

Once deployed, you can verify the application is running by accessing:

- Frontend: `https://your-domain.com`
- API Health Check: `https://your-domain.com/api/health`

The health check endpoint will return:
- 200 OK: Application is running correctly
- 503 Service Unavailable: Application is not ready or has issues

## Production Considerations

1. **SSL/TLS**: 
   - Ensure you have SSL/TLS configured for production deployments
   - WebSocket connections should use WSS (WebSocket Secure) in production

2. **Backups**: 
   - Set up regular backups of your PostgreSQL database
   - Consider point-in-time recovery options
   - Test your backup restoration process

3. **Monitoring and Logging**:
   - Set up application monitoring
   - Configure centralized logging
   - Set up alerts for critical errors
   - Monitor system resources (CPU, memory, disk usage)

4. **Scaling and High Availability**:
   - The application can be scaled horizontally by running multiple instances
   - Use a load balancer to distribute traffic
   - Consider using container orchestration platforms (e.g., Kubernetes)
   - Implement health checks and automatic instance recovery

5. **Security**:
   - Use strong, unique passwords for all services
   - Regularly rotate credentials
   - Follow the principle of least privilege for service accounts
   - Keep the Docker image up to date with security patches

## Troubleshooting

If you encounter any issues during deployment:

1. Check the container logs:
   ```bash
   docker logs sourcetool
   ```

2. Common issues to check:
   - Environment variables are correctly set
   - Database and Redis connection details are correct
   - Network connectivity between services
   - WebSocket connections are properly configured and working
   - Sufficient system resources (CPU, memory)
   - Correct permissions for service accounts

3. WebSocket-specific issues:
   - Check if your proxy/load balancer configuration supports WebSocket
   - Verify timeout settings are appropriate
   - Monitor for connection drops or failures
   - Check client-side console for WebSocket connection errors

For additional support or questions, please refer to our [GitHub repository](https://github.com/trysourcetool/sourcetool) or contact our support team. 