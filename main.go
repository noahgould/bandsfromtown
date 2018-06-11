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
	spotifyController := api.NewSpotifyController(artistStore, locationStore)

	r.HandleFunc("/artist/{artist}", artistController.LookupArtist).Methods("GET", "OPTIONS")
	r.HandleFunc("/artist", artistController.Index)
	r.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))
	r.Handle("/", http.RedirectHandler("/frontend/artistLookup.html", 301))
	r.HandleFunc("/spotify/auth/", spotifyController.AuthorizationRequest)
	r.HandleFunc("/spotify/login/", spotifyController.AuthorizationCallback)
	r.HandleFunc("/spotify/locations/{accessToken}", spotifyController.FindUserArtistLocations)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Println(err)
	}
}
