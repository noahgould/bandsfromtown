package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noahgould/bandsfromtown/dal"
)

type ArtistController struct {
	artistStore   dal.ArtistStore
	locationStore dal.LocationStore
}

func NewArtistController(newArtistStore dal.ArtistStore, newLocationStore dal.LocationStore) *ArtistController {

	return &ArtistController{
		artistStore:   newArtistStore,
		locationStore: newLocationStore,
	}
}

func (ac *ArtistController) Register() {
	http.HandleFunc("/artist", ac.LookupArtist)
	fmt.Println("register reached.")
}

func (ac *ArtistController) LookupArtist(w http.ResponseWriter, r *http.Request) {

	artistName := mux.Vars(r)["artist"]

	if artistName == "" {
		w.Write([]byte("No artist entered."))
	} else {

		artists, err := ac.artistStore.GetArtistsByName(artistName)
		fmt.Println("post db call reached.")

		if err != nil {
			log.Fatal(err)
		}

		if artists == nil {
			artistLocation := LookupArtistLocation(artistName)

			artistLocation.ID, err = ac.locationStore.AddLocation(artistLocation)

			if err != nil {
				log.Fatal(err)
			}

			newArtist := dal.Artist{
				Name:     artistName,
				Location: artistLocation,
			}

			newArtist.ID, err = ac.artistStore.AddArtist(newArtist)
			artists = append(artists, newArtist)
		}

		if err := json.NewEncoder(w).Encode(artists); err != nil {
			log.Fatal(err)
		}
	}

}
