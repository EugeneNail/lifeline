COMPOSE_FILE := docker/docker-compose.yaml
SHELL_TARGET := $(word 2,$(MAKECMDGOALS))

.PHONY: up down shell

up:
	docker compose -f $(COMPOSE_FILE) up -d --build

down:
	docker compose -f $(COMPOSE_FILE) down

shell:
	@if [ -z "$(SHELL_TARGET)" ]; then echo "usage: make shell <container-name-without-lifeline-prefix>"; exit 1; fi
	docker exec -it lifeline-$(SHELL_TARGET) sh

%:
	@:
