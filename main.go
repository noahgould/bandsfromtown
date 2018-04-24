package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/noahgould/bandsfromtown/api"
	"github.com/noahgould/bandsfromtown/dal"
)

func main() {
	r := mux.NewRouter()
	db, err := dal.StartDB(os.Getenv("CLEARDB_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	artistStore := dal.NewArtistStore(db)
	locationStore := dal.NewLocationStore(db)
	artistController := api.NewArtistController(artistStore, locationStore)

	r.HandleFunc("/artist/{artist}", artistController.LookupArtist)
	r.HandleFunc("/artist", artistController.Index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, r)

}
