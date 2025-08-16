-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    status_id INTEGER NOT NULL,
    retries INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (status_id) REFERENCES task_status(id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE tasks;