.PHONY: migrate-up migrate-down generate-sqlc migrate-create mocks

MIGRATION_DIR = db/sqlite/migrations
DATABASE_TYPE = sqlite3
DATABASE_URL = db/sqlite/kumo.db

migrate-up:
	goose -dir $(MIGRATION_DIR) $(DATABASE_TYPE) $(DATABASE_URL) up

migrate-down:
	goose -dir $(MIGRATION_DIR) $(DATABASE_TYPE) $(DATABASE_URL) down

generate-sqlc:
	sqlc generate

migrate-create:
	goose -dir $(MIGRATION_DIR) $(DATABASE_TYPE) $(DATABASE_URL) create $(NAME) sql

mocks:
	mockery