-- name: GetHardwareInfo :one
SELECT * from hardwareinfo
WHERE mac = $1 LIMIT 1;


-- name: CreateHardwareInfo :one
INSERT INTO hardwareinfo (mac, info)
VALUES ($1, $2) RETURNING *;
