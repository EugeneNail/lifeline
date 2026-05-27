COMPOSE_FILE := docker/docker-compose.yaml
ENV_FILE := docker/.env
ENV_EXAMPLE_FILE := docker/.env.example
NETWORK_NAME := lifeline
MIGRATION_COMMAND_TARGET := $(word 2,$(MAKECMDGOALS))
MIGRATION_NAME := $(wordlist 3,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
SHELL_TARGET := $(word 2,$(MAKECMDGOALS))
GO_RUN_CONTAINER := docker run --rm -v "$(CURDIR)":/workspace -w /workspace golang:1.26.1
GO_RUN_NETWORK_CONTAINER := docker run --rm --network "$(NETWORK_NAME)" --env-file "$(ENV_FILE)" -v "$(CURDIR)":/workspace -w /workspace golang:1.26.1

.PHONY: up down shell env network create migrate rollback

up: network
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d --build

down:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down

shell:
	@if [ -z "$(SHELL_TARGET)" ]; then echo "usage: make shell <container-name-without-lifeline-prefix>"; exit 1; fi
	docker exec -it lifeline-$(SHELL_TARGET) sh

env:
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

network:
	@docker network inspect "$(NETWORK_NAME)" >/dev/null 2>&1 || docker network create "$(NETWORK_NAME)"

create:
	@if [ "$(MIGRATION_COMMAND_TARGET)" != "migration" ]; then echo "usage: make create migration <multi word name>"; exit 1; fi
	@if [ -z "$(strip $(MIGRATION_NAME))" ]; then echo "usage: make create migration <multi word name>"; exit 1; fi
	@$(GO_RUN_CONTAINER) go run ./cmd/migrator create "$(strip $(MIGRATION_NAME))"

migrate: env network
	@$(GO_RUN_NETWORK_CONTAINER) go run ./cmd/migrator migrate

rollback: env network
	@$(GO_RUN_NETWORK_CONTAINER) go run ./cmd/migrator rollback

%:
	@:
