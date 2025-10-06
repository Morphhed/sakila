-- name: CreateCity :execresult
INSERT INTO city (city, country_id) VALUES (?, ?);

-- name: GetCity :one
SELECT * FROM city WHERE city_id = ? LIMIT 1;

-- name: ListCities :many
SELECT * FROM city ORDER BY city_id;

-- name: UpdateCity :exec
UPDATE city SET city = ?, country_id = ? WHERE city_id = ?;

-- name: DeleteCity :exec
DELETE FROM city WHERE city_id = ?;

-- name: ListCitiesByCountry :many
SELECT * FROM city WHERE country_id = ? ORDER BY city;
