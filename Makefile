COMPOSE_FILE := docker/docker-compose.yaml
ENV_FILE := docker/.env
ENV_EXAMPLE_FILE := docker/.env.example
HTTP_CLIENT_ENV_FILE := doc/routes/http-client.env.json
HTTP_CLIENT_ENV_EXAMPLE_FILE := doc/routes/http-client.env.json.example
NETWORK_NAME := lifeline
MIGRATION_COMMAND_TARGET := $(word 2,$(MAKECMDGOALS))
MIGRATION_NAME := $(wordlist 3,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
SHELL_TARGET := $(word 2,$(MAKECMDGOALS))
GO_RUN_CONTAINER := docker run --rm -v "$(CURDIR)":/workspace -w /workspace golang:1.26.1
GO_RUN_NETWORK_CONTAINER := docker run --rm --network "$(NETWORK_NAME)" --env-file "$(ENV_FILE)" -v "$(CURDIR)":/workspace -w /workspace golang:1.26.1

.PHONY: help up down shell envs networks fetch deploy create migrate rollback

help:
	@printf '%-14s %-28s %s\n' "Target" "Interface" "Description"
	@printf '%-14s %-28s %s\n' "------" "--------- " "-----------"
	@printf '%-14s %-28s %s\n' "help" "" "Show this help table."
	@printf '%-14s %-28s %s\n' "up" "" "Start the full stack with Docker Compose."
	@printf '%-14s %-28s %s\n' "down" "" "Stop the full stack."
	@printf '%-14s %-28s %s\n' "shell" "<container>" "Open a shell inside a running Lifeline container."
	@printf '%-14s %-28s %s\n' "envs" "" "Create or refresh local environment files from examples."
	@printf '%-14s %-28s %s\n' "networks" "" "Create the shared Docker network if it is missing."
	@printf '%-14s %-28s %s\n' "fetch" "" "Pull the latest master branch from git."
	@printf '%-14s %-28s %s\n' "deploy" "" "Fetch the latest code, refresh local env files, and start the stack."
	@printf '%-14s %-28s %s\n' "create" "migration <name>" "Create a new database migration."
	@printf '%-14s %-28s %s\n' "migrate" "" "Apply pending database migrations."
	@printf '%-14s %-28s %s\n' "rollback" "" "Roll back the latest database migration."

up: networks
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d --build

down:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down

shell:
	@if [ -z "$(SHELL_TARGET)" ]; then echo "usage: make shell <container-name-without-lifeline-prefix>"; exit 1; fi
	docker exec -it lifeline-$(SHELL_TARGET) sh

envs:
	@if [ ! -f "$(ENV_FILE)" ]; then \
		cp "$(ENV_EXAMPLE_FILE)" "$(ENV_FILE)"; \
		echo "created $(ENV_FILE)"; \
	else \
		while IFS= read -r line || [ -n "$$line" ]; do \
			case "$$line" in \
				''|\#*) continue ;; \
			esac; \
			key=$${line%%=*}; \
			if ! grep -Eq "^$${key}=" "$(ENV_FILE)"; then \
				printf '\n%s\n' "$$line" >> "$(ENV_FILE)"; \
			fi; \
		done < "$(ENV_EXAMPLE_FILE)"; \
	fi
	@if [ ! -f "$(HTTP_CLIENT_ENV_FILE)" ]; then \
		cp "$(HTTP_CLIENT_ENV_EXAMPLE_FILE)" "$(HTTP_CLIENT_ENV_FILE)"; \
		echo "created $(HTTP_CLIENT_ENV_FILE)"; \
	fi

networks:
	@docker network inspect "$(NETWORK_NAME)" >/dev/null 2>&1 || docker network create "$(NETWORK_NAME)"

fetch:
	git pull --ff-only origin main

deploy:
	@$(MAKE) fetch
	@$(MAKE) networks
	@$(MAKE) envs
	@$(MAKE) migrate
	@$(MAKE) up

create:
	@if [ "$(MIGRATION_COMMAND_TARGET)" != "migration" ]; then echo "usage: make create migration <multi word name>"; exit 1; fi
	@if [ -z "$(strip $(MIGRATION_NAME))" ]; then echo "usage: make create migration <multi word name>"; exit 1; fi
	@$(GO_RUN_CONTAINER) go run ./cmd/migrator create "$(strip $(MIGRATION_NAME))"

migrate: envs networks
	@$(GO_RUN_NETWORK_CONTAINER) go run ./cmd/migrator migrate

rollback: envs networks
	@$(GO_RUN_NETWORK_CONTAINER) go run ./cmd/migrator rollback

%:
	@:
