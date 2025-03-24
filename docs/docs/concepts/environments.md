---
sidebar_position: 2
---

# Environments

Environments in Sourcetool allow you to manage different deployment contexts for your applications. They help you separate development, testing, and production workloads while maintaining consistent code across all stages.

## What are Environments?

An environment in Sourcetool represents a specific deployment context with its own:

- API keys
- Configuration settings
- Access controls

Environments enable you to:

1. Develop and test features without affecting production users
2. Maintain separate configurations for different stages of your application lifecycle
3. Control access to sensitive production environments
4. Deploy changes progressively through your development pipeline

## Environment Types

Sourcetool supports three common environment types:

### Development Environment

- Used for local development and testing
- Typically accessed by developers only
- May contain experimental features and work-in-progress code

### Staging Environment

- Used for pre-production testing
- Mirrors the production environment as closely as possible
- Used for final testing before deploying to production
- May be accessed by QA teams, product managers, and stakeholders

### Production Environment

- Serves your end users
- Requires the highest level of stability and security
- Changes are carefully reviewed and tested before deployment
- Often has stricter access controls

## Next Steps

Now that you understand environments, learn about:

- [Organizations](./organizations) for managing users and access control
- [Pages](./pages) for building your application's UI
<!-- TODO: Add API Keys documentation -->
<!-- - [API Keys](../reference/api-keys) for authenticating your application -->
