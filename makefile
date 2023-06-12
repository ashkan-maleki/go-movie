# ==============================================================================
# docker compose
# ==============================================================================

compose-up:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml up \
	-d

compose-up-infra:
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

compose-down-infra:
	docker compose -f docker/compose/monitor/docker-compose.yaml \
	-f docker/compose/monitor/docker-compose.override.yaml \
	-f docker/compose/infra/docker-compose.yaml \
	-f docker/compose/infra/docker-compose.override.yaml down