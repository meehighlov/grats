.PHONY: migrate
migrate:
	goose -dir=migrations sqlite3 grats.db up

.PHONY: run
run:
	go run cmd/grats/main.go

.PHONY: ea  # escape snalysis
ea:
	go run -gcflags=-m cmd/grats/main.go