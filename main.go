package main

import (
	"log"
	"net/http"

	"github.com/noahgould/bandsfromtown/api"
	"github.com/noahgould/bandsfromtown/dal"
)

func startWebServer() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func main() {

	db, err := dal.StartDB("noah:bigDBpass@tcp(localhost)/bandsfromtown")
	if err != nil {
		log.Fatal(err)
	}

	artistStore := dal.NewArtistStore(db)
	artistController := api.NewArtistController(artistStore)
	artistController.Register()

	startWebServer()

}
