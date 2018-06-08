package api

import (
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
			locationPtr, err := gMC.NormalizeLocation(artistLocation)

			artistLocation = *locationPtr
			if err != nil {
				log.Println(err)
			}

			artistLocation = *ac.checkForExistingLocation(artistLocation)

			newArtist := dal.Artist{
				Name:     artistName,
				Location: artistLocation,
			}

			newArtist.ID, err = ac.artistStore.AddArtist(newArtist)
			artists = append(artists, newArtist)
		} else {
			for i, artist := range artists {
				artists[i].Location, err = ac.locationStore.GetLocationByID(artist.Location.ID)

				if artist.Location.Latitude == 0 && artist.Location.Longitude == 0 {
					gMC := NewGoogleMapsController()
					artistLocationPtr, err := gMC.GetCoordinates(artists[i].Location)

					if err != nil {
						log.Println(err)
					} else {
						artists[i].Location = *artistLocationPtr
						ac.locationStore.UpdateLocation(artists[i].Location)
					}
				}

				if err != nil {
					log.Println(err)
				}
			}
		}

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

	artistLocation, err := googleMapController.NormalizeLocation(newLocation)

	artistLocation = ac.checkForExistingLocation(*artistLocation)

	if artistLocation.Latitude == 0 {
		artistLocation, err = googleMapController.GetCoordinates(*artistLocation)
	}

	artistToUpdate, err := ac.artistStore.GetArtistByID(artistId)

	if err != nil {
		log.Printf("artistController line 103, artist: %d .err: %s", artistToUpdate.ID, err)
		log.Println(artistToUpdate)
	}

	artistToUpdate.Location = *artistLocation

	ac.artistStore.UpdateArtist(artistToUpdate)
}

func (ac *ArtistController) checkForExistingLocation(locationToCheck dal.Location) *dal.Location {
	alreadyStored, location := ac.locationStore.CheckForExistingLocation(locationToCheck)

	if !alreadyStored {
		gMC := NewGoogleMapsController()
		artistLocationPtr, err := gMC.GetCoordinates(locationToCheck)
		if err != nil {
			log.Println(err)
		} else {
			locationToCheck = *artistLocationPtr
			locationToCheck.ID, err = ac.locationStore.AddLocation(locationToCheck)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		locationToCheck = location
	}

	return &locationToCheck
}
