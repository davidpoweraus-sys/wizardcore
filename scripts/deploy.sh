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

# Function to validate environment file
validate_env() {
    local env_file="$1"
    if [[ ! -f "$env_file" ]]; then
        echo "Error: Environment file $env_file not found."
        echo "Please create it from .env.example or run setup script."
        exit 1
    fi

    # List of required variables (non-empty)
    local required_vars=(
        "NEXT_PUBLIC_SUPABASE_URL"
        "NEXT_PUBLIC_SUPABASE_ANON_KEY"
        "SUPABASE_JWT_SECRET"
        "DATABASE_URL"
        "REDIS_URL"
        "JUDGE0_API_URL"
        "JUDGE0_API_KEY"
    )

    local missing=()
    for var in "${required_vars[@]}"; do
        if ! grep -q "^$var=" "$env_file" || [[ -z "$(grep "^$var=" "$env_file" | cut -d= -f2-)" ]]; then
            missing+=("$var")
        fi
    done

    if [[ ${#missing[@]} -gt 0 ]]; then
        echo "Warning: The following required environment variables are missing or empty in $env_file:"
        printf '%s\n' "${missing[@]}"
        echo "Please set them before deployment."
        echo "Continue anyway? (y/N)"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        echo "Environment file validation passed."
    fi
}

# Validate environment file if it exists
if [[ -f "$ENV_FILE" ]]; then
    validate_env "$ENV_FILE"
else
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

echo "Waiting for services to become healthy (max 2 minutes)..."
timeout=120
interval=5
elapsed=0
healthy=false

while [[ $elapsed -lt $timeout ]]; do
    # Check if all services are running (not restarting) and health status if available
    status_output=$(docker compose --project-name "$PROJECT_NAME" -f "$COMPOSE_FILE" ps --format json 2>/dev/null || true)
    if [[ -z "$status_output" ]]; then
        echo "Could not retrieve container status. Retrying..."
        sleep $interval
        ((elapsed+=interval))
        continue
    fi

    # Count total services and healthy ones
    total=$(echo "$status_output" | jq -r 'length' 2>/dev/null || echo "0")
    if [[ $total -eq 0 ]]; then
        sleep $interval
        ((elapsed+=interval))
        continue
    fi

    # Determine if all services are healthy (or at least running)
    all_running=true
    for i in $(seq 0 $((total - 1))); do
        service=$(echo "$status_output" | jq -r ".[$i].Service")
        state=$(echo "$status_output" | jq -r ".[$i].State")
        health=$(echo "$status_output" | jq -r ".[$i].Health // \"\"")
        if [[ "$state" != "running" ]]; then
            all_running=false
            echo "Service $service is not running (state: $state)"
            break
        fi
        if [[ -n "$health" && "$health" != "healthy" ]]; then
            all_running=false
            echo "Service $service is not healthy (health: $health)"
            break
        fi
    done

    if $all_running; then
        healthy=true
        break
    fi

    sleep $interval
    ((elapsed+=interval))
done

if $healthy; then
    echo "All services are healthy."
else
    echo "Warning: Some services may not be fully healthy after waiting $timeout seconds."
fi

echo "Deployment completed successfully."
echo "Running containers:"
docker compose --project-name "$PROJECT_NAME" -f "$COMPOSE_FILE" ps