package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/noahgould/bandsfromtown/dal"
)

type ArtistController struct {
	artistStore dal.ArtistStore
}

func NewArtistController(newArtistStore dal.ArtistStore) *ArtistController {

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

	if err := json.NewEncoder(w).Encode(artists); err != nil {
		log.Fatal(err)
	}

}
