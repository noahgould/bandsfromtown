package api

import (
	"fmt"
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
		{Name: "Kali Uchis", SpotifyID: "1U1el3k54VvEUzo3ybLPlM"},
		{Name: "Foster the People"},
		{Name: "Glass Animals", SpotifyID: "4yvcSjfu4PC0CYQyLy4wSq"},
	}

	outputArtists := []dal.Artist{
		{Name: "Kali Uchis", Location: dal.Location{City: "Alexandria", State: "Virginia", Country: "USA", FullLocation: "Alexandria, Virginia, USA", GooglePlaceID: "ChIJ8aukkz5NtokRLAHB24Ym9dc"}, SpotifyID: "1U1el3k54VvEUzo3ybLPlM"},
		{Name: "Foster the People", Location: dal.Location{City: "Los Angeles", State: "California", Country: "U.S", FullLocation: "Los Angeles, California, U.S.", GooglePlaceID: "ChIJE9on3F3HwoAR9AhGJW_fL-I"}, SpotifyID: "7gP3bB2nilZXLfPHJhMdvc"},
		{Name: "Glass Animals", Location: dal.Location{City: "Oxford", State: "England", Country: "UK", FullLocation: "Oxford, England, UK", GooglePlaceID: "ChIJrx_ErYAzcUgRAnRUy6jbIMg"}, SpotifyID: "4yvcSjfu4PC0CYQyLy4wSq"},
	}

	artists := make(chan dal.Artist)

	go func() {
		for _, a := range inputArtists {
			artists <- a
		}
		fmt.Println("closing artists")
		close(artists)
	}()

	spotifyController := NewSpotifyController(artistStore, locationStore)

	artistResults := spotifyController.getArtistLocations(artists)

	for j, artistWithLocation := range artistResults {
		if artistWithLocation.FullLocation != outputArtists[j].FullLocation {
			t.Errorf("Location incorrect, got: %s. Want: %s.", artistWithLocation.FullLocation, outputArtists[j].FullLocation)
		}
		if artistWithLocation.Location.GooglePlaceID != outputArtists[j].Location.GooglePlaceID {
			t.Errorf("Google Place ID Mismatch., got: %s. Want: %s.", artistWithLocation.GooglePlaceID, outputArtists[j].GooglePlaceID)
		}
	}

}
