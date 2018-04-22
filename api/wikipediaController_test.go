package api

import (
	"testing"

	"github.com/noahgould/bandsfromtown/dal"
)

func TestLookupArtistLocation(t *testing.T) {
	table := []struct {
		artistName string
		location   dal.Location
	}{
		{"Kali Uchis", dal.Location{
			City: "Pereira", State: "unknown", Country: "Colombia",
		}},
		{"The Shins", dal.Location{
			City: "Albuquerque", State: "New Mexico", Country: "United States",
		}},
		{"Unknown Mortal Orchestra", dal.Location{
			City: "Auckland", State: "unknown", Country: "New Zealand",
		}},
		{"King_Krule", dal.Location{
			City: "Southwark", State: "London", Country: "England",
		}},
		{"Ghostface Killah", dal.Location{
			City: "Southwark", State: "London", Country: "England",
		}},
	}

	for _, testArtist := range table {
		location := LookupArtistLocation(testArtist.artistName)

		if location.City != testArtist.location.City || location.State != testArtist.location.State || location.Country != testArtist.location.Country {
			t.Errorf("Location incorrect, got: %s, %s, %s. Want: %s, %s, %s.", location.City, location.State, location.Country, testArtist.location.City, testArtist.location.State, testArtist.location.Country)
		}
	}

}
