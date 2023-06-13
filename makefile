# ==============================================================================
# run
# ==============================================================================

tidy:
	go mod tidy
	go mod vendor

run-movie:
	INFRA_CONFIG_FILE=configs/base.yaml APP_CONFIG_FILE=metadata/configs/base.yaml \
	go run movie/cmd/*.go

run-metadata:
	INFRA_CONFIG_FILE=configs/base.yaml APP_CONFIG_FILE=metadata/configs/base.yaml \
	go run metadata/cmd/main.go

run-rating:
	INFRA_CONFIG_FILE=configs/base.yaml APP_CONFIG_FILE=metadata/configs/base.yaml \
	go run rating/cmd/*.go


# ==============================================================================
# docker compose
# ==============================================================================

compose-build:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml \
	-f docker/compose/app/docker-compose.yaml \
	-f docker/compose/app/docker-compose.override.yaml \
	 up -d --build

compose-up:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml \
	-f docker/compose/app/docker-compose.yaml \
	-f docker/compose/app/docker-compose.override.yaml \
	 up -d

compose-infra-up:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml \
	up -d

compose-down:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml \
	-f docker/compose/app/docker-compose.yaml \
	-f docker/compose/app/docker-compose.override.yaml \
	down

compose-infra-down:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml \
	down