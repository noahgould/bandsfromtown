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

func (gmc *GoogleMapsController) NormalizeLocation(location dal.Location) (*dal.Location, error) {

	if location.FullLocation == "nil, nil, nil" {
		log.Println("wiki couldn't find location.")
		location.FullLocation = "location could not be found"
		location.GooglePlaceID = "-1"
		location.ID = 0
		return &location, nil
	}
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
		log.Printf("loc string: %v. Error: %v.\n", locationString, err)
		normalizedLocation = location
		normalizedLocation.FullLocation = "location could not be found"
		normalizedLocation.GooglePlaceID = "-1"
		normalizedLocation.ID = 0
		return &normalizedLocation, err
	}

	normalizedLocation = dal.Location{
		FullLocation:  placeResult.Predictions[0].Description,
		GooglePlaceID: placeResult.Predictions[0].PlaceID,
	}

	if len(placeResult.Predictions[0].Terms) < 3 {
		if len(placeResult.Predictions[0].Terms) > 1 {
			normalizedLocation.City = placeResult.Predictions[0].Terms[0].Value
			normalizedLocation.State = ""
			normalizedLocation.Country = placeResult.Predictions[0].Terms[1].Value
		} else {
			normalizedLocation.City = ""
			normalizedLocation.State = ""
			normalizedLocation.Country = placeResult.Predictions[0].Terms[0].Value
		}
	} else {
		normalizedLocation.City = placeResult.Predictions[0].Terms[0].Value
		normalizedLocation.State = placeResult.Predictions[0].Terms[1].Value
		normalizedLocation.Country = placeResult.Predictions[0].Terms[2].Value
	}

	return &normalizedLocation, nil
}

//GeocodeLocation returns a location object with long/lat coordinates as well as a normalized name, neighborhood, city, state and country.
func (gmc *GoogleMapsController) GeocodeLocation(location dal.Location) (dal.Location, error) {

	if location.FullLocation == "nil, nil, nil" {
		location.FullLocation = "location could not be found"
		location.GooglePlaceID = "-1"
		location.ID = 0
		return location, nil
	}

	place := &maps.GeocodingRequest{
		Address: location.FullLocation,
	}

	placeResult, err := gmc.mapsClient.Geocode(context.Background(), place)

	if err != nil {
		log.Println(err)

		location.FullLocation = "location could not be found"
		location.GooglePlaceID = "-1"
		location.ID = 0

		return location, err
	}

	var sublocality, adminLevelThree string
	for _, segment := range placeResult[0].AddressComponents {
		for _, segmentType := range segment.Types {
			switch segmentType {
			case "locality":
				location.City = segment.LongName
			case "country":
				location.Country = segment.LongName
			case "administrative_area_level_1":
				location.State = segment.LongName
			case "neighborhood":
				location.Neighborhood = segment.LongName
			case "sublocality_level_1":
				sublocality = segment.LongName
			case "administrative_area_level_3":
				adminLevelThree = segment.LongName
			}
		}
	}

	if location.Neighborhood == "" {
		if sublocality != "" {
			location.Neighborhood = sublocality
		} else if adminLevelThree != "" {
			location.Neighborhood = adminLevelThree
		}
	}

	location.Latitude = placeResult[0].Geometry.Location.Lat
	location.Longitude = placeResult[0].Geometry.Location.Lng

	return location, nil
}

func (gmc *GoogleMapsController) GetCoordinates(location dal.Location) (*dal.Location, error) {

	if location.GooglePlaceID == "" || location.GooglePlaceID == "-1" {
		return nil, errors.New("location does not have a valid place id")
	}

	place := &maps.GeocodingRequest{
		PlaceID: location.GooglePlaceID,
	}

	placeResult, err := gmc.mapsClient.ReverseGeocode(context.Background(), place)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	location.Latitude = placeResult[0].Geometry.Location.Lat
	location.Longitude = placeResult[0].Geometry.Location.Lng

	return &location, nil

}
