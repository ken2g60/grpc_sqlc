# Paths to Docker Compose files
LOCAL_COMPOSE_FILE := "./docker-compose.yml"

# Docker Compose command function
define dc_command
docker compose -f $(1) $(2)
endef

# Ensure Docker networks are created
create_local_network:
	docker network inspect microservices_nginx_network >/dev/null 2>&1 || docker network create microservices_nginx_network

# Common commands
up: create_local_network
	$(call dc_command,$(LOCAL_COMPOSE_FILE),up -d --build)

down:
	$(call dc_command,$(LOCAL_COMPOSE_FILE),down -v --remove-orphans)

logs:
	$(call dc_command,$(LOCAL_COMPOSE_FILE),logs -f)

ps:
	$(call dc_command,$(LOCAL_COMPOSE_FILE),ps)


# Migration commands
MIGRATE_CMD := docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate:v4.17.0 -path=/migrations -database "postgres://grpc_sqlc:s3cr3t$password@localhost:5432/grpc_sqlc_db?sslmode=disable"

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down 1

migrate-version:
	$(MIGRATE_CMD) version

# Create a new migration - usage: make migrate-create name=add_users_table
migrate-create:
	@mkdir -p migrations
	@docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate:v4.17.0 create -ext sql -dir /migrations -seq $(name)
	@echo "Created migration files for: $(name)"

# Generate sqlc code
sqlc-generate:
	sqlc generate

.PHONY: up down logs ps migrate-up migrate-down migrate-version migrate-create sqlc-generate