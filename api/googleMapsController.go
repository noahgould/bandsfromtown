package api

import (
	"context"
	"errors"
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

	var normalizedLocation dal.Location
	if err != nil {
		log.Print(err)
		normalizedLocation = location
		normalizedLocation.FullLocation = "location could not be found"
		normalizedLocation.GooglePlaceID = "-1"
		normalizedLocation.ID = 0

	} else {
		normalizedLocation = dal.Location{
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

	}

	return &normalizedLocation
}

func (gmc *GoogleMapsController) GetCoordinates(location dal.Location) (*dal.Location, error) {

	if location.GooglePlaceID == "" || location.GooglePlaceID == "-1" {
		return nil, errors.New("Location does not have a valid place id.")
	}

	place := &maps.GeocodingRequest{
		PlaceID: location.GooglePlaceID,
	}

	placeResult, err := gmc.mapsClient.Geocode(context.Background(), place)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	location.Latitude = placeResult[0].Geometry.Location.Lat
	location.Longitude = placeResult[0].Geometry.Location.Lng

	return location, nil

}
