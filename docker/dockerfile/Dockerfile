ARG service
FROM golang:1.20-alpine AS builder
ARG service

WORKDIR /movie

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY ./ ./


# Build the binary.
RUN go build -ldflags="-s -w" -v -o service ./${service}/cmd/

ARG service
FROM alpine:latest AS runtime
ARG service

#COPY --from=builder /mallbots/docker/wait-for .
#RUN chmod +x /wait-for

# Copy the binary to the production image from the builder stage.
COPY --from=builder /movie/service /movie/service
#COPY ../${service}/configs/base.yaml /movie/base.yaml

CMD [ "/movie/service" ]
#CMD ["ls", "/movie/service"]
