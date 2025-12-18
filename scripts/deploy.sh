#!/bin/bash
set -e

# Deployment script for wizardcore on Coolify
# Usage: ./scripts/deploy.sh [ENV_FILE]
# If ENV_FILE is provided, it will be passed to docker compose via --env-file

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

ENV_FILE="${1:-.env}"
COMPOSE_FILE="docker-compose.prod.yml"
PROJECT_NAME="${COOLIFY_RESOURCE_UUID:-wizardcore}"

echo "=== WizardCore Deployment ==="
echo "Project root: $PROJECT_ROOT"
echo "Compose file: $COMPOSE_FILE"
echo "Project name: $PROJECT_NAME"
echo "Environment file: $ENV_FILE"

if [[ ! -f "$ENV_FILE" ]]; then
    echo "Warning: Environment file $ENV_FILE not found. Proceeding without it."
    ENV_FILE=""
fi

# Stop and remove any existing containers from the same project
echo "Stopping existing containers..."
docker compose --project-name "$PROJECT_NAME" -f "$COMPOSE_FILE" down --remove-orphans

# Build and start services
echo "Building and starting services..."
if [[ -n "$ENV_FILE" ]]; then
    docker compose --project-name "$PROJECT_NAME" -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --build
else
    docker compose --project-name "$PROJECT_NAME" -f "$COMPOSE_FILE" up -d --build
fi

echo "Deployment completed successfully."
echo "Running containers:"
docker compose --project-name "$PROJECT_NAME" -f "$COMPOSE_FILE" ps