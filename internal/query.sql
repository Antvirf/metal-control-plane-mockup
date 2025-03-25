-- name: GetHardwareInfoByBmcMacAddress :one
SELECT * from hardwareinfo
WHERE bmcMac = $1 LIMIT 1;

-- name: GetServerEthernetInterfaceMacAddressesByBmcMacAddress :many
SELECT 
  interface->>'MACAddress'
FROM                                       
  hardwareinfo,
  jsonb_array_elements(info->'RedFishData'->'EthernetInterfaces') AS interface
WHERE bmcMac = $1;

-- name: GetHardwareInfoByEthernetInterfaceMacAddresses :one
SELECT bmcMac, info
FROM hardwareinfo
WHERE EXISTS (
    SELECT 1
    FROM jsonb_array_elements(info->'RedFishData'->'EthernetInterfaces') AS interface
    WHERE LOWER(interface->>'MACAddress') = LOWER($1)
);

-- name: CreateHardwareInfo :one
INSERT INTO hardwareinfo (bmcMac, info)
VALUES ($1, $2)
ON CONFLICT(bmcMac)
DO UPDATE SET info = EXCLUDED.info
RETURNING *;
