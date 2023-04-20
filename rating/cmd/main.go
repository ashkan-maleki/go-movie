package main

import (
	"github.com/mamalmaleki/go_movie/rating/internal/controller/rating"
	httpHandler "github.com/mamalmaleki/go_movie/rating/internal/handler/http"
	"github.com/mamalmaleki/go_movie/rating/internal/repository/memory"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
