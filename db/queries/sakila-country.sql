-- name: CreateCountry :execresult
INSERT INTO country (country) VALUES (?);

-- name: GetCountry :one
SELECT * FROM country WHERE country_id = ? LIMIT 1;

-- name: ListCountries :many
SELECT * FROM country ORDER BY country_id;

-- name: UpdateCountry :exec
UPDATE country SET country = ? WHERE country_id = ?;

-- name: DeleteCountry :exec
DELETE FROM country WHERE country_id = ?;
