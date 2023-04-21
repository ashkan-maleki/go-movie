# go_movie
microservices based movie catalog developed using golang


## commands:

docker compose -f docker-compose.yml -f docker-compose.override.yml up -d

docker compose -f docker-compose.yml -f docker-compose.override.yml up -d --scale recipe_worker=5

docker-compose -f docker-compose.yml -f docker-compose.override.yml down

docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d --build

docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d --build  --scale recipe_worker=5