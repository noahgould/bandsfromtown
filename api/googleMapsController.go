package api

import (
	"context"
	"log"
	"strings"

	"github.com/noahgould/bandsfromtown/dal"
	"googlemaps.github.io/maps"
)

type GoogleMapsController struct {
	mapsClient *maps.Client
}

func NewGoogleMapsController() *GoogleMapsController {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyA8oJVjkQZenxQIvA0EMXBAomiYjJwEqRE"))
	if err != nil {
		log.Fatal(err)
	}

	return &GoogleMapsController{
		mapsClient: c,
	}
}

func (gmc *GoogleMapsController) NormalizeLocation(location dal.Location) *dal.Location {

	var locationString string

	if location.State == "unknown" {
		locationString = strings.Join([]string{location.City, location.Country}, ",")
	} else {
		locationString = strings.Join([]string{location.City, location.State, location.Country}, ",")
	}

	place := &maps.PlaceAutocompleteRequest{
		Input: locationString,
		Types: maps.AutocompletePlaceTypeRegions,
	}

	placeResult, err := gmc.mapsClient.PlaceAutocomplete(context.Background(), place)

	if err != nil {
		log.Fatal(err)
	}

	normalizedLocation := &dal.Location{
		FullLocation:  placeResult.Predictions[0].Description,
		GooglePlaceID: placeResult.Predictions[0].PlaceID,
	}

	if len(placeResult.Predictions[0].Terms) < 3 {
		normalizedLocation.City = placeResult.Predictions[0].Terms[0].Value
		normalizedLocation.State = ""
		normalizedLocation.Country = placeResult.Predictions[0].Terms[1].Value
	} else {
		normalizedLocation.City = placeResult.Predictions[0].Terms[0].Value
		normalizedLocation.State = placeResult.Predictions[0].Terms[1].Value
		normalizedLocation.Country = placeResult.Predictions[0].Terms[2].Value
	}

	return normalizedLocation
}
