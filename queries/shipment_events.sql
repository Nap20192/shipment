-- name: AddShipmentEvent :one
INSERT INTO
    shipment_events (
        id,
        shipment_id,
        status,
        description
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetShipmentEventHistory :many
SELECT *
FROM shipment_events
WHERE
    shipment_id = $1
ORDER BY created_at ASC;
