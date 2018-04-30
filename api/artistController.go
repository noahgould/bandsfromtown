package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (ac *ArtistController) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func parseArtistName(name string) string {
	parsedName := strings.Replace(name, "%20", " ", -1)
	parsedName = strings.Replace(parsedName, "_", " ", -1)
	parsedName = strings.Title(parsedName)
	return parsedName
}

func (ac *ArtistController) LookupArtist(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	artistName := parseArtistName(mux.Vars(r)["artist"])
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

			if err != nil {
				log.Print(err)
			}

			artistLocation = *ac.checkForExistingLocation(artistLocation.GooglePlaceID)

			newArtist := dal.Artist{
				Name:     artistName,
				Location: artistLocation,
			}

			newArtist.ID, err = ac.artistStore.AddArtist(newArtist)
			artists = append(artists, newArtist)
		}

		log.Println(json.Marshal(artists[0]))
		log.Println("Location:")
		log.Println(json.Marshal(artists[0].Location))

		for i, artist := range artists {
			artists[i].Location, err = ac.locationStore.GetLocationByID(artist.Location.ID)
			if err != nil {
				log.Println(err)
			}
		}

		log.Println(json.Marshal(artists[0]))
		log.Println("Location:")
		log.Println(json.Marshal(artists[0].Location))

		if err := json.NewEncoder(w).Encode(artists); err != nil {
			log.Fatal(err)
		}
	}

}

func (ac *ArtistController) UpdateArtistLocation(w http.ResponseWriter, r *http.Request) {
	artistId, err := strconv.Atoi(mux.Vars(r)["artistID"])
	if err != nil {
		log.Printf("artistController line 84, %s", err)
	}
	newLocationString := mux.Vars(r)["location"]
	newLocationArray := strings.Split(newLocationString, ",")
	newLocation := dal.Location{
		City:    newLocationArray[0],
		State:   newLocationArray[1],
		Country: newLocationArray[2],
	}

	googleMapController := NewGoogleMapsController()

	artistLocation := googleMapController.NormalizeLocation(newLocation)

	artistLocation = ac.checkForExistingLocation(newLocation.GooglePlaceID)

	artistToUpdate, err := ac.artistStore.GetArtistByID(artistId)

	if err != nil {
		log.Printf("artistController line 103, artist: %d .err: %s", artistToUpdate.ID, err)
		log.Print(artistToUpdate)
	}

	artistToUpdate.Location = *artistLocation

	ac.artistStore.UpdateArtist(artistToUpdate)
}

func (ac *ArtistController) checkForExistingLocation(locationGoogleID string) *dal.Location {
	locationAlreadyStored, err := ac.locationStore.GetLocationByGoogleID(locationGoogleID)

	var artistLocation dal.Location
	if err != nil {
		if err == sql.ErrNoRows {
			artistLocation.ID, err = ac.locationStore.AddLocation(artistLocation)
		} else {
			log.Fatal(err)
		}
	} else {
		artistLocation = locationAlreadyStored
	}

	return &artistLocation

}
