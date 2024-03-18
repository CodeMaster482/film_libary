include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/app/app.go
.PHONY: swag-v1

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

docker-rm-volume: ### remove docker volume
	docker volume rm film_libary_pg-data
.PHONY: docker-rm-volume

docker-it-db:
	docker exec -it postgres psql -U CodeMaster482 -d FLibraryDB 
.PHONY: docker-it-db

lint: linter-golangci linter-hadolint linter-dotenv ### run all linters
.PHONY: lint

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

mock: ### run mockgen
	~/go/bin/mockgen -source=./internal/actor/actor.go -destination=./internal/actor/mocks/mocks.go
	~/go/bin/mockgen -source=./internal/film/film.go -destination=./internal/film/mocks/mocks.go
.PHONY: mock

easyjson: ### run easyjson
	~/go/bin/easyjson -all internal/model/actor.go
	~/go/bin/easyjson -all internal/model/film.go
	~/go/bin/easyjson -all pkg/response/response.go
.PHONY: easyjson

bin-dep:
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest
	GOBIN=$(LOCAL_BIN) go install github.com/mailru/easyjson/...@latest
.PHONY: bin-dep
