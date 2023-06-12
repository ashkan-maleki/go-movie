# ==============================================================================
# run
# ==============================================================================

run-movie:
	go run movie/cmd/*.go

run-metadata:
	go run metadata/cmd/*.go

run-rating:
	go run rating/cmd/*.go


# ==============================================================================
# docker compose
# ==============================================================================

compose-up:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml up \
	-d

compose-infra-up:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml up \
	-d

compose-down:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml down

compose-infra-down:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml down