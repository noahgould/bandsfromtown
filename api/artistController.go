package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
}

func (ac *ArtistController) LookupArtist(w http.ResponseWriter, r *http.Request) {

	artistName := strings.Title(mux.Vars(r)["artist"])
	if artistName == "" {
		w.Write([]byte("No artist entered."))
	} else {

		artists, err := ac.artistStore.GetArtistsByName(artistName)

		if err != nil {
			log.Fatal(err)
		}

		if artists == nil {
			artistLocation := LookupArtistLocation(artistName)
			gMC := NewGoogleMapsController()
			artistLocation = *gMC.NormalizeLocation(artistLocation)

			locationAlreadyStored, err := ac.locationStore.GetLocationByGoogleID(artistLocation.GooglePlaceID)

			if err != nil {
				if err == sql.ErrNoRows {
					artistLocation.ID, err = ac.locationStore.AddLocation(artistLocation)
				} else {
					log.Fatal(err)
				}
			} else {
				artistLocation = locationAlreadyStored
			}

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

		for i, artist := range artists {
			artists[i].Location, err = ac.locationStore.GetLocationByID(artist.Location.ID)
			if err != nil {
				log.Fatal(err)
			}
		}

		if err := json.NewEncoder(w).Encode(artists); err != nil {
			log.Fatal(err)
		}
	}

}
