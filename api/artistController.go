package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/noahgould/bandsfromtown/dal"
)

type ArtistController struct {
	artistStore   dal.ArtistStore
	locationStore dal.LocationStore
}

func NewArtistController(newArtistStore dal.ArtistStore, newLocationStore dal.LocationStore) *ArtistController {

	return &ArtistController{
		artistStore: newArtistStore,
	}
}

func (ac *ArtistController) Register() {
	http.HandleFunc("/artist", ac.lookupArtist)
}

func (ac *ArtistController) lookupArtist(w http.ResponseWriter, r *http.Request) {
	artistName := strings.SplitN(r.URL.Path, "/", 3)[2]

	artists, err := ac.artistStore.GetArtistsByName(artistName)

	if err != nil {
		log.Fatal(err)
	}

	if artists == nil {
		artistLocationStrings := strings.SplitN(LookupArtistLocation(artistName), " ", 3)
		artistLocation := dal.Location{
			City:    artistLocationStrings[0],
			State:   artistLocationStrings[1],
			Country: artistLocationStrings[2],
		}

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
