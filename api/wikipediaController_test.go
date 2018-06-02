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
		{"Porches", dal.Location{
			City: "Pleasantville", State: "New York", Country: "United States",
		}},
		{"Kali Uchis", dal.Location{
			City: "Alexandria", State: "unknown", Country: "Virginia",
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
			City: "Staten Island", State: "New York", Country: "U.S.",
		}},
		{"Jeff Rosenstock", dal.Location{
			City: "Long Island", State: "unknown", Country: "United States",
		}},
	}

	for _, testArtist := range table {
		location := LookupArtistLocation(testArtist.artistName)

		if location.City != testArtist.location.City || location.State != testArtist.location.State || location.Country != testArtist.location.Country {
			t.Errorf("Location incorrect, got: %s, %s, %s. Want: %s, %s, %s.", location.City, location.State, location.Country, testArtist.location.City, testArtist.location.State, testArtist.location.Country)
		}
	}

}
