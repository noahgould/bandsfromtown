package api

import (
	"testing"

	"github.com/noahgould/bandsfromtown/dal"
)

func TestNormalizeLocation(t *testing.T) {
	table := []struct {
		inputLocation  dal.Location
		outputlocation dal.Location
	}{
		{dal.Location{
			City: "Pereira", State: "unknown", Country: "Colombia",
		}, dal.Location{
			City: "Pereira", State: "Risaralda", Country: "Colombia", FullLocation: "Pereira, Risaralda, Colombia", GooglePlaceID: "ChIJ_cFW60iHOI4RvN_x-RCUs5U",
		}},
		{dal.Location{
			City: "Albuquerque", State: "New Mexico", Country: "United States",
		}, dal.Location{
			City: "Albuquerque", State: "New Mexico", Country: "United States", FullLocation: "Albuquerque, New Mexico, United States", GooglePlaceID: "ChIJe4MJ090KIocR_fbZuM7408A",
		}},
		{dal.Location{
			City: "Auckland", State: "unknown", Country: "New Zealand",
		}, dal.Location{
			City: "Auckland", State: "", Country: "New Zealand", FullLocation: "Auckland, New Zealand", GooglePlaceID: "ChIJ--acWvtHDW0RF5miQ2HvAAU",
		}},
		{dal.Location{
			City: "Southwark", State: "London", Country: "England",
		}, dal.Location{
			City: "London", State: "England", Country: "UK", GooglePlaceID: "ChIJdd4hrwug2EcRmSrV3Vo6llI", FullLocation: "London, England, UK",
		}},
		{dal.Location{
			City: "Staten Island", State: " New York City", Country: "New York (state)",
		}, dal.Location{
			City: "nil", State: "nil", Country: "nil", GooglePlaceID: "-1", FullLocation: "location could not be found",
		}},
	}

	gMController := NewGoogleMapsController()

	for _, location := range table {

		locationResult := gMController.NormalizeLocation(location.inputLocation)

		if locationResult.FullLocation != location.outputlocation.FullLocation || locationResult.GooglePlaceID != location.outputlocation.GooglePlaceID {
			t.Errorf("Location incorrect, got: %s, %s. Want: %s, %s.", locationResult.GooglePlaceID, locationResult.FullLocation, location.outputlocation.GooglePlaceID, location.outputlocation.FullLocation)

		}

	}
}
