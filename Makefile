include .env.example.local
export $(shell sed 's/=.*//' .env.example.local)

# ==============================================================================
# Docker commands

FILES := $(shell docker ps -aq)

.PHONY: up
up:
	docker-compose -f docker-compose.yaml up -d --build && \
 	sleep 3; migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5450/postgres?sslmode=disable' up

.PHONY: down
down:
	docker-compose -f docker-compose.yaml down

.PHONY: local-up
local-up:
	docker-compose -f docker-compose.local.yaml up -d --build  && \
	sleep 3; migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5450/postgres?sslmode=disable' up


.PHONY: local-down
local-down:
	docker-compose -f docker-compose.local.yaml down

.PHONY: clean
docker-clean:
	docker system prune -f
	docker stop $(FILES)
	docker rm $(FILES)


# ==============================================================================
# Golang migrate postgresql

.PHONY: migrate-up
migrate-up:
		migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5450/postgres?sslmode=disable' up

.PHONY: migrate-down
migrate-down:
		migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5450/postgres?sslmode=disable' down

.PHONY: migrate-drop
migrate-drop:
		migrate -path ./migrations -database 'postgres://postgres:qwerty@0.0.0.0:5450/postgres?sslmode=disable' drop -f

# ==============================================================================
# Main

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build main.go

.PHONY: test
test:
	go test -cover ./... -count=1

.PHONY: test-short
test-short:
	go test -cover -short ./... -count=1

.PHONY: gen
gen:
	go generate ./... 2>/dev/null


# ==============================================================================
# Tools commands

.PHONY: linter
linter:
	golangci-lint run ./... --config=.golangci.yaml

.PHONY: linter-fix
linter-fix:
	golangci-lint run ./... --config=.golangci.yaml --fix


# ==============================================================================
# Modules support

.PHONY: deps-reset
deps-reset:
	git checkout -- go.mod
	go mod tidy

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: deps-upgrade
deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

.PHONY: deps-cleancache
deps-cleancache:
	go clean -modcache
