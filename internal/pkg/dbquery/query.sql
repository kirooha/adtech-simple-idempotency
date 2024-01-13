-- name: CreateCampaign :one
INSERT INTO campaigns (
    name, description
) VALUES (
    $1, $2
)
ON CONFLICT(name) DO UPDATE
    SET name = excluded.name
RETURNING *;

-- name: GetCampaign :one
SELECT *
FROM campaigns
WHERE id = $1;

-- name: UpdateCampaign :one
UPDATE campaigns
SET adserver_id = $1
WHERE id = $2
RETURNING *;