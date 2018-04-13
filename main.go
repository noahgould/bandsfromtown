package main

import (
	"log"
	"net/http"

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

	startWebServer()

	db, err := dal.StartDB("noah:bigDBpass@tcp(localhost)/bandsfromtown")
	if err != nil {
		log.Fatal(err)
	}

	newLocation := dal.Location{
		State:   "Nebraska",
		City:    "Omaha",
		Country: "U>S>A"}

	newArtist := dal.Artist{
		Name:     "cheeseballs",
		Location: newLocation,
		Genre:    "good"}

	artistStore := dal.NewArtistStore(db)
	locationStore := dal.NewLocationStore(db)

	newLocation.ID, err = locationStore.AddLocation(newLocation)
	newArtist.Location.ID = newLocation.ID

	if err != nil {
		log.Fatal(err)
	}

	newArtist.ID, err = artistStore.AddArtist(newArtist)

	if err != nil {
		log.Fatal(err)
	}

}
