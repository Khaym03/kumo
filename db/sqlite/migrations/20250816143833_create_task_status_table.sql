-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE task_status (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO task_status (id, name) VALUES
(1, 'PENDING'),
(2, 'IN_PROGRESS'),
(3, 'COMPLETED'),
(4, 'FAILED');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE task_status;