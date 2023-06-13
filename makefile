# ==============================================================================
# variables
# ==============================================================================

VERSION := 1.0.0

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
# docker
# ==============================================================================

docker-build-metadata:
	sudo docker build \
    		-f docker/dockerfile/Dockerfile \
    		-t ashkanmaleki/metadata:$(VERSION) \
    		--build-arg service=metadata \
    		--build-arg BUILD_REF=$(VERSION) \
    		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    		.


docker-build: docker-build-metadata


# ==============================================================================
# docker compose
# ==============================================================================

compose-build: docker-build
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