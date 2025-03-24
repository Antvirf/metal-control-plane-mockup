.PHONY: containers
containers:
	cd compose && docker compose up --build


.PHONY: worker
worker:
	go run main.go -mode worker

.PHONY: control-plane
control-plane:
	go run main.go -mode controlplane

.PHONY: sqlc
sqlc:
	cd internal/ && sqlc generate


.PHONY: migrate
migrate:
	cat ./internal/schema.sql | docker exec cp_postgres psql -U controlplane -d controlplane -c
