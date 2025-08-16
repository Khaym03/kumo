-- name: GetStatusIDByName :one
SELECT id FROM task_status WHERE name = ?;

-- name: AddTask :exec
INSERT INTO tasks (id, url, status_id)
VALUES (?, ?, ?);

-- name: GetTaskForProcessing :one
UPDATE tasks
SET status_id = ?
WHERE id IN (
    SELECT id
    FROM tasks AS t
    WHERE t.status_id = ?
    ORDER BY t.created_at
    LIMIT 1
)
RETURNING *;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status_id = ?
WHERE id = ?;

-- name: FailTask :exec
UPDATE tasks
SET status_id = ?, retries = retries + 1
WHERE id = ?;

-- name: ListPendingOrFailedTasks :many
SELECT *
FROM tasks
WHERE status_id IN ( ?, ? ) AND retries < 3
ORDER BY created_at;