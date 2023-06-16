# go-movie
microservices based movie catalog developed using golang.<br/>
1) Metadata service: It stores and retrieves movie metadata<br/>
gRPC, MySQL, Go
2) Rating service: It stores and retrieves movie rating <br/>
gRPC, MySQL, Go
3) Movie service: It is an API Gateway. It sits between Client and other services.<br/>
gRPC, MySQL, Go

# requirements:
Golang, Docker, Docker compose, Makefile, and zap-pretty

# How to run:
1) make docker-build<br/>
2) make compose-start<br/>

## grpc url commands:

grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie"}' localhost:8082 RatingService/GetAggregatedRating

grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie", "user_id": "alex", "rating_value": 5}' localhost:8082 RatingService/PutRating

grpcurl -plaintext -d '{"movie_id":"1"}' localhost:8083 MovieService/GetMovieDetails

## other commands:
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=5

## links:

Prometheus Dashboard: http://localhost:9090/ <br/>
Alert Manager UI: http://localhost:9093 <br/>
Portainer: http://localhost:9000/ <br/>
Grafana: http://localhost:3000 <br/>



