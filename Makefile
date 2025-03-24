.PHONY: containers
containers:
	cd compose && docker compose up --build


.PHONY: worker
worker:
	go run main.go -mode worker

.PHONY: control-plane
control-plane:
	go run main.go -mode controlplane