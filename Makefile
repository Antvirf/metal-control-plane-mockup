.PHONY: containers
containers:
	cd compose && docker compose up --build


.PHONY: wipe
wipe:
	cd compose && docker compose down --volumes

.PHONY: worker
worker:
	go run main.go -mode worker

.PHONY: control-plane
control-plane:
	go run main.go -mode controlplane

.PHONY: generate
generate:
	cd internal/ && sqlc generate
	cd internal/pixieapi && ~/go/bin/stringer -type=ServerType


.PHONY: migrate
migrate:
	docker exec cp_postgres bash -c 'psql -U controlplane -d controlplane -f /mnt/internal/schema.sql'
