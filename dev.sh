#!/bin/bash

# Development setup script for WebSocket Chat Application

set -e

echo "WebSocket Chat Application - Development Setup"
echo "================================================"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${BLUE}$1${NC}"
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check if .env exists
if [ ! -f ".env" ]; then
    print_info "Creating .env file from .env.example"
    cp .env.example .env
    echo "Please edit .env with your MongoDB URI"
fi

# Function to start backend
start_backend() {
    print_step "Starting Go backend server..."

    # Install Go dependencies
    go mod download

    # Create default users
    print_step "Creating default users..."
    go run create_users.go

    # Start backend server
    print_step "Starting backend on port 8080..."
    go run . &
    BACKEND_PID=$!

    # Wait for backend to start
    sleep 3
    print_success "Backend server started (PID: $BACKEND_PID)"
}

# Function to start frontend
start_frontend() {
    print_step "Starting React frontend..."

    cd frontend

    # Install npm dependencies
    if [ ! -d "node_modules" ]; then
        print_step "Installing npm dependencies..."
        npm install
    fi

    # Start frontend server
    print_step "Starting frontend on port 3000..."
    npm start &
    FRONTEND_PID=$!

    cd ..

    # Wait for frontend to start
    sleep 5
    print_success "Frontend server started (PID: $FRONTEND_PID)"
}

# Function to start with Docker
start_docker() {
    print_step "Starting with Docker Compose..."

    # Build and start all services
    docker-compose up --build -d

    print_success "All services started with Docker"
    print_info "Frontend: http://localhost:3000"
    print_info "Backend API: http://localhost:8080"
    print_info "Health Check: http://localhost:8080/health"
}

# Function to stop services
stop_services() {
    print_step "Stopping all services..."

    # Stop Docker services
    docker-compose down 2>/dev/null || true

    # Kill background processes
    pkill -f "go run" 2>/dev/null || true
    pkill -f "npm start" 2>/dev/null || true

    print_success "All services stopped"
}

# Function to show logs
show_logs() {
    print_step "Showing application logs..."
    docker-compose logs -f
}

# Function to test the application
test_app() {
    print_step "Testing application endpoints..."

    # Test health endpoint
    echo "Testing health endpoint..."
    curl -s http://localhost:8080/health | python3 -m json.tool || echo "Health check failed"

    # Test login endpoint
    echo -e "\nTesting login endpoint..."
    curl -s -X POST http://localhost:8080/login \
        -H "Content-Type: application/json" \
        -d '{"username":"alice","password":"password123"}' | \
        python3 -m json.tool || echo "Login test failed"

    print_success "Basic tests completed"
}

# Main menu
case "${1:-menu}" in
    "backend")
        start_backend
        print_info "Backend running on http://localhost:8080"
        print_info "Press Ctrl+C to stop"
        wait
        ;;
    "frontend")
        start_frontend
        print_info "Frontend running on http://localhost:3000"
        print_info "Press Ctrl+C to stop"
        wait
        ;;
    "both")
        start_backend
        start_frontend
        print_info "Both services running:"
        print_info "Frontend: http://localhost:3000"
        print_info "Backend: http://localhost:8080"
        print_info "Press Ctrl+C to stop both"
        wait
        ;;
    "docker")
        start_docker
        ;;
    "stop")
        stop_services
        ;;
    "logs")
        show_logs
        ;;
    "test")
        test_app
        ;;
    "menu"|*)
        echo ""
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  backend   - Start only Go backend server"
        echo "  frontend  - Start only React frontend"
        echo "  both      - Start both backend and frontend"
        echo "  docker    - Start with Docker Compose"
        echo "  stop      - Stop all services"
        echo "  logs      - Show Docker logs"
        echo "  test      - Test application endpoints"
        echo ""
        echo "Examples:"
        echo "  $0 both     # Start both services for development"
        echo "  $0 docker   # Start with Docker for production-like testing"
        echo "  $0 test     # Test the running application"
        echo ""
        ;;
esac
