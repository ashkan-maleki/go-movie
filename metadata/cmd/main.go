package main

import (
	"github.com/mamalmaleki/go_movie/metadata/internal/controller/metadata"
	httpHandler "github.com/mamalmaleki/go_movie/metadata/internal/handler/http"
	"github.com/mamalmaleki/go_movie/metadata/internal/repository/memory"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httpHandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
