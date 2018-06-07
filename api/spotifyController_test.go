package api

import (
	"log"
	"os"
	"testing"

	"github.com/noahgould/bandsfromtown/dal"
)

func TestGetArtistLocations(t *testing.T) {
	db, err := dal.StartDB(os.Getenv("LOCAL_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	artistStore := dal.NewArtistStore(db)
	locationStore := dal.NewLocationStore(db)

	inputArtists := []dal.Artist{
		{Name: "Kali Uchis"},
		{Name: "Foster the People"},
	}
	outputArtists := []dal.Artist{
		{Name: "Kali Uchis", Location: dal.Location{City: "Alexandria", State: "Virginia", Country: "USA", FullLocation: "Alexandria, Virginia, USA.", GooglePlaceID: "ChIJ8aukkz5NtokRLAHB24Ym9dc"}},
		{Name: "Foster the People"},
	}

	spotifyController := NewSpotifyController(artistStore, locationStore)

	artistResults := spotifyController.getArtistLocations(inputArtists)

	for j, artistWithLocation := range artistResults {
		if artistWithLocation.FullLocation != outputArtists[j].FullLocation {
			t.Errorf("Location incorrect, got: %s. Want: %s.", artistWithLocation.FullLocation, outputArtists[j].FullLocation)
		}
		if artistWithLocation.Location.GooglePlaceID != outputArtists[j].Location.GooglePlaceID {
			t.Errorf("Google Place ID Mismatch., got: %s. Want: %s.", artistWithLocation.GooglePlaceID, outputArtists[j].GooglePlaceID)
		}
	}

}
