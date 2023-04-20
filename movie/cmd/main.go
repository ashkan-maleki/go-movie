package main

import (
	"github.com/mamalmaleki/go_movie/movie/internal/controller/movie"
	metadataGatewayPkg "github.com/mamalmaleki/go_movie/movie/internal/gateway/metadata/http"
	ratingGatewayPkg "github.com/mamalmaleki/go_movie/movie/internal/gateway/rating/http"
	httpHandler "github.com/mamalmaleki/go_movie/movie/internal/handler/http"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting the movie service")
	metadataGateway := metadataGatewayPkg.New("localhost:8081")
	ratingGateway := ratingGatewayPkg.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
