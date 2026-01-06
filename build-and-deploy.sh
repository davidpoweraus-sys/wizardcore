#!/bin/bash
set -e

echo "========================================="
echo "WizardCore Build & Deploy Script"
echo "========================================="
echo "Builds images locally on server"
echo "Pulls external dependencies from Docker Hub"
echo "========================================="

# Configuration
BACKEND_IMAGE="wizardcore-backend:local"
FRONTEND_IMAGE="wizardcore-frontend:local"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${GREEN}[STEP]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker daemon."
    exit 1
fi

# Function to pull external images
pull_external_images() {
    print_step "Pulling external dependency images..."
    
    # List of external images we need
    EXTERNAL_IMAGES=(
        "postgres:15-alpine"
        "postgres:16-alpine"
        "redis:7-alpine"
        "supabase/gotrue:v2.184.0"
        "judge0/judge0:latest"
    )
    
    for image in "${EXTERNAL_IMAGES[@]}"; do
        print_info "Pulling $image..."
        if docker pull "$image" > /dev/null 2>&1; then
            print_info "✓ $image pulled successfully"
        else
            print_error "Failed to pull $image"
            # Continue anyway - Docker will pull during compose up if needed
        fi
    done
}

# Function to build backend
build_backend() {
    print_step "Building wizardcore-backend..."
    
    if [ -d "wizardcore-backend" ]; then
        cd wizardcore-backend
        
        # Clean up old images
        docker rmi -f "$BACKEND_IMAGE" 2>/dev/null || true
        
        # Build with timestamp for cache busting
        print_info "Building backend image with timestamp: $TIMESTAMP"
        docker build \
            --build-arg BUILD_TIMESTAMP="$TIMESTAMP" \
            -t "$BACKEND_IMAGE" \
            -t "wizardcore-backend:$TIMESTAMP" \
            .
        
        if [ $? -eq 0 ]; then
            print_info "✓ Backend built successfully: $BACKEND_IMAGE"
        else
            print_error "Backend build failed"
            exit 1
        fi
        
        cd ..
    else
        print_error "wizardcore-backend directory not found"
        exit 1
    fi
}

# Function to build frontend
build_frontend() {
    print_step "Building wizardcore-frontend..."
    
    if [ -f "Dockerfile.nextjs" ]; then
        # Clean up old images
        docker rmi -f "$FRONTEND_IMAGE" 2>/dev/null || true
        
        # Build with timestamp for cache busting
        print_info "Building frontend image with timestamp: $TIMESTAMP"
        docker build \
            --build-arg BUILD_TIMESTAMP="$TIMESTAMP" \
            -f Dockerfile.nextjs \
            -t "$FRONTEND_IMAGE" \
            -t "wizardcore-frontend:$TIMESTAMP" \
            .
        
        if [ $? -eq 0 ]; then
            print_info "✓ Frontend built successfully: $FRONTEND_IMAGE"
        else
            print_error "Frontend build failed"
            exit 1
        fi
    else
        print_error "Dockerfile.nextjs not found"
        exit 1
    fi
}

# Function to update docker-compose.yml
update_docker_compose() {
    print_step "Updating docker-compose.yml to use local images..."
    
    # Create backup
    cp docker-compose.yml docker-compose.yml.backup.$TIMESTAMP
    
    # Update backend image
    sed -i 's|image: limpet/wizardcore-backend:.*|image: wizardcore-backend:local|g' docker-compose.yml
    
    # Update frontend image
    sed -i 's|image: limpet/wizardcore-frontend:.*|image: wizardcore-frontend:local|g' docker-compose.yml
    
    # Remove pull_policy for local images (not needed)
    sed -i '/pull_policy: always/d' docker-compose.yml
    
    print_info "✓ docker-compose.yml updated to use local images"
    print_info "  Backend: wizardcore-backend:local"
    print_info "  Frontend: wizardcore-frontend:local"
}

# Function to deploy
deploy() {
    print_step "Deploying with docker-compose..."
    
    # Stop existing services
    print_info "Stopping existing services..."
    docker-compose down --remove-orphans
    
    # Start services
    print_info "Starting services..."
    docker-compose up -d
    
    # Wait for services to be healthy
    print_info "Waiting for services to be healthy..."
    sleep 10
    
    # Check service status
    print_info "Checking service status..."
    docker-compose ps
    
    print_info "✓ Deployment complete"
}

# Function to show logs
show_logs() {
    print_step "Showing recent logs..."
    
    echo "========================================="
    echo "Backend logs (last 20 lines):"
    echo "========================================="
    docker-compose logs backend --tail=20
    
    echo ""
    echo "========================================="
    echo "Frontend logs (last 20 lines):"
    echo "========================================="
    docker-compose logs frontend --tail=20
}

# Main execution
main() {
    case "${1:-}" in
        "pull-only")
            pull_external_images
            ;;
        "build-only")
            build_backend
            build_frontend
            ;;
        "deploy-only")
            deploy
            ;;
        "logs")
            show_logs
            ;;
        "full"|"")
            print_step "Starting full build and deploy process..."
            pull_external_images
            build_backend
            build_frontend
            update_docker_compose
            deploy
            show_logs
            ;;
        *)
            echo "Usage: $0 [pull-only|build-only|deploy-only|logs|full]"
            echo ""
            echo "Options:"
            echo "  pull-only    - Pull external images only"
            echo "  build-only   - Build backend/frontend only"
            echo "  deploy-only  - Deploy only (assumes images built)"
            echo "  logs         - Show recent logs"
            echo "  full         - Full build and deploy (default)"
            exit 1
            ;;
    esac
    
    print_step "Script completed successfully!"
}

# Run main function
main "$@"