-- name: CreateShipment :one
INSERT INTO
    shipments (
        id,
        origin,
        destination,
        status,
        cost,
        revenue,
        weight,
        dimension_length,
        dimension_width,
        dimension_height,
        driver_name
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11
    ) RETURNING *;

-- name: GetShipmentByID :one
SELECT * FROM shipments WHERE id = $1;

-- name: UpdateShipmentStatus :one
UPDATE shipments
SET status = $2
WHERE id = $1
RETURNING *;
