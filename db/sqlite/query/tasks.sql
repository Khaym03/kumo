-- name: GetStatusIDByName :one
SELECT id FROM task_status WHERE name = ?;

-- name: GetTaskStatus :many
SELECT *
FROM task_status;

-- name: AddTask :exec
INSERT INTO tasks (url, status_id)
VALUES (?, ?);

-- name: ListTasksByStatusID :many
SELECT * FROM tasks WHERE status_id = ?;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status_id = ?
WHERE id = ?;

-- name: FailTask :exec
UPDATE tasks
SET status_id = (
    SELECT id FROM task_status WHERE name = 'FAILED'
), retries = retries + 1
WHERE tasks.id = ?;

-- name: ListPendingOrFailedTasks :many
SELECT *
FROM tasks
WHERE status_id IN (
    SELECT id FROM task_status WHERE name IN ('PENDING', 'FAILED')
)
AND retries <= 3
ORDER BY created_at;
