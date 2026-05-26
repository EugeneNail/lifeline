COMPOSE_FILE := docker/docker-compose.yaml
ENV_FILE := docker/.env
ENV_EXAMPLE_FILE := docker/.env.example
NETWORK_NAME := lifeline
SHELL_TARGET := $(word 2,$(MAKECMDGOALS))

.PHONY: up down shell env network

up: network
	docker compose -f $(COMPOSE_FILE) up -d --build

down:
	docker compose -f $(COMPOSE_FILE) down

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

%:
	@:
