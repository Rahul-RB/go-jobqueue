.PHONY: build up down logs clean help

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build all Docker images"
	@echo "  up       - Start all services"
	@echo "  down     - Stop all services"
	@echo "  logs     - Show logs from all services"
	@echo "  clean    - Remove all containers and images"
	@echo "  scale    - Scale app instances (usage: make scale N=5)"

# Build all images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Clean up everything
clean:
	docker-compose down -v --rmi all --volumes --remove-orphans

# Scale app instances (default to 3)
scale:
	docker-compose up -d --scale app=$${N:-3}

# Show status
status:
	docker-compose ps

# Access Traefik dashboard
dashboard:
	@echo "Traefik dashboard available at: http://localhost:8080"

# Test the application
test:
	@echo "Testing application endpoints..."
	@echo "Creating a job..."
	@curl -X POST http://localhost/v1/job
	@echo ""
	@echo "Application is running at: http://localhost"
