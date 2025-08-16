.PHONY: migrate-up migrate-down

migrate-up:
	goose -dir db/sqlite/migrations sqlite3 db/sqlite/kumo.db up

migrate-down:
	goose -dir db/sqlite/migrations sqlite3 db/sqlite/kumo.db down

generate-sqlc:
	sqlc generate