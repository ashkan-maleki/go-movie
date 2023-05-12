FROM golang:1.20
WORKDIR /go/src/github.com/mamalmaleki/go_movie
COPY go.sum .
COPY go.mod .
RUN go mod download

COPY . .
#RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/*


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/mamalmaleki/go_movie/app .
#CMD ["./app"]