##
# Go commands
##
.PHONY: run build test clean proto

dev:
	make env-local
	gow -c run cmd/main.go

build:
	go build -o battledak-server cmd/main.go



##
# Docker commands
##

# Docker compose command - use only the modern version (with space)
DOCKER_COMPOSE := docker compose --env-file ./.env.local -p battledak-server

# Build and start the containers
.PHONY: docker-up docker-down docker-build docker-rebuild docker-clean docker-status docker-logs 
docker-up:
	$(DOCKER_COMPOSE) up -d

# Stop the containers
docker-down:
	$(DOCKER_COMPOSE) down

# Build the containers
docker-build:
	$(DOCKER_COMPOSE) build

# Rebuild the containers without cache
docker-rebuild:
	$(DOCKER_COMPOSE) build --no-cache

# Stop and remove all containers, networks, images, and volumes
docker-clean:
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

# Show the status of the containers
docker-status:
	$(DOCKER_COMPOSE) ps

# Follow logs of all containers
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Build dev
gcloud-build-dev: env-dev
	cp -f .env.dev .env
	gcloud builds submit --config=cloudbuild.yaml --substitutions=_ENV_FILE=.env,_ENV=dev,_PORT=8080

# Build prod
gcloud-build-prod: env-prod
	cp -f .env.prod .env
	gcloud builds submit --config=cloudbuild.yaml --substitutions=_ENV_FILE=.env,_ENV=prod,_PORT=8080

# Deploy dev
deploy-dev: gcloud-build-dev
	gcloud config set run/region us-east1
	gcloud run deploy battledak-server-dev --image gcr.io/gateway-dashboard-front/battledak-server-dev --platform managed --allow-unauthenticated

# Deploy prod
deploy-prod: gcloud-build-prod
	gcloud config set run/region us-east1
	gcloud run deploy battledak-server-prod --image gcr.io/gateway-dashboard-front/battledak-server-prod --platform managed --allow-unauthenticated

## ENV MANAGE
env-local:
	./scripts/env.sh local
env-dev:
	./scripts/env.sh dev
env-staging:
	./scripts/env.sh staging
env-prod:
	./scripts/env.sh prod

env-diff-local:
	./scripts/env.sh -d local
env-diff-dev:
	./scripts/env.sh -d dev
env-diff-staging:
	./scripts/env.sh -d staging
env-diff-prod:
	./scripts/env.sh -d prod

env-save-local:
	./scripts/env.sh -s local
env-save-dev:
	./scripts/env.sh -s dev
env-save-staging:
	./scripts/env.sh -s staging
env-save-prod:
	./scripts/env.sh -s prod

.PHONY: test coverage
.PHONY: testsum

# Command to detect operating system
OS := $(shell uname -s)
SKIP_GO := $(filter-out 0,$(SKIP_GO))

prepare-env:
ifeq ($(SKIP_GO),)
ifeq ($(OS),Linux)
	@echo "Detected OS: Linux"
	# Install Go (assuming Go is already available in the package repository)
	@sudo apt-get update
	@sudo apt-get install -y golang-go
endif

ifeq ($(OS),Darwin)
	@echo "Detected OS: macOS"
	# Install Go using Homebrew
	@brew update
	@brew install go
endif

ifeq ($(OS),Windows_NT)
	@echo "Detected OS: Windows"
	# Check if Chocolatey is installed
	@if exist "%ProgramData%\chocolatey\bin\choco.exe" (\
		echo "Chocolatey is already installed";\
	) else (\
		@echo "Installing Chocolatey...";\
		@powershell -NoProfile -ExecutionPolicy Bypass -Command "Set-ExecutionPolicy Bypass -Scope Process; [System.Net.ServicePointManager]::SecurityProtocol = 'Tls12'; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))";\
	)
	# Install Go using Chocolatey
	@choco install golang -y
endif
endif

	# Ensure Go dependencies are installed
	@go mod tidy

	# Install gow
	@go install github.com/mitranim/gow@latest

	@echo "Environment setup complete"