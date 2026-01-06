#!/bin/bash
set -e

echo "========================================="
echo "WizardCore Git Source Deployment Script"
echo "========================================="
echo "Builds images from git source on server"
echo "Deploys using docker-compose with build context"
echo "========================================="

# Configuration
BACKEND_IMAGE="wizardcore-backend"
FRONTEND_IMAGE="wizardcore-frontend"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BUILD_TAG="${BUILD_TAG:-$TIMESTAMP}"

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

# Function to verify git source
verify_git_source() {
    print_step "Verifying git source..."
    
    if [ ! -d ".git" ]; then
        print_error "Not a git repository. This script should run in a git clone."
        print_info "If deploying from Dokploy, ensure repository is cloned properly."
        exit 1
    fi
    
    print_info "Git repository verified"
    print_info "Current commit: $(git rev-parse --short HEAD)"
    print_info "Branch: $(git branch --show-current)"
}

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
            print_info "Continuing - Docker will pull during compose up if needed"
        fi
    done
}

# Function to build backend with timestamp
build_backend() {
    print_step "Building wizardcore-backend..."
    
    if [ -d "wizardcore-backend" ]; then
        cd wizardcore-backend
        
        # Clean up old images with same tag
        docker rmi -f "$BACKEND_IMAGE:$BUILD_TAG" 2>/dev/null || true
        
        # Build with timestamp for cache busting
        print_info "Building backend image with tag: $BUILD_TAG"
        docker build \
            --build-arg BUILD_TIMESTAMP="$TIMESTAMP" \
            -t "$BACKEND_IMAGE:$BUILD_TAG" \
            -t "$BACKEND_IMAGE:latest" \
            .
        
        if [ $? -eq 0 ]; then
            print_info "✓ Backend built successfully: $BACKEND_IMAGE:$BUILD_TAG"
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

# Function to build frontend with timestamp
build_frontend() {
    print_step "Building wizardcore-frontend..."
    
    if [ -f "Dockerfile.nextjs" ]; then
        # Clean up old images with same tag
        docker rmi -f "$FRONTEND_IMAGE:$BUILD_TAG" 2>/dev/null || true
        
        # Build with timestamp for cache busting
        print_info "Building frontend image with tag: $BUILD_TAG"
        docker build \
            --build-arg BUILD_TIMESTAMP="$TIMESTAMP" \
            --build-arg CACHE_BUST="$TIMESTAMP" \
            -f Dockerfile.nextjs \
            -t "$FRONTEND_IMAGE:$BUILD_TAG" \
            -t "$FRONTEND_IMAGE:latest" \
            .
        
        if [ $? -eq 0 ]; then
            print_info "✓ Frontend built successfully: $FRONTEND_IMAGE:$BUILD_TAG"
        else
            print_error "Frontend build failed"
            exit 1
        fi
    else
        print_error "Dockerfile.nextjs not found"
        exit 1
    fi
}

# Function to set BUILD_TAG in environment
set_build_tag() {
    print_step "Setting BUILD_TAG environment variable..."
    
    # Export for current shell
    export BUILD_TAG="$BUILD_TAG"
    
    # Create/update .env file with BUILD_TAG
    if [ -f ".env" ]; then
        if grep -q "BUILD_TAG=" .env; then
            sed -i "s/BUILD_TAG=.*/BUILD_TAG=$BUILD_TAG/" .env
        else
            echo "BUILD_TAG=$BUILD_TAG" >> .env
        fi
    else
        echo "BUILD_TAG=$BUILD_TAG" > .env
    fi
    
    print_info "BUILD_TAG set to: $BUILD_TAG"
    print_info "docker-compose.yml will use: wizardcore-*:\${BUILD_TAG:-latest}"
}

# Function to deploy with docker-compose build
deploy_with_build() {
    print_step "Deploying with docker-compose (building from source)..."
    
    # Stop existing services
    print_info "Stopping existing services..."
    docker-compose down --remove-orphans
    
    # Build and deploy
    print_info "Building and starting services..."
    BUILD_TAG="$BUILD_TAG" docker-compose up -d --build
    
    # Wait for services to be healthy
    print_info "Waiting for services to be healthy..."
    sleep 15
    
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

# Function to verify deployment
verify_deployment() {
    print_step "Verifying deployment..."
    
    # Check if containers are running
    BACKEND_STATUS=$(docker-compose ps backend | grep -c "Up")
    FRONTEND_STATUS=$(docker-compose ps frontend | grep -c "Up")
    
    if [ "$BACKEND_STATUS" -eq 1 ] && [ "$FRONTEND_STATUS" -eq 1 ]; then
        print_info "✓ Both backend and frontend are running"
        
        # Test backend health endpoint
        print_info "Testing backend health endpoint..."
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            print_info "✓ Backend health check passed"
        else
            print_error "Backend health check failed"
        fi
        
        # Test frontend
        print_info "Testing frontend..."
        if curl -s http://localhost:3001 > /dev/null 2>&1; then
            print_info "✓ Frontend is responding"
        else
            print_error "Frontend is not responding"
        fi
    else
        print_error "Some services are not running"
        docker-compose ps
    fi
}

# Function for Dokploy-style deployment
dokploy_deploy() {
    print_step "Dokploy-style deployment (build from source)..."
    
    # This is what Dokploy would run
    print_info "Running: docker-compose -f ./docker-compose.yml up -d --build --remove-orphans"
    
    # Set BUILD_TAG for the compose command
    BUILD_TAG="$BUILD_TAG" docker-compose -f ./docker-compose.yml up -d --build --remove-orphans
    
    if [ $? -eq 0 ]; then
        print_info "✓ Dokploy-style deployment completed"
        docker-compose ps
    else
        print_error "Dokploy-style deployment failed"
        exit 1
    fi
}

# Main execution
main() {
    case "${1:-}" in
        "verify-source")
            verify_git_source
            ;;
        "pull-only")
            pull_external_images
            ;;
        "build-only")
            verify_git_source
            build_backend
            build_frontend
            set_build_tag
            ;;
        "deploy-only")
            deploy_with_build
            ;;
        "dokploy")
            verify_git_source
            pull_external_images
            build_backend
            build_frontend
            set_build_tag
            dokploy_deploy
            ;;
        "logs")
            show_logs
            ;;
        "verify")
            verify_deployment
            ;;
        "full"|"")
            print_step "Starting full deployment from git source..."
            verify_git_source
            pull_external_images
            build_backend
            build_frontend
            set_build_tag
            deploy_with_build
            show_logs
            verify_deployment
            ;;
        *)
            echo "Usage: $0 [verify-source|pull-only|build-only|deploy-only|dokploy|logs|verify|full]"
            echo ""
            echo "Options:"
            echo "  verify-source - Verify git repository"
            echo "  pull-only     - Pull external images only"
            echo "  build-only    - Build backend/frontend only"
            echo "  deploy-only   - Deploy only (assumes images built)"
            echo "  dokploy       - Dokploy-style deployment (build from source)"
            echo "  logs          - Show recent logs"
            echo "  verify        - Verify deployment status"
            echo "  full          - Full build and deploy from git source (default)"
            echo ""
            echo "Environment variables:"
            echo "  BUILD_TAG     - Tag for built images (default: timestamp)"
            exit 1
            ;;
    esac
    
    print_step "Script completed successfully!"
    echo ""
    echo "========================================="
    echo "Deployment Summary"
    echo "========================================="
    echo "Build Tag: $BUILD_TAG"
    echo "Backend Image: $BACKEND_IMAGE:$BUILD_TAG"
    echo "Frontend Image: $FRONTEND_IMAGE:$BUILD_TAG"
    echo "Git Commit: $(git rev-parse --short HEAD)"
    echo "========================================="
}

# Run main function
main "$@"