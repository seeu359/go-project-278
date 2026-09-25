-- name: GetLinks :many
SELECT * FROM links;

-- name: GetLink :one
SELECT * FROM links WHERE id = $1;

-- name: CreateLink :one
INSERT INTO links (url, short_name)
VALUES ($1, $2)
RETURNING id;

-- name: UpdateLink :execresult
UPDATE links
SET url = $1,
    short_name = $2
WHERE id = $3;

-- name: DeleteLink :execresult
DELETE FROM links
WHERE id = $1;

-- name: GetVisits :many
SELECT * FROM visits;

-- name: AddVisit :one
INSERT INTO visits (
  link_id,
  ip,
  user_agent,
  status
) VALUES ($1, $2, $3, $4)
RETURNING id;
