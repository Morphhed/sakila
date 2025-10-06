-- name: ListActors :many
SELECT actor_id, first_name, last_name, last_update FROM actor;

-- name: GetActor :one
SELECT actor_id, first_name, last_name, last_update FROM actor WHERE actor_id = ? LIMIT 1;

-- name: CreateActor :execresult
INSERT INTO actor (first_name, last_name) VALUES (?, ?);

-- name: UpdateActor :exec
UPDATE actor
SET first_name = ?, last_name = ?
WHERE actor_id = ?;

-- name: DeleteActor :exec
DELETE FROM actor WHERE actor_id = ?;

