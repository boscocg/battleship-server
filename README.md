# Battleship Server

A Go-based backend server application with Redis integration, containerized with Docker and deployed to Google Cloud Run.

## 🚀 Features

- Go API server with robust configuration management
- Redis integration for caching and data storage
- Docker containerization for consistent environments
- Automatic deployment to Google Cloud Run
- Environment-based configuration system

## 📋 Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Google Cloud SDK (gcloud)
- Redis (for local development)
- Make (for automation)

## 🏗️ Project Structure

```
battledak-server/
├── cmd/
│   └── main.go         # Application entry point
├── configs/            # Configuration management
│   └── env.go          # Environment loading utilities
├── internal/           # Internal application code
│   └── routes/         # API routes
├── scripts/            # Utility scripts
│   ├── env.sh          # Environment file management
│   ├── deploy.sh       # Main deployment script
│   ├── deploy-base.sh  # Base deployment functionality
│   ├── deploy-dev.sh   # Development deployment
│   ├── deploy-prod.sh  # Production deployment
├── .env                # Environment variables (not in git)
├── Dockerfile          # Container definition
├── docker-compose.yml  # Local development services
├── cloudbuild.yaml     # Google Cloud Build config
├── makefile            # Build and deployment tasks
└── README.md           # Project documentation
```

## 🛠️ Setup and Installation

### Local Development

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd battledak-server
   ```

2. Set up your local environment:
   ```bash
   # Create local environment file
   make env-local
   
   # Start Redis and other required services
   make docker-up
   ```

3. Run the application:
   ```bash
   # Using gow for hot reloading during development
   make dev
   
   # Or build and run manually
   make build
   ./battledak-server
   ```

### Environment Configuration

Environment variables are managed through `.env.<environment>` files:

- `make env-local` - Generate local environment file
- `make env-dev` - Generate development environment file
- `make env-prod` - Generate production environment file
- `make env-staging` - Generate staging environment file

To save changes to environment files in Google Cloud Secret Manager:
```bash
make env-save-dev    # Save development environment
make env-save-prod   # Save production environment
```

## 🐳 Docker

The application is containerized for consistent development and deployment environments.

### Local Docker Development

```bash
# Build Docker containers
make docker-build

# Start all services
make docker-up

# View container status
make docker-status

# View logs
make docker-logs

# Stop containers
make docker-down

# Clean up containers, images, and volumes
make docker-clean
```

## 🚢 Deployment

The application uses Google Cloud Run for deployment with automatic CI/CD.

### Manual Deployment

Deploy to different environments:
```bash
make deploy-dev      # Deploy to development
make deploy-staging  # Deploy to staging
make deploy-prod     # Deploy to production
```

Or use the unified script with environment parameter:
```bash
./scripts/deploy.sh dev|staging|prod
```

### CI/CD Integration

The repository includes GitHub Actions workflows for automated deployment to Cloud Run:
- Pushes to `develop` branch are deployed to development environment
- Pushes to `main` branch are deployed to production environment

## 🔧 Configuration Management

Configuration is loaded from environment variables with environment-specific defaults:

- `ENV` - Specifies the current environment (local, dev, staging, prod)
- `PORT` - The port the server listens on (default: 8080)
- `REDIS_ADDR` - Redis connection string

## 🤝 Contributing

1. Create a feature branch from develop
2. Implement your changes
3. Run tests
4. Create a pull request to develop

## 📄 License

[Add your license information here]
