# Database
include .env
export

# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/backend/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s | tr A-Z a-z)
export ARCH := $(shell uname -m)

# --- Tooling & Variables ----------------------------------------------------------------
# include ./misc/make/tools.Makefile
# include ./misc/make/help.Makefile

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

up: dev-env             ## Startup / Spinup Docker Compose
down: docker-stop               ## Stop Docker
destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)

install-deps: ## Install Development Dependencies (localy).
	brew install golang-migrate
	go install github.com/cosmtrek/air@latest

dev-env: ## Bootstrap Environment (with a Docker-Compose help).
	@ docker-compose up -d

dev-air: ## Starts AIR ( Continuous Development app).
	cd backend && air

docker-stop:
	@ docker-compose down

docker-teardown:
	@ docker-compose down --remove-orphans -v

# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

lint: ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	cd backend && golangci-lint run ./...

build: ## Builds binary
	@ printf "Building aplication... "
	@ cd backend && go build \
		-trimpath  \
		-o engine \
		./main.go
	@ echo "done"

tests:
	@ cd backend && go test ./... -race -v

# ~~~ Database Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

migrate-up: ## Apply all migrations.
	migrate -database $(DATABASE_URL) -path=backend/db/migrations up

.PHONY: migrate-down
migrate-down: ## Rollback 1 migration.
	migrate -database $(DATABASE_URL) -path=backend/db/migrations down 1

.PHONY: migrate-drop
migrate-drop: ## Drop everything inside the database.
	migrate -database $(DATABASE_URL) -path=backend/db/migrations drop

.PHONY: migrate-create
migrate-create: ## Create a set of up/down migrations with a specified name.
	migrate create -ext sql -dir backend/db/migrations -seq $(name)

# ~~~ Cleans ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-artifacts clean-docker

clean-artifacts: ## Removes Artifacts (*.out)
	@printf "Cleanning artifacts... "
	@rm -f backend/*.out
	@echo "done."

clean-docker: ## Removes dangling docker images
	@ docker image prune -f
