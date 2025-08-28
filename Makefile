.PHONY: migrate-up migrate-down migrate-create generate-sqlc mocks machine

# Database Variables
MIGRATION_DIR ?= db/sqlite/migrations
DATABASE_TYPE ?= sqlite3
DATABASE_URL ?= db/sqlite/kumo.db
MIGRATION_TOOL = goose

# Generation Variables
SQLC_CMD = sqlc generate
MOCKERY_CMD = mockery
MACHINE_CMD = go run cmd/external-device/main.go

## Database Migration Rules
## -----------------------

migrate-up:
	@echo "Applying database migrations..."
	$(MIGRATION_TOOL) -dir $(MIGRATION_DIR) $(DATABASE_TYPE) $(DATABASE_URL) up

migrate-down:
	@echo "Reverting last migration..."
	$(MIGRATION_TOOL) -dir $(MIGRATION_DIR) $(DATABASE_TYPE) $(DATABASE_URL) down

migrate-create:
	@echo "Creating new migration file: $(NAME)..."
	@if [ -z "$(NAME)" ]; then echo "Error: NAME=<migration_name> is required" && exit 1; fi
	$(MIGRATION_TOOL) -dir $(MIGRATION_DIR) $(DATABASE_TYPE) $(DATABASE_URL) create $(NAME) sql

## Code Generation Rules
## ---------------------

generate-sqlc:
	@echo "Generating SQL code with sqlc..."
	$(SQLC_CMD)

mocks:
	@echo "Generating mocks with mockery..."
	$(MOCKERY_CMD)

## Application Execution Rules
## ---------------------------

machine:
	@echo "Starting external-device machine..."
	$(MACHINE_CMD)