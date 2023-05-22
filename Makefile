

.PHONY: run generate fmt

run:
	docker-compose -f docker-compose.yml up --build -d

logs:
	docker-compose -f docker-compose.yml logs -f

fmt: fmt-docs

fmt-docs:
	swag fmt

generate: generate-sqlc

generate-docs:
	swag init -g .\cmd\ttto\main.go

generate-sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate
