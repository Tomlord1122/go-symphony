# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Database commands
db-start:
	@echo "Starting Supabase local development..."
	supabase start

db-stop:
	@echo "Stopping Supabase local development..."
	supabase stop

db-up:
	@echo "Running database migrations..."
	supabase migration up

# Database migration commands
db-schema:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make db-schema name=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(name)"
	@migrate create -ext sql -dir ./supabase/migrations $(name)


.PHONY: all build run test clean watch db-start db-stop db-up db-schema
