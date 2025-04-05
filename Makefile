.PHONY: migrate
migrate:
	goose -dir=migrations postgres "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable" up

downgrade:
	goose -dir=migrations postgres "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable" down

.PHONY: run
run:
	go run cmd/grats/main.go

.PHONY: format
format:
	gofmt -w .
