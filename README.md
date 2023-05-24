# go_movie
microservices based movie catalog developed using golang


## docker compose commands:

docker compose -f docker-compose.yml -f docker-compose.override.yml up -d

docker compose -f docker-compose.yml -f docker-compose.override.yml up -d --scale recipe_worker=5

docker-compose -f docker-compose.yml -f docker-compose.override.yml down

docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d --build

docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d --build  --scale recipe_worker=5

## docker - mysql commands:

docker exec -i db mysql movie -h localhost -P 3306 --protocol=tcp -uroot -pmauFJcuf5dhRMQrjj < schema/schema.sql

docker exec -i db mysql movie -h localhost -P 3306 --protocol=tcp -uroot -pmauFJcuf5dhRMQrjj -e "SHOW tables"

## grpc url commands:

grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie"}' localhost:8082 RatingService/GetAggregatedRating

grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie", "user_id": "alex", "rating_value": 5}' localhost:8082 RatingService/PutRating

grpcurl -plaintext -d '{"movie_id":"1"}' localhost:8083 MovieService/GetMovieDetails

## links:

Prometheus Dashboard: http://localhost:9090/
Alert Manager UI: http://localhost:9093
Portainer: http://localhost:9000/



