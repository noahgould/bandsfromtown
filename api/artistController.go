package api

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
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
	jsonOnlyString := mux.Vars(r)["jsonOnly"]
	var jsonOnly = false

	if jsonOnlyString == "t" || jsonOnlyString == "T" {
		jsonOnly = true
	}

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

		if jsonOnly {
			if err := json.NewEncoder(w).Encode(artists); err != nil {
				log.Println(err)
			}
		} else {
			t, err := template.ParseFiles("./frontend/artistLookup.html")
			if err != nil {
				log.Print("template parsing error: ", err)
			}
			err = t.Execute(w, artists[0])
			if err != nil {
				log.Print("template executing error: ", err)
			}
		}

	}

}

//UpdateArtistLocation Updates the location for an existing artist.
func (ac *ArtistController) UpdateArtistLocation(w http.ResponseWriter, r *http.Request) {
	artistID, err := strconv.Atoi(mux.Vars(r)["artistID"])

	if err != nil {
		log.Printf("artistController line 84, %s", err)
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Println(readErr)
	}

	newLocation := dal.Location{}

	jsonErr := json.Unmarshal(body, &newLocation)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	artistLocation := ac.checkForExistingLocation(newLocation)

	if artistLocation.Latitude == 0 {
		newLocation.ID = artistLocation.ID
		ac.locationStore.UpdateLocation(newLocation)
	}

	artistToUpdate, err := ac.artistStore.GetArtistByID(artistID)

	if err != nil {
		log.Printf("artistController line 103, artist: %d .err: %s", artistToUpdate.ID, err)
	}

	artistToUpdate.Location = *artistLocation

	artistToUpdate.ID, err = ac.artistStore.UpdateArtist(artistToUpdate)
	if err != nil {
		log.Println(err)
	} else {
		if err := json.NewEncoder(w).Encode(artistToUpdate); err != nil {
			log.Fatal(err)
		}
	}
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
