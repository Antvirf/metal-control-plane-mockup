.PHONY: containers
containers:
	cd compose && docker compose up --build


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
	cat ./internal/schema.sql | docker exec cp_postgres psql -U controlplane -d controlplane -c
