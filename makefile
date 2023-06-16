SHELL := /bin/bash
# ==============================================================================
# variables
# ==============================================================================

METADATA_VERSION := 1.0.0

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
	go run metadata/cmd/main.go | zap-pretty --all

run-rating:
	INFRA_CONFIG_FILE=configs/base.yaml APP_CONFIG_FILE=metadata/configs/base.yaml \
	go run rating/cmd/*.go

# ==============================================================================
# docker
# ==============================================================================

docker-build-metadata:
	sudo docker build \
    		-f docker/dockerfile/Dockerfile \
    		-t ashkanmaleki/metadata:$(METADATA_VERSION) \
    		--build-arg service=metadata \
    		--build-arg BUILD_REF=$(METADATA_VERSION) \
    		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    		.

docker-build: docker-build-metadata

docker-logs-metadata:
	sudo docker logs metadata -f --tail=100  | zap-pretty --all

# ==============================================================================
# docker compose
# ==============================================================================

compose-env-update:
	cd docker/compose/app; echo "METADATA_VERSION=$(METADATA_VERSION)" > .env

compose-up: compose-env-update
	docker compose -f docker/compose/app/docker-compose.yaml \
	-f docker/compose/app/docker-compose.override.yaml \
	--env-file docker/compose/app/.env \
	 up -d

compose-down:
	docker compose -f docker/compose/app/docker-compose.yaml \
	-f docker/compose/app/docker-compose.override.yaml \
	--env-file docker/compose/app/.env \
	down

compose-build-metadata: docker-build-metadata compose-down compose-up
compose-rebuild-metadata: compose-stop docker-build-metadata compose-start

compose-monitor-up:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	up -d

compose-monitor-down:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
  	-f docker/compose/monitor/docker-compose.override.yaml \
  	down


compose-infra-up:
	docker compose -f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml \
	up -d

compose-infra-down:
	docker compose -f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml \
	down

compose-start: compose-monitor-up compose-infra-up compose-up
compose-stop: compose-down compose-infra-down
compose-restart: compose-stop compose-start
compose-shutdown: compose-monitor-down compose-stop

