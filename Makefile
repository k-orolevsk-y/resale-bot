.PHONY: all build up down

all: build up

check-env-file:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env and configure it."; \
		exit 1; \
	fi

build: check-env-file
	docker-compose build

up: check-env-file
	docker-compose up -d

down:
	docker-compose down

set-manager:
	@if [ -z "$(id)" ]; then \
    	echo "Usage: make set-manager id=<telegram_user_id>"; \
    	exit 1; \
    fi
	@docker-compose exec postgres psql -U bot -d bot -c "UPDATE users SET is_manager = true WHERE id = $(id);" | \
    	  grep "UPDATE 1" > /dev/null 2>&1; \
    	  if [ $$? -eq 0 ]; then \
    	  	echo "SUCCESS"; \
    	  else \
    	  	echo "ERROR"; \
    	fi

help:
	@echo "Available targets:"
	@echo "  - all:          Build and start the application."
	@echo "  - build:        Build Docker containers."
	@echo "  - up:           Start Docker containers in the background."
	@echo "  - down:         Stop and remove Docker containers."
	@echo "  - set-manager:  Set user a manager in bot. Usage: make set-manager id=<user_id>"