package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/noahgould/bandsfromtown/api"
	"github.com/noahgould/bandsfromtown/dal"
)

// func startWebServer(r *Router) {
// 	http.ListenAndServe(":8080", r)
// }

func main() {
	r := mux.NewRouter()
	db, err := dal.StartDB("noah:bigDBpass@tcp(localhost)/bandsfromtown")
	if err != nil {
		log.Fatal(err)
	}

	artistStore := dal.NewArtistStore(db)
	locationStore := dal.NewLocationStore(db)
	artistController := api.NewArtistController(artistStore, locationStore)

	r.HandleFunc("/artist/{artist}", artistController.LookupArtist)
	r.HandleFunc("/artist", artistController.Index)

	http.ListenAndServe(":8080", r)

	//startWebServer(artistController, r)

}
