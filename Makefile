include .env

migrate-new:
ifndef seq
	$(error seq is required, usage: make migrate-create seq=name_of_migration)
endif
	@docker compose run --rm lab-track-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@docker compose run --rm lab-track-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@lab-track-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

swagger:
	cd src/backend && swag init -g cmd/main.go -o ./internal/docs